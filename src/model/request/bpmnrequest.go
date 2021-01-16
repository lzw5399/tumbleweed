/**
 * @Author: lzw5399
 * @Date: 2021/1/15 21:24
 * @Desc: 接收bpmn2.0的struct
 */
package request

import (
	"time"
	"workflow/src/global/constant"
	"workflow/src/model"
)

type BpmnRequest struct {
	Data string `json:"data"`
}

// bpmn2.0的struct
type Definitions struct {
	Process ProcessRequest `xml:"process"`
}

type ProcessRequest struct {
	ID              string `xml:"id,attr"`              // 流程标识
	Name            string `xml:"name,attr"`            // 流程名字
	ProcessCategory string `xml:"processCategory,attr"` // 流程类别
	StartEvent      []struct {
		ID       string   `xml:"id,attr"`
		Name     string   `xml:"name,attr"`
		Incoming []string `xml:"incoming"`
		Outgoing []string `xml:"outgoing"`
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
		ID              string   `xml:"id,attr"`
		Name            string   `xml:"name,attr"`
		FormKey         string   `xml:"formKey,attr"`
		CandidateUsers  string   `xml:"candidateUsers,attr"`  // 候选人, 逗号分割
		CandidateGroups string   `xml:"candidateGroups,attr"` // 候选组， 逗号分割
		Assignee        string   `xml:"assignee,attr"`        // 指定人员
		Incoming        []string `xml:"incoming"`
		Outgoing        []string `xml:"outgoing"`
	} `xml:"userTask"`
	ExclusiveGateway []struct {
		ID       string   `xml:"id,attr"`
		Incoming []string `xml:"incoming"`
		Outgoing []string `xml:"outgoing"`
	} `xml:"exclusiveGateway"`
	EndEvent []struct {
		ID       string   `xml:"id,attr"`
		Name     string   `xml:"name,attr"`
		Incoming []string `xml:"incoming"`
		Outgoing []string `xml:"outgoing"`
	} `xml:"endEvent"`
}

func (d *ProcessRequest) ToEvents() []model.Event {
	var events []model.Event
	for _, v := range d.StartEvent {
		event := model.Event{
			DbBase: model.DbBase{
				Id:         v.ID,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			},
			Name:     v.Name,
			Incoming: v.Incoming,
			Outgoing: v.Outgoing,
			Type:     constant.StartEvent,
		}
		events = append(events, event)
	}

	for _, v := range d.EndEvent {
		event := model.Event{
			DbBase: model.DbBase{
				Id:         v.ID,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			},
			Name:     v.Name,
			Incoming: v.Incoming,
			Outgoing: v.Outgoing,
			Type:     constant.EndEvent,
		}
		events = append(events, event)
	}

	return events
}

func (d *ProcessRequest) ToProcess(originXml string) model.Process {
	return model.Process{
		DbBase: model.DbBase{
			Id:         d.ID,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		},
		Name:                d.Name,
		Category:            d.ProcessCategory,
		Version:             1, // 默认是版本1
		Resource:            originXml,
		StartEventIds:       nil,
		SequenceFlowIds:     nil,
		UserTaskIds:         nil,
		ExclusiveGatewayIds: nil,
		EndEventIds:         nil,
	}
}
