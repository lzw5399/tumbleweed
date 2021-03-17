/**
 * @Author: lzw5399
 * @Date: 2021/1/16 21:59
 * @Desc: 流程实例 process instance
 */
package request

import (
	"time"

	"workflow/src/model"
)

type ProcessInstanceRequest struct {
	Title               string `json:"title" form:"title"`                             // 流程实例标题
	ProcessDefinitionId int    `json:"processDefinitionId" form:"processDefinitionId"` // 流程ID
}

func (i *ProcessInstanceRequest) ToProcessInstance(currentUserId uint) model.ProcessInstance {
	return model.ProcessInstance{
		AuditableBase: model.AuditableBase{
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
			CreateBy:   currentUserId,
			UpdateBy:   currentUserId,
		},
		Title:               i.Title,
		ProcessDefinitionId: i.ProcessDefinitionId,
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
