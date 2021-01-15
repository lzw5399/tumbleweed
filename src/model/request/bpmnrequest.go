/**
 * @Author: lzw5399
 * @Date: 2021/1/15 21:24
 * @Desc: 接收bpmn2.0的struct
 */
package request

import (
	"time"
	"workflow/src/model"
)

type BpmnRequest struct {
	Data string `json:"data"`
}

// bpmn2.0的struct
type Definitions struct {
	Process struct {
		ID              string `xml:"id,attr"`              // 流程标识
		Name            string `xml:"name,attr"`            // 流程名字
		ProcessCategory string `xml:"processCategory,attr"` // 流程类别
		StartEvent      struct {
			ID       string `xml:"id,attr"`
			Name     string `xml:"name,attr"`
			Outgoing string `xml:"outgoing"`
		} `xml:"startEvent"` // 开始事件
		SequenceFlow []struct {
			ID                  string `xml:"id,attr"`
			SourceRef           string `xml:"sourceRef,attr"` // 上一个节点
			TargetRef           string `xml:"targetRef,attr"` // 目标节点
			ConditionExpression struct {
				Text string `xml:",chardata"` // 表达式
				Type string `xml:"type,attr"`
			} `xml:"conditionExpression"`
		} `xml:"sequenceFlow"`
		UserTask []struct {
			ID             string   `xml:"id,attr"`
			Name           string   `xml:"name,attr"`
			FormKey        string   `xml:"formKey,attr"`
			Assignee       string   `xml:"assignee,attr"`
			CandidateUsers string   `xml:"candidateUsers,attr"`
			Incoming       []string `xml:"incoming"`
			Outgoing       []string `xml:"outgoing"`
		} `xml:"userTask"`
		ExclusiveGateway []struct {
			ID       string   `xml:"id,attr"`
			Incoming []string `xml:"incoming"`
			Outgoing []string `xml:"outgoing"`
		} `xml:"exclusiveGateway"`
		EndEvent struct {
			ID       string   `xml:"id,attr"`
			Name     string   `xml:"name,attr"`
			Incoming []string `xml:"incoming"`
			Outgoing []string `xml:"outgoing"`
		} `xml:"endEvent"`
	} `xml:"process"`
}

func (d *Definitions) ToProcess() model.ProcessDefinition {
	return model.ProcessDefinition{
		DbBase: model.DbBase{
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		},
		Name:       "",
		Version:    0,
		Resource:   "",
		Userid:     "",
		Username:   "",
		Company:    "",
		DeployTime: "",
	}
}
