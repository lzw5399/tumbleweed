/**
 * @Author: lzw5399
 * @Date: 2021/3/17 23:15
 * @Desc:
 */
package engine

import (
	"encoding/json"
	"errors"

	"workflow/src/model"
)

type InstanceEngine struct {
	processDefinition   model.ProcessDefinition
	definitionStructure DefinitionStructure
}

func NewInstanceEngine(p model.ProcessDefinition) (*InstanceEngine, error) {
	var definitionStructure DefinitionStructure
	err := json.Unmarshal(p.Structure, &definitionStructure)
	if err != nil {
		return nil, err
	}

	return &InstanceEngine{
		processDefinition:   model.ProcessDefinition{},
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

	state["id"] = nextNode["id"]
	state["processMethod"] = nextNode["assignType"]
	state["processor"] = nextNode["assignValue"]
	state["label"] = nextNode["label"]
	states[0] = state

	return states, nil
}
