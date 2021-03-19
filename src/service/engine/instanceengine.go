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
}

func NewInstanceEngine(p model.ProcessDefinition, currentUserId uint) (*InstanceEngine, error) {
	var definitionStructure DefinitionStructure
	err := json.Unmarshal(p.Structure, &definitionStructure)
	if err != nil {
		return nil, err
	}

	return &InstanceEngine{
		ProcessDefinition:   model.ProcessDefinition{},
		definitionStructure: definitionStructure,
		currentUserId:       currentUserId,
	}, nil
}

func NewInstanceEngineByInstanceId(processInstanceId uint, currentUserId uint) (*InstanceEngine, error) {
	var processInstance model.ProcessInstance
	var processDefinition model.ProcessDefinition

	err := global.BankDb.
		Model(model.ProcessInstance{}).
		Where("id = ?", processInstanceId).
		First(&processInstance).
		Error
	if err != nil {
		return nil, fmt.Errorf("找不到当前processInstanceId为 %v 的记录", processInstanceId)
	}

	err = global.BankDb.
		Model(model.ProcessDefinition{}).
		Where("id = ?", processInstance.ProcessDefinitionId).
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
	initialStates := i.GenStates([]map[string]interface{}{nextNode})

	return initialStates, nil
}

// 验证入参合法性
func (i *InstanceEngine) ValidateHandleRequest(r *request.HandleInstancesRequest, currentUserId uint) error {
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
		return fmt.Errorf("当前传入的edgeId: %v 不合法, 请检查", r.EdgeID)
	}

	// 判断当前角色是否有权限
	processors, succeed := state["processor"].([]interface{})
	if !succeed {
		return errors.New("当前processInstance的state状态不合法, 请检查")
	}

	hasPermission := false
	for _, processor := range processors {
		if uint(processor.(float64)) == currentUserId {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return errors.New("当前用户无权限进行当前操作")
	}

	return nil
}

// 根据流程处理
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
	switch targetNode["clazz"] {
	case constant.USER_TASK:
		newStates := i.GenStates([]map[string]interface{}{targetNode})
		return i.CommonProcessing(edge, targetNode, newStates)
	}

	// 更新上一条的流转历史的CostDuration
	var lastCirculation model.CirculationHistory
	err = global.BankDb.
		Model(&model.CirculationHistory{}).
		Where("process_instance_id = ?", r.ProcessInstanceId).
		Order("create_time desc").
		Find(lastCirculation).
		Limit(1).
		Error
	if err != nil {
		return err
	}
	duration := util.FmtDuration(time.Since(lastCirculation.CreateTime))
	err = global.BankDb.
		Model(&model.CirculationHistory{}).
		Where("id = ?", lastCirculation.Id).
		Updates(map[string]interface{}{
			"cost_duration": duration,
			"update_time":   time.Now().Local(),
			"update_by":     i.currentUserId,
		}).
		Error
	if err != nil {
		return err
	}

	// 创建新的一条流转历史
	//cirHistory := model.CirculationHistory{
	//	Title:             h.workOrderDetails.Title,
	//	ProcessInstanceId: 0,
	//	SourceState:       "",
	//	SourceId:          "",
	//	TargetId:          "",
	//	Circulation:       circulationValue,
	//	ProcessorId:       tools.GetUserId(c),
	//	CostDuration:      costDurationValue,
	//	Remarks:           remarks,
	//}

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
func (i *InstanceEngine) GenStates(nodes []map[string]interface{}) []map[string]interface{} {
	states := make([]map[string]interface{}, 0)
	for _, node := range nodes {
		state := make(map[string]interface{})
		state["id"] = node["id"]
		state["processMethod"] = node["assignType"]
		state["processor"] = node["assignValue"]
		state["label"] = node["label"]

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

	return states
}
