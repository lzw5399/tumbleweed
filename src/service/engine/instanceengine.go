/**
 * @Author: lzw5399
 * @Date: 2021/3/17 23:15
 * @Desc:
 */
package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"workflow/src/global"
	"workflow/src/global/constant"
	"workflow/src/model"
	"workflow/src/model/request"
	"workflow/src/util"
)

type InstanceEngine struct {
	ProcessInstance     model.ProcessInstance   // 流程实例
	ProcessDefinition   model.ProcessDefinition // 流程定义
	definitionStructure DefinitionStructure     // 流程定义中的结构
	currentUserId       uint                    // 当前用户id
	tenantId            uint                    //租户id
}

func NewInstanceEngine(p model.ProcessDefinition, currentUserId uint, tenantId uint) (*InstanceEngine, error) {
	var definitionStructure DefinitionStructure
	err := json.Unmarshal(p.Structure, &definitionStructure)
	if err != nil {
		return nil, err
	}

	return &InstanceEngine{
		ProcessDefinition:   model.ProcessDefinition{},
		definitionStructure: definitionStructure,
		currentUserId:       currentUserId,
		tenantId:            tenantId,
	}, nil
}

func NewInstanceEngineByInstanceId(processInstanceId uint, currentUserId uint, tenantId uint) (*InstanceEngine, error) {
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

	var definitionStructure DefinitionStructure
	err = json.Unmarshal(processDefinition.Structure, &definitionStructure)
	if err != nil {
		return nil, err
	}

	return &InstanceEngine{
		ProcessInstance:     processInstance,
		ProcessDefinition:   processDefinition,
		definitionStructure: definitionStructure,
		currentUserId:       currentUserId,
		tenantId:            tenantId,
	}, nil
}

// 获取instance的初始state
func (i *InstanceEngine) GetInstanceInitialState() ([]map[string]interface{}, error) {
	startNode := i.definitionStructure["nodes"][0]
	startNodeId := startNode["id"].(string)
	var firstEdge map[string]interface{}

	// 获取firstEdge
	for _, edge := range i.definitionStructure["edges"] {
		if edge["source"].(string) == startNodeId {
			firstEdge = edge
			break
		}
	}
	if firstEdge == nil {
		return nil, errors.New("流程模板结构不合法, 请检查初始流程节点和初始顺序流")
	}

	firstEdgeTargetId := firstEdge["target"].(string)
	var nextNode map[string]interface{}
	// 获取接下来的节点nextNode
	for _, node := range i.definitionStructure["nodes"] {
		if node["id"].(string) == firstEdgeTargetId {
			nextNode = node
			break
		}
	}
	if nextNode == nil {
		return nil, errors.New("流程模板结构不合法, 请检查初始流程节点和初始顺序流")
	}

	// 获取初始的states
	return i.GenStates([]map[string]interface{}{nextNode})
}

// 验证入参合法性
func (i *InstanceEngine) ValidateHandleRequest(r *request.HandleInstancesRequest) error {
	currentEdge, err := i.GetEdge(r.EdgeID)
	if err != nil {
		return err
	}

	var currentInstanceState []map[string]interface{}
	err = json.Unmarshal(i.ProcessInstance.State, &currentInstanceState)
	if err != nil {
		return errors.New("当前processInstance的state状态不合法, 请检查")
	}

	// todo 这里先判断[0]
	state := currentInstanceState[0]
	if currentEdge["source"].(string) != state["id"].(string) {
		return errors.New("当前审批不合法, 请检查")
	}

	// 判断当前流程实例状态是否已结束或者被否决
	if i.ProcessInstance.IsEnd {
		return errors.New("当前流程已结束, 不能进行审批操作")
	}

	if i.ProcessInstance.IsDenied {
		return errors.New("当前流程已被否决, 不能进行审批操作")
	}

	// 判断当前用户是否有权限
	hasPermission := i.EnsurePermission(state)
	if !hasPermission {
		return errors.New("当前用户无权限进行当前操作")
	}

	return nil
}

// 验证否决请求的入参
func (i *InstanceEngine) ValidateDenyRequest() error {
	var currentInstanceState []map[string]interface{}
	err := json.Unmarshal(i.ProcessInstance.State, &currentInstanceState)
	if err != nil {
		return errors.New("当前processInstance的state状态不合法, 请检查")
	}

	// todo 这里先判断[0]
	state := currentInstanceState[0]

	// 判断当前流程实例状态是否已结束或者被否决
	if i.ProcessInstance.IsEnd {
		return errors.New("当前流程已结束, 不能进行审批操作")
	}

	if i.ProcessInstance.IsDenied {
		return errors.New("当前流程已被否决, 不能进行审批操作")
	}

	// 判断是否有权限
	hasPermission := i.EnsurePermission(state)
	if !hasPermission {
		return errors.New("当前用户无权限进行当前操作")
	}

	return nil
}

// 判断当前用户是否有权限
func (i *InstanceEngine) EnsurePermission(state map[string]interface{}) bool {
	// 判断当前角色是否有权限
	processors, succeed := state["processor"].([]interface{})
	if !succeed {
		return false
	}

	hasPermission := false
	for _, processor := range processors {
		if uint(processor.(float64)) == i.currentUserId {
			hasPermission = true
			break
		}
	}

	return hasPermission
}

