/**
 * @Author: lzw5399
 * @Date: 2021/3/17 23:15
 * @Desc:
 */
package engine

import (
	"errors"
	"fmt"

	. "github.com/ahmetb/go-linq/v3"
	"gorm.io/gorm"

	"workflow/src/global"
	"workflow/src/global/constant"
	"workflow/src/model"
	"workflow/src/model/dto"
	"workflow/src/model/request"
	"workflow/src/util"
)

type ProcessEngine struct {
	tx                  *gorm.DB                // 数据库事务对象
	userIdentifier      string                  // 当前用户外部系统的id
	tenantId            int                     // 当前租户id
	sourceNode          *dto.Node               // 流转的源node
	targetNode          *dto.Node               // 流转的目标node
	linkEdge            *dto.Edge               // sourceNode和targetNode中间连接的edge
	ProcessInstance     model.ProcessInstance   // 流程实例
	ProcessDefinition   model.ProcessDefinition // 流程定义
	DefinitionStructure dto.Structure           // ProcessDefinition.Structure的快捷方式
}

// 初始化流程引擎
func NewProcessEngine(p model.ProcessDefinition, instance model.ProcessInstance, userIdentifier string, tenantId int, tx *gorm.DB) (*ProcessEngine, error) {
	return &ProcessEngine{
		ProcessDefinition:   p,
		ProcessInstance:     instance,
		DefinitionStructure: p.Structure,
		userIdentifier:      userIdentifier,
		tenantId:            tenantId,
		tx:                  tx,
	}, nil
}

// 初始化流程引擎(带process_instance)
func NewProcessEngineByInstanceId(processInstanceId int, userIdentifier string, tenantId int, tx *gorm.DB) (*ProcessEngine, error) {
	var processInstance model.ProcessInstance
	var processDefinition model.ProcessDefinition

	err := global.BankDb.
		Model(model.ProcessInstance{}).
		Where("id = ?", processInstanceId).
		Where("tenant_id = ?", tenantId).
		First(&processInstance).
		Error
	if err != nil {
		return nil, fmt.Errorf("找不到当前processInstanceId为 %v 的记录", processInstanceId)
	}

	err = global.BankDb.
		Model(model.ProcessDefinition{}).
		Where("id = ?", processInstance.ProcessDefinitionId).
		Where("tenant_id = ?", tenantId).
		First(&processDefinition).
		Error
	if err != nil {
		return nil, fmt.Errorf("找不到当前processDefinitionId为 %v 的记录", processInstance.ProcessDefinitionId)
	}

	return &ProcessEngine{
		ProcessInstance:     processInstance,
		ProcessDefinition:   processDefinition,
		DefinitionStructure: processDefinition.Structure,
		userIdentifier:      userIdentifier,
		tenantId:            tenantId,
		tx:                  tx,
	}, nil
}

// 获取instance的初始state
func (engine *ProcessEngine) GetInstanceInitialState() (dto.StateArray, error) {
	var startNode dto.Node
	for _, node := range engine.DefinitionStructure.Nodes {
		if node.Clazz == constant.START {
			startNode = node
			break
		}
	}

	// 获取firstEdge
	firstEdge := dto.Edge{}
	for _, edge := range engine.DefinitionStructure.Edges {
		if edge.Source == startNode.Id {
			firstEdge = edge
			break
		}
	}

	if firstEdge.Id == "" {
		return nil, errors.New("流程模板结构不合法, 请检查初始流程节点和初始顺序流")
	}

	firstEdgeTargetId := firstEdge.Target
	nextNode := dto.Node{}
	// 获取接下来的节点nextNode
	for _, node := range engine.DefinitionStructure.Nodes {
		if node.Id == firstEdgeTargetId {
			nextNode = node
			break
		}
	}
	if nextNode.Id == "" {
		return nil, errors.New("流程模板结构不合法, 请检查初始流程节点和初始顺序流")
	}

	// 获取初始的states
	return engine.GenNewStates([]dto.Node{nextNode})
}

