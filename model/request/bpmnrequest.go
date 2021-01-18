/**
 * @Author: lzw5399
 * @Date: 2021/1/15 21:24
 * @Desc: 接收bpmn2.0的struct
 */
package request

import (
	"strings"
	"workflow/global/constant"
	"workflow/model"
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

func (p *ProcessRequest) Events(processId uint) []model.Event {
	var events []model.Event
	for _, v := range p.StartEvent {
		event := model.Event{
			Code:      v.ID, // 导入的id对应code
			Name:      v.Name,
			Incoming:  v.Incoming,
			Outgoing:  v.Outgoing,
			Type:      constant.StartEvent,
			ProcessId: processId,
		}
		events = append(events, event)
	}

	for _, v := range p.EndEvent {
		event := model.Event{
			Code:     v.ID, // 导入的id对应code
			Name:     v.Name,
			Incoming: v.Incoming,
			Outgoing: v.Outgoing,
			Type:     constant.EndEvent,
		}
		events = append(events, event)
	}

	return events
}

func (p *ProcessRequest) SequenceFlows(processId uint) []model.SequenceFlow {
	var flows []model.SequenceFlow
	for _, v := range p.SequenceFlow {
		flow := model.SequenceFlow{
			Code:                v.ID,
			SourceRef:           v.SourceRef,
			TargetRef:           v.TargetRef,
			ConditionExpression: v.ConditionExpression.Text,
			ProcessId:           processId,
		}
		flows = append(flows, flow)
	}

	return flows
}

func (p *ProcessRequest) Tasks(processId uint) []model.UserTask {
	var tasks []model.UserTask
	for _, v := range p.UserTask {
		userTask := model.UserTask{
			Code:            v.ID,
			Name:            v.Name,
			FormKey:         v.FormKey,
			Assignee:        v.Assignee,
			CandidateUsers:  nil,
			CandidateGroups: nil,
			Incoming:        v.Incoming,
			Outgoing:        v.Outgoing,
			ProcessId:       processId,
		}
		if v.CandidateGroups != "" {
			userTask.CandidateGroups = strings.Split(v.CandidateGroups, ",")
		}
		if v.CandidateUsers != "" {
			userTask.CandidateUsers = strings.Split(v.CandidateUsers, ",")
		}
		tasks = append(tasks, userTask)
	}

	return tasks
}

func (p *ProcessRequest) ExclusiveGateways(processId uint) []model.ExclusiveGateway {
	var gateways []model.ExclusiveGateway
	for _, v := range p.ExclusiveGateway {
		gateway := model.ExclusiveGateway{
			Code:      v.ID,
			Incoming:  v.Incoming,
			Outgoing:  v.Outgoing,
			ProcessId: processId,
		}
		gateways = append(gateways, gateway)
	}

	return gateways
}

func (p *ProcessRequest) Process(originXml string) model.Process {
	return model.Process{
		Code:     p.ID, // 导入的id对应code
		Name:     p.Name,
		Category: p.ProcessCategory,
		Version:  1,         // 默认是版本1
		Resource: originXml, // 原始的xml存档
	}
}
