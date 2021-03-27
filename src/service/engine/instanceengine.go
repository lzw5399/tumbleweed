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

type InstanceEngine struct {
	tx                  *gorm.DB                // 数据库事务对象
	currentUserId       uint                    // 当前用户id
	tenantId            uint                    // 当前租户id
	sourceNode          *dto.Node               // 流转的源node
	targetNode          *dto.Node               // 流转的目标node
	linkEdge            *dto.Edge               // sourceNode和targetNode中间连接的edge
	ProcessInstance     model.ProcessInstance   // 流程实例
	ProcessDefinition   model.ProcessDefinition // 流程定义
	DefinitionStructure dto.Structure           // ProcessDefinition.Structure的快捷方式
}

// 初始化流程引擎
func NewInstanceEngine(p model.ProcessDefinition, instance model.ProcessInstance, currentUserId uint, tenantId uint, tx *gorm.DB) (*InstanceEngine, error) {
	return &InstanceEngine{
		ProcessDefinition:   p,
		ProcessInstance:     instance,
		DefinitionStructure: p.Structure,
		currentUserId:       currentUserId,
		tenantId:            tenantId,
		tx:                  tx,
	}, nil
}

// 初始化流程引擎(带process instance)
func NewInstanceEngineByInstanceId(processInstanceId uint, currentUserId uint, tenantId uint, tx *gorm.DB) (*InstanceEngine, error) {
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

	return &InstanceEngine{
		ProcessInstance:     processInstance,
		ProcessDefinition:   processDefinition,
		DefinitionStructure: processDefinition.Structure,
		currentUserId:       currentUserId,
		tenantId:            tenantId,
		tx:                  tx,
	}, nil
}