// 流程处理
func (engine *ProcessEngine) Handle(r *request.HandleInstancesRequest) error {
	// 获取edge
	edge, err := engine.GetEdge(r.EdgeId)
	if err != nil {
		return err
	}

	// 获取两个node
	sourceNode, _ := engine.GetNode(edge.Source)
	targetNode, err := engine.GetTargetNodeByEdgeId(r.EdgeId)
	if err != nil {
		return err
	}

	// 设置当前的节点和顺序流信息
	engine.SetCurrentNodeEdgeInfo(&sourceNode, &edge, &targetNode)
	engine.UpdateRelatedPerson()

	// handle内部(有递归操作，针对比如网关后还是网关等场景)
	return engine.handleInternal(r, 1)
}

// 流程处理内部(用于递归)
func (engine *ProcessEngine) handleInternal(r *request.HandleInstancesRequest, deepLevel int) error {
	// 判断当前节点是否会签
	isCounterSign, isLastProcessor, err := engine.JudgeCounterSign()
	if err != nil {
		return err
	}

	// 添加流转历史记录
	err = engine.CreateHistory(r.Remarks, false)
	if err != nil {
		return err
	}

	// 是会签并且不是最后一个人
	// 则不需要判下面目标节点相关的逻辑,直接退出
	if isCounterSign && !isLastProcessor {
		return nil
	}

	// 递归中, 当sourceNode为【结束事件】的情况下targetNode会为空
	if engine.targetNode == nil {
		return nil
	}

	// 判断目标节点的类型，有不同的处理方式
	switch engine.targetNode.Clazz {
	case constant.UserTask:
		// 只有第一次进来，针对userTask才需要跳转. 后续的直接退出
		if deepLevel > 1 {
			break
		}
		newStates, err := engine.MergeStates(engine.sourceNode.Id, []dto.Node{*engine.targetNode})
		if err != nil {
			return err
		}
		err = engine.Circulation(newStates)
		if err != nil {
			return err
		}

	case constant.End:
		// 只有第一次进来，才需要Circulation跳转
		// 非第一次的递归记一条结束的日志就退出
		if deepLevel == 1 {
			newStates, err := engine.MergeStates(engine.sourceNode.Id, []dto.Node{*engine.targetNode})
			if err != nil {
				return err
			}
			err = engine.Circulation(newStates)
			if err != nil {
				return err
			}
		}
		engine.SetCurrentNodeEdgeInfo(engine.targetNode, nil, nil)
		return engine.handleInternal(r, deepLevel+1)

	case constant.ExclusiveGateway:
		err := engine.ProcessingExclusiveGateway(*engine.targetNode, r)
		if err != nil {
			return err
		}

		// 递归处理
		return engine.handleInternal(r, deepLevel+1)

	case constant.ParallelGateway:
		relationInfos, err := engine.ProcessParallelGateway()
		if err != nil {
			return err
		}

		// 递归处理
		for _, info := range relationInfos {
			engine.SetCurrentNodeEdgeInfo(&info.SourceNode, &info.LinkedEdge, &info.TargetNode)
			err = engine.handleInternal(r, deepLevel+1)
			if err != nil {
				return err
			}
		}
		return nil

	default:
		return fmt.Errorf("目前的下一步节点类型：%v，暂不支持", engine.targetNode.Clazz)
	}

	return nil
}

func (engine *ProcessEngine) SetCurrentNodeEdgeInfo(sourceNode *dto.Node, edge *dto.Edge, targetNode *dto.Node) {
	engine.sourceNode = sourceNode
	engine.linkEdge = edge
	engine.targetNode = targetNode
}

// 获取edge
func (engine *ProcessEngine) GetEdge(edgeId string) (dto.Edge, error) {
	if len(engine.DefinitionStructure.Edges) == 0 {
		return dto.Edge{}, errors.New("当前模板结构不合法, 缺少edges, 请检查")
	}

	for _, edge := range engine.DefinitionStructure.Edges {
		if edge.Id == edgeId {
			return edge, nil
		}
	}

	return dto.Edge{}, fmt.Errorf("当前edgeId为:%s的edge不存在", edgeId)
}

