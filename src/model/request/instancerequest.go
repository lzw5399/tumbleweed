/**
 * @Author: lzw5399
 * @Date: 2021/1/16 21:59
 * @Desc: 流程实例 process instance
 */
package request

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"

	"workflow/src/model"
)

type ProcessInstanceRequest struct {
	Title               string          `json:"title" form:"title"`                                      // 工单标题
	Priority            int             `json:"priority" form:"priority"`                                // 工单优先级 1，正常 2，紧急 3，非常紧急
	ProcessDefinitionId int             `json:"processDefinitionId" form:"processDefinitionId"`          // 流程ID
	Classify            int             `json:"classify" form:"classify"`                                // 分类ID
	State               json.RawMessage `json:"state" form:"state" swaggertype:"string"`                 // 状态信息
	RelatedPerson       json.RawMessage `json:"relatedPerson" form:"relatedPerson" swaggertype:"string"` // 工单所有处理人
	SourceState         string          `json:"sourceState"`
	Tasks               json.RawMessage `json:"tasks" swaggertype:"string"`
	Source              string          `json:"source"`
}

func (i *ProcessInstanceRequest) ToProcessInstance(currentUserId uint) model.ProcessInstance {
	return model.ProcessInstance{
		AuditableBase: model.AuditableBase{
			CreateTime: time.Now(),
			CreateBy:   currentUserId,
		},
		Title:               i.Title,
		Priority:            i.Priority,
		ProcessDefinitionId: i.ProcessDefinitionId,
		Classify:            i.Classify,
		State:               datatypes.JSON(i.State),
		RelatedPerson:       datatypes.JSON(i.RelatedPerson),
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