// 流程处理
func (i *InstanceEngine) Handle(r *request.HandleInstancesRequest) error {
	edge, err := i.GetEdge(r.EdgeID)
	if err != nil {
		return err
	}

	targetNode, err := i.GetTargetNodeByEdgeId(r.EdgeID)
	if err != nil {
		return err
	}

	// 判断目标节点的类型，有不同的处理方式
	switch targetNode["clazz"].(string) {
	case constant.UserTask, constant.End:
		newStates, err := i.GenStates([]map[string]interface{}{targetNode})
		if err != nil {
			return err
		}
		err = i.CommonProcessing(edge, targetNode, newStates)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("目前的下一步节点类型：%v，暂不支持", targetNode["clazz"])
	}

	// 获取上一条的流转历史的CreateTime来计算CostDuration
	var lastCirculation model.CirculationHistory
	err = global.BankDb.
		Where("process_instance_id = ?", r.ProcessInstanceId).
		Order("create_time desc").
		Select("create_time").
		First(&lastCirculation).
		Error
	if err != nil {
		return err
	}
	duration := util.FmtDuration(time.Since(lastCirculation.CreateTime))

	// 创建新的一条流转历史
	sourceNode, _ := i.GetNode(edge["source"].(string))
	cirHistory := model.CirculationHistory{
		AuditableBase: model.AuditableBase{
			CreateBy: i.currentUserId,
			UpdateBy: i.currentUserId,
		},
		Title:             i.ProcessInstance.Title,
		ProcessInstanceId: i.ProcessInstance.Id,
		SourceState:       sourceNode["label"].(string),
		SourceId:          sourceNode["id"].(string),
		TargetId:          targetNode["id"].(string),
		Circulation:       edge["label"].(string),
		ProcessorId:       i.currentUserId,
		CostDuration:      duration,
		Remarks:           r.Remarks,
	}

	err = global.BankDb.
		Model(&model.CirculationHistory{}).
		Create(&cirHistory).
		Error

	if err != nil {
		return err
	}

	return nil
}

// 获取edge
func (i *InstanceEngine) GetEdge(edgeId string) (map[string]interface{}, error) {
	if i.definitionStructure["edges"] == nil {
		return nil, errors.New("当前模板结构不合法, 缺少edges, 请检查")
	}

	for _, edge := range i.definitionStructure["edges"] {
		if edge["id"].(string) == edgeId {
			return edge, nil
		}
	}

	return nil, fmt.Errorf("当前edgeId为:%s的edge不存在", edgeId)
}

// 获取node
func (i *InstanceEngine) GetNode(nodeId string) (map[string]interface{}, error) {
	if i.definitionStructure["nodes"] == nil {
		return nil, errors.New("当前模板结构不合法, 缺少nodes, 请检查")
	}

	for _, edge := range i.definitionStructure["nodes"] {
		if edge["id"].(string) == nodeId {
			return edge, nil
		}
	}

	return nil, fmt.Errorf("当前nodeId为:%s的node不存在", nodeId)
}

// 获取edge上的targetNode
func (i *InstanceEngine) GetTargetNodeByEdgeId(edgeId string) (map[string]interface{}, error) {
	edge, err := i.GetEdge(edgeId)
	if err != nil {
		return nil, err
	}

	return i.GetNode(edge["target"].(string))
}

// 获取当前ProcessInstance的State
func (i *InstanceEngine) GetCurrentInstanceState() ([]map[string]interface{}, error) {
	var currentInstanceStates []map[string]interface{}
	err := json.Unmarshal(i.ProcessInstance.State, &currentInstanceStates)

	return currentInstanceStates, err
}

// 获取数据库process_instance表存储的state字段的对象
func (i *InstanceEngine) GenStates(nodes []map[string]interface{}) ([]map[string]interface{}, error) {
	states := make([]map[string]interface{}, 0)
	for _, node := range nodes {
		state := make(map[string]interface{})
		state["id"] = node["id"]
		state["label"] = node["label"]

		switch node["assignType"].(string) {
		case "role": // 审批者是role, 需要转成person
			state["processMethod"] = "person"
			processors, err := i.getUserIdsByRoleIds(node["assignValue"])
			if err != nil {
				return nil, err
			}
			state["processor"] = processors
			state["originProcessMethod"] = node["assignType"]
			state["originProcessor"] = node["assignValue"]
			break
		case "person": // 审批者是person的话直接用原值
			state["processMethod"] = node["assignType"]
			state["processor"] = node["assignValue"]
			break
		default:
			return nil, fmt.Errorf("不支持的处理人类型: %s", node["assignType"].(string))
		}

		// 获取可用的edge
		availableEdges := make([]map[string]interface{}, 0, 1)
		for _, edge := range i.definitionStructure["edges"] {
			if edge["source"].(string) == node["id"].(string) {
				availableEdges = append(availableEdges, edge)
			}
		}
		state["availableEdges"] = availableEdges

		states = append(states, state)
	}

	return states, nil
}

func (i *InstanceEngine) getUserIdsByRoleIds(ids interface{}) ([]int, error) {
	bytes := util.MarshalToBytes(ids)
	var roleIds []int
	err := json.Unmarshal(bytes, &roleIds)
	if err != nil {
		return nil, err
	}

	var roleUsersList []model.RoleUsers
	err = global.BankDb.
		Model(&model.RoleUsers{}).
		Where("tenant_id = ?", i.tenantId).
		Where("role_id in ?", roleIds).
		Find(&roleUsersList).
		Error
	if err != nil {
		log.Printf("查询roleuser失败，原因:%s", err.Error())
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