// i.GetEdges("userTask123", "source") 获取所有source为userTask123的edges
// i.GetEdges("userTask123", "target") 获取所有target为userTask123的edges
func (engine *ProcessEngine) GetEdges(nodeId string, nodeIdType string) []dto.Edge {
	edges := make([]dto.Edge, 0)
	for _, edge := range engine.DefinitionStructure.Edges {
		switch nodeIdType {
		case "source":
			if edge.Source == nodeId {
				edges = append(edges, edge)
			}
		case "target":
			if edge.Target == nodeId {
				edges = append(edges, edge)
			}
		}
	}

	// 根据sort排序
	From(edges).OrderByT(func(i dto.Edge) interface{} {
		return i.Sort
	})

	return edges
}

// 获取node
func (engine *ProcessEngine) GetNode(nodeId string) (dto.Node, error) {
	if len(engine.DefinitionStructure.Nodes) == 0 {
		return dto.Node{}, errors.New("当前模板结构不合法, 缺少nodes, 请检查")
	}

	for _, node := range engine.DefinitionStructure.Nodes {
		if node.Id == nodeId {
			return node, nil
		}
	}

	return dto.Node{}, fmt.Errorf("当前nodeId为:%s的node不存在", nodeId)
}

// GetNodes(edges, "source") 获取edges中source属性指向的node的集合
// GetNodes(edges, "target") 获取edges中target属性指向的node的集合
func (engine *ProcessEngine) GetNodesByEdges(edges []dto.Edge, edgeType string) []dto.Node {
	nodes := make([]dto.Node, 0)
	for _, node := range engine.DefinitionStructure.Nodes {
		for _, edge := range edges {
			switch edgeType {
			case "source":
				if edge.Source == node.Id {
					nodes = append(nodes, node)
				}
			case "target":
				if edge.Target == node.Id {
					nodes = append(nodes, node)
				}
			}
		}
	}

	return nodes
}

// 获取edge上的targetNode
func (engine *ProcessEngine) GetTargetNodeByEdgeId(edgeId string) (dto.Node, error) {
	edge, err := engine.GetEdge(edgeId)
	if err != nil {
		return dto.Node{}, err
	}

	return engine.GetNode(edge.Target)
}

// 获取全新的 数据库process_instance表存储的state字段的对象
func (engine *ProcessEngine) GenNewStates(nodes []dto.Node) (dto.StateArray, error) {
	states := dto.StateArray{}
	for _, node := range nodes {
		state := dto.State{
			Id:                 node.Id,
			Label:              node.Label,
			Processor:          []string{},
			CompletedProcessor: []string{},
			ProcessMethod:      node.AssignType,  // 处理方式(角色 用户等)
			AssignValue:        node.AssignValue, // 指定的处理者(用户的id或者角色的id)
			AvailableEdges:     []dto.Edge{},
			IsCounterSign:      node.IsCounterSign,
		}

		// 审批者是role的需要在这里转成person
		if node.AssignType != "" && node.AssignValue != nil {
			switch node.AssignType {
			case "role": // 审批者是role, 需要转成person
				processors, err := engine.GetUserIdsByRoleIds(node.AssignValue)
				if err != nil {
					return nil, err
				}
				state.Processor = processors
				state.UnCompletedProcessor = processors
				break
			case "person": // 审批者是person的话直接用原值
				state.Processor = node.AssignValue
				state.UnCompletedProcessor = node.AssignValue
				break
			default:
				return nil, fmt.Errorf("不支持的处理人类型: %s", node.AssignType)
			}
		}

		// 获取可用的edge
		availableEdges := make([]dto.Edge, 0, 1)
		for _, edge := range engine.DefinitionStructure.Edges {
			if edge.Source == node.Id {
				availableEdges = append(availableEdges, edge)
			}
		}
		state.AvailableEdges = availableEdges

		states = append(states, state)
	}

	return states, nil
}

