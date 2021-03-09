/**
 * @Author: lzw5399
 * @Date: 2021/1/16 21:59
 * @Desc: 流程实例 process instance
 */
package request

import (
	"encoding/json"
	"workflow/src/model"
)

type ProcessInstanceRequest struct {
	model.ProcessInstance
	FormId      int             `json:"formId" form:"formId"` // 表单的标识
	SourceState string          `json:"sourceState"`
	Tasks       json.RawMessage `json:"tasks"`
	Source      string          `json:"source"`
}

//func (i *ProcessInstanceRequest) ProcessInstance(processId uint) model.ProcessInstance {
//	return model.ProcessInstance{
//		ProcessId:  processId,
//		Variables:  util.StructToBytes(i.Variables),
//		IsFinished: false,
//	}
//}

type GetVariableRequest struct {
	InstanceId   uint   `json:"instanceId,omitempty" form:"instanceId,omitempty"`
	VariableName string `json:"variableName,omitempty" form:"variableName,omitempty"`
}

type GetVariableListRequest struct {
	PagingRequest
	InstanceId uint `json:"instanceId,omitempty" form:"instanceId,omitempty"`
}
