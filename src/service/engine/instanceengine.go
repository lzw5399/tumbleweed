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

	"workflow/src/global"
	"workflow/src/model"
	"workflow/src/model/request"
)

type InstanceEngine struct {
	ProcessInstance     model.ProcessInstance
	ProcessDefinition   model.ProcessDefinition
	definitionStructure DefinitionStructure
}

func NewInstanceEngine(p model.ProcessDefinition) (*InstanceEngine, error) {
	var definitionStructure DefinitionStructure
	err := json.Unmarshal(p.Structure, &definitionStructure)
	if err != nil {
		return nil, err
	}

	return &InstanceEngine{
		ProcessDefinition:   model.ProcessDefinition{},
		definitionStructure: definitionStructure,
	}, nil
}

func NewInstanceEngineByInstanceId(processInstanceId uint) (*InstanceEngine, error) {
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
	}, nil
}

// 获取instance的初始state
func (i *InstanceEngine) GetInstanceInitialState() ([]map[string]interface{}, error) {
	states := make([]map[string]interface{}, 1)
	state := make(map[string]interface{})

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

	// 获取接下来的节点nextNode的可用的edge
	availableEdges := make([]map[string]interface{}, 0, 1)
	nextNodeId := nextNode["id"].(string)
	for _, edge := range i.definitionStructure["edges"] {
		if edge["source"].(string) == nextNodeId {
			availableEdges = append(availableEdges, edge)
		}
	}

	state["id"] = nextNode["id"]
	state["processMethod"] = nextNode["assignType"]
	state["processor"] = nextNode["assignValue"]
	state["label"] = nextNode["label"]
	state["availableEdges"] = availableEdges
	states[0] = state

	return states, nil
}

// 验证入参合法性
func (i *InstanceEngine) ValidateHandleRequest(r *request.HandleInstancesRequest, currentUserId uint) error {
	var currentEdge map[string]interface{}
	for _, edge := range i.definitionStructure["edges"] {
		if edge["id"].(string) == r.EdgeID {
			currentEdge = edge
			break
		}
	}
	if currentEdge == nil {
		return fmt.Errorf("入参错误，当前edgeId:%s不存在", r.EdgeID)
	}

	var currentInstanceState []map[string]interface{}
	err := json.Unmarshal(i.ProcessInstance.State, &currentInstanceState)
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
		return fmt.Errorf("当前用户id:%v, 没有权限进行当前操作", currentUserId)
	}

	return nil
}
