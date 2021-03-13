/**
 * @Author: lzw5399
 * @Date: 2021/1/16 21:59
 * @Desc: 流程实例 process instance
 */
package request

import (
	"encoding/json"
	"time"
	"workflow/src/model"
)

type ProcessInstanceRequest struct {
	model.ProcessInstance
	SourceState string          `json:"sourceState"`
	Tasks       json.RawMessage `json:"tasks"`
	Source      string          `json:"source"`
}

func (i *ProcessInstanceRequest) ToProcessInstance(currentUserId uint) model.ProcessInstance {
	return model.ProcessInstance{
		AuditableBase: model.AuditableBase{
			EntityBase: model.EntityBase{
				Id: i.Id,
			},
			CreateTime: time.Now(),
			CreateBy:   currentUserId,
		},
		Title:               i.Title,
		Priority:            i.Priority,
		ProcessDefinitionId: i.ProcessDefinitionId,
		Classify:            i.Classify,
		IsEnd:               i.IsEnd,
		IsDenied:            i.IsDenied,
		State:               i.State,
		RelatedPerson:       i.RelatedPerson,
		UrgeCount:           i.UrgeCount,
		UrgeLastTime:        i.UrgeLastTime,
	}
}

type GetVariableRequest struct {
	InstanceId   uint   `json:"instanceId,omitempty" form:"instanceId,omitempty"`
	VariableName string `json:"variableName,omitempty" form:"variableName,omitempty"`
}

type GetVariableListRequest struct {
	PagingRequest
	InstanceId uint `json:"instanceId,omitempty" form:"instanceId,omitempty"`
}