// 合并state字段，返回合并之后的
func (engine *ProcessEngine) MergeStates(removeNodeId string, newNodes []dto.Node) (dto.StateArray, error) {
	finalMergedStates := make(dto.StateArray, 0)

	// 原来的state排除掉removeNodeId
	for _, state := range engine.ProcessInstance.State {
		if state.Id == removeNodeId {
			continue
		}
		finalMergedStates = append(finalMergedStates, state)
	}

	// 新的node生成state
	newStates, err := engine.GenNewStates(newNodes)
	if err != nil {
		return nil, err
	}

	// 合并
	finalMergedStates = append(finalMergedStates, newStates...)

	return finalMergedStates, nil
}

// 通过角色id获取用户id
func (engine *ProcessEngine) GetUserIdsByRoleIds(roleIds []string) ([]string, error) {
	var userIds []string
	err := global.BankDb.Model(&model.Role{}).
		Joins("inner join wf.user_role on user_role.role_identifier = role.identifier").
		Joins("inner join wf.user on user_role.user_identifier = user.identifier").
		Where("role.tenant_id = ? and role.identifier in ?", engine.tenantId, roleIds).
		Select("user.identifier").
		Scan(&userIds).
		Error

	return userIds, err
}

// 合并更新变量
func (engine *ProcessEngine) MergeVariables(newVariables []model.InstanceVariable) {
	// 反序列化出来
	originVariables := util.UnmarshalToInstanceVariables(engine.ProcessInstance.Variables)

	// 查询优化先整理成map
	originVariableMap := make(map[string]model.InstanceVariable, len(originVariables))
	for _, v := range originVariables {
		originVariableMap[v.Name] = v
	}

	for _, v := range newVariables {
		originVariableMap[v.Name] = v
	}

	finalVariables := make([]model.InstanceVariable, 0)
	for _, v := range originVariableMap {
		finalVariables = append(finalVariables, v)
	}

	engine.ProcessInstance.Variables = util.MarshalToDbJson(finalVariables)
}

// 获取初始节点
func (engine *ProcessEngine) GetInitialNode() (dto.Node, error) {
	startNode := dto.Node{}
	for _, node := range engine.DefinitionStructure.Nodes {
		if node.Clazz == constant.START {
			startNode = node
		}
	}

	if startNode.Id == "" {
		return startNode, errors.New("当前结构不合法")
	}

	return startNode, nil
}

// 通过sourceNode获取TargetNodes列表
func (engine *ProcessEngine) GetTargetNodes(sourceNode dto.Node) ([]dto.Node, error) {
	edges := engine.GetEdges(sourceNode.Id, "source")

	nextNodes := make([]dto.Node, 0, 1)

	for _, edge := range edges {
		node, err := engine.GetTargetNodeByEdgeId(edge.Id)
		if err != nil {
			return nil, err
		}

		nextNodes = append(nextNodes, node)
	}

	return nextNodes, nil
}

// 根据nodeId获取state
func (engine *ProcessEngine) GetStateByNodeId(nodeId string) (dto.State, error) {
	state := dto.State{}
	for _, s := range engine.ProcessInstance.State {
		if s.Id == nodeId {
			state = s
		}
	}

	if state.Id == "" {
		return state, errors.New("当前流程状态错误，对应的nodeId错误")
	}

	return state, nil
}

// 通过edgeId获取state
func (engine *ProcessEngine) GetStateByEdgeId(edge dto.Edge) (dto.State, error) {
	state := dto.State{}
	for _, s := range engine.ProcessInstance.State {
		if s.Id == edge.Source {
			state = s
			break
		}
	}

	if state.Id == "" {
		return state, util.BadRequest.New("当前审批不合法, 请检查")
	}

	return state, nil
}