// 获取instance的初始state
func (i *InstanceEngine) GetInstanceInitialState() (dto.StateArray, error) {
	var startNode dto.Node
	for _, node := range i.DefinitionStructure.Nodes {
		if node.Clazz == constant.START {
			startNode = node
		}
	}

	// 获取firstEdge
	firstEdge := dto.Edge{}
	for _, edge := range i.DefinitionStructure.Edges {
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
	for _, node := range i.DefinitionStructure.Nodes {
		if node.Id == firstEdgeTargetId {
			nextNode = node
			break
		}
	}
	if nextNode.Id == "" {
		return nil, errors.New("流程模板结构不合法, 请检查初始流程节点和初始顺序流")
	}

	// 获取初始的states
	return i.GenStates([]dto.Node{nextNode})
}

// 流程处理
func (i *InstanceEngine) Handle(r *request.HandleInstancesRequest) error {
	// 获取edge
	edge, err := i.GetEdge(r.EdgeID)
	if err != nil {
		return err
	}

	// 获取两个node
	sourceNode, _ := i.GetNode(edge.Source)
	targetNode, err := i.GetTargetNodeByEdgeId(r.EdgeID)
	if err != nil {
		return err
	}

	// 设置当前的节点和顺序流信息
	i.SetCurrentNodeEdgeInfo(&sourceNode, &edge, &targetNode)
	i.UpdateRelatedPerson()

	// 判断当前节点是否会签
	isCounterSign, isLastProcessor, err := i.JudgeCounterSign()
	if err != nil {
		return err
	}

	// 添加历史记录(这条只是保底的, 后续还会有其他的判断)
	err = i.CreateCirculationHistory(r.Remarks)
	if err != nil {
		return err
	}

	// 是会签并且不是最后一个人
	// 则不需要判下面目标节点相关的逻辑,直接退出
	if isCounterSign && !isLastProcessor {
		return nil
	}

	// 判断目标节点的类型，有不同的处理方式
	switch targetNode.Clazz {
	case constant.UserTask, constant.End:
		newStates, err := i.GenStates([]dto.Node{targetNode})
		if err != nil {
			return err
		}
		err = i.Circulation(newStates)
		if err != nil {
			return err
		}
	case constant.ExclusiveGateway:
		err := i.ProcessingExclusiveGateway(targetNode, r)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("目前的下一步节点类型：%v，暂不支持", targetNode.Clazz)
	}

	// TODO: 这里可以跟 【目前的handle, 如果排他网关后面还是排他网关，会有问题】一起优化掉，应该需要递归
	originTargetNode := targetNode
	switch originTargetNode.Clazz {
	case constant.ExclusiveGateway:
		// 由于排他网关理论上会跳至少两次【原节点->排他网关->后续节点】
		// 所以需要再
		err := i.CreateCirculationHistory(r.Remarks)
		if err != nil {
			return err
		}

		// 判断二次是否是end
		if i.targetNode != nil && i.targetNode.Clazz == constant.End {
			i.SetCurrentNodeEdgeInfo(i.targetNode, nil, nil)
			err = i.CreateCirculationHistory(r.Remarks)
			if err != nil {
				return err
			}
		}
	case constant.End:
		i.SetCurrentNodeEdgeInfo(i.targetNode, nil, nil)
		err = i.CreateCirculationHistory(r.Remarks)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *InstanceEngine) SetCurrentNodeEdgeInfo(sourceNode *dto.Node, edge *dto.Edge, targetNode *dto.Node) {
	i.sourceNode = sourceNode
	i.linkEdge = edge
	i.targetNode = targetNode
}

// 获取edge
func (i *InstanceEngine) GetEdge(edgeId string) (dto.Edge, error) {
	if len(i.DefinitionStructure.Edges) == 0 {
		return dto.Edge{}, errors.New("当前模板结构不合法, 缺少edges, 请检查")
	}

	for _, edge := range i.DefinitionStructure.Edges {
		if edge.Id == edgeId {
			return edge, nil
		}
	}

	return dto.Edge{}, fmt.Errorf("当前edgeId为:%s的edge不存在", edgeId)
}

// i.GetEdges("userTask123", "source") 获取所有source为userTask123的edges
// i.GetEdges("userTask123", "target") 获取所有target为userTask123的edges
func (i *InstanceEngine) GetEdges(nodeId string, nodeIdType string) []dto.Edge {
	edges := make([]dto.Edge, 0)
	for _, edge := range i.DefinitionStructure.Edges {
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
func (i *InstanceEngine) GetNode(nodeId string) (dto.Node, error) {
	if len(i.DefinitionStructure.Nodes) == 0 {
		return dto.Node{}, errors.New("当前模板结构不合法, 缺少nodes, 请检查")
	}

	for _, node := range i.DefinitionStructure.Nodes {
		if node.Id == nodeId {
			return node, nil
		}
	}

	return dto.Node{}, fmt.Errorf("当前nodeId为:%s的node不存在", nodeId)
}

// 获取edge上的targetNode
func (i *InstanceEngine) GetTargetNodeByEdgeId(edgeId string) (dto.Node, error) {
	edge, err := i.GetEdge(edgeId)
	if err != nil {
		return dto.Node{}, err
	}

	return i.GetNode(edge.Target)
}

// 获取数据库process_instance表存储的state字段的对象
func (i *InstanceEngine) GenStates(nodes []dto.Node) (dto.StateArray, error) {
	states := dto.StateArray{}
	for _, node := range nodes {
		state := dto.State{
			Id:                 node.Id,
			Label:              node.Label,
			Processor:          []int{},
			CompletedProcessor: []int{},
			ProcessMethod:      node.AssignType,  // 处理方式(角色 用户等)
			AssignValue:        node.AssignValue, // 指定的处理者(用户的id或者角色的id)
			AvailableEdges:     []dto.Edge{},
			IsCounterSign:      node.IsCounterSign,
		}

		// 审批者是role的需要在这里转成person
		if node.AssignType != "" && node.AssignValue != nil {
			switch node.AssignType {
			case "role": // 审批者是role, 需要转成person
				processors, err := i.GetUserIdsByRoleIds(node.AssignValue)
				if err != nil {
					return nil, err
				}
				state.Processor = processors
				break
			case "person": // 审批者是person的话直接用原值
				state.Processor = node.AssignValue
				break
			default:
				return nil, fmt.Errorf("不支持的处理人类型: %s", node.AssignType)
			}
		}

		// 获取可用的edge
		availableEdges := make([]dto.Edge, 0, 1)
		for _, edge := range i.DefinitionStructure.Edges {
			if edge.Source == node.Id {
				availableEdges = append(availableEdges, edge)
			}
		}
		state.AvailableEdges = availableEdges

		states = append(states, state)
	}

	return states, nil
}

// 通过角色id获取用户id
func (i *InstanceEngine) GetUserIdsByRoleIds(roleIds []int) ([]int, error) {
	var roleUsersList []model.RoleUsers
	err := global.BankDb.
		Model(&model.RoleUsers{}).
		Where("tenant_id = ?", i.tenantId).
		Where("role_id in ?", roleIds).
		Find(&roleUsersList).
		Error
	if err != nil {
		global.BankLogger.Errorln("查询role_user失败", err)
		return nil, err
	}

	// 使用map来提高查询效率
	finalUserIdMap := make(map[int64]bool, 0)
	for _, roleUsers := range roleUsersList {
		for _, userId := range roleUsers.UserIds {
			if _, present := finalUserIdMap[userId]; !present {
				finalUserIdMap[userId] = true
			}
		}
	}

	// 转成[]string
	finalUserIds := make([]int, 0)
	for k, _ := range finalUserIdMap {
		finalUserIds = append(finalUserIds, int(k))
	}

	return finalUserIds, nil
}

// 合并更新变量
func (i *InstanceEngine) UpdateVariables(newVariables []model.InstanceVariable) {
	// 反序列化出来
	originVariables := util.UnmarshalToInstanceVariables(i.ProcessInstance.Variables)

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

	i.ProcessInstance.Variables = util.MarshalToDbJson(finalVariables)
}

// 获取初始节点
func (i *InstanceEngine) GetInitialNode() (dto.Node, error) {
	startNode := dto.Node{}
	for _, node := range i.DefinitionStructure.Nodes {
		if node.Clazz == constant.START {
			startNode = node
		}
	}

	if startNode.Id == "" {
		return startNode, errors.New("当前结构不合法")
	}

	return startNode, nil
}

func (i *InstanceEngine) GetNextNodes(sourceNode dto.Node) ([]dto.Node, error) {
	edges := i.GetEdges(sourceNode.Id, "source")

	nextNodes := make([]dto.Node, 0, 1)

	for _, edge := range edges {
		node, err := i.GetTargetNodeByEdgeId(edge.Id)
		if err != nil {
			return nil, err
		}

		nextNodes = append(nextNodes, node)
	}

	return nextNodes, nil
}
