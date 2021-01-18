/**
 * @Author: lzw5399
 * @Date: 2021/1/16 21:59
 * @Desc: 流程实例 process instance
 */
package request

import (
	"workflow/model"
	"workflow/util"
)

type InstanceRequest struct {
	ProcessCode string                   `json:"processCode"` // 流程模板唯一标识code
	Variables   []model.InstanceVariable `json:"variables"`   // 流程启动时的初始变量
}

func (i *InstanceRequest) ProcessInstance(processId uint) model.ProcessInstance {
	return model.ProcessInstance{
		ProcessId:  processId,
		Variables:  util.StructToBytes(i.Variables),
		IsFinished: false,
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
