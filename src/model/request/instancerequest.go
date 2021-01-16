/**
 * @Author: lzw5399
 * @Date: 2021/1/16 21:59
 * @Desc: 流程实例 process instance
 */
package request

import (
	"workflow/src/model"
	"workflow/src/util"
)

type InstanceRequest struct {
	ProcessCode string                 `json:"processCode"` // 流程模板唯一标识code
	Variables   map[string]interface{} `json:"variables"`   // 流程启动时的初始变量
}

func (i *InstanceRequest) ProcessInstance(processId uint) model.ProcessInstance {
	return model.ProcessInstance{
		ProcessId:  processId,
		Variables:  util.MapToBytes(i.Variables),
		IsFinished: false,
	}
}
