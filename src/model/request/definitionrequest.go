/**
 * @Author: lzw5399
 * @Date: 2021/3/5 11:05
 * @Desc:
 */
package request

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"

	"workflow/src/model"
	"workflow/src/model/dto"
)

type ProcessDefinitionRequest struct {
	Id         int             `json:"id" form:"id"`
	Name       string          `json:"name" form:"name"`                                // 流程名称
	FormId     int             `json:"formId" form:"formId"`                            // 对应的表单的id(仅对外部系统做一个标记)
	Structure  json.RawMessage `json:"structure" form:"structure" swaggertype:"string"` // 流程结构
	ClassifyId int             `json:"classifyId" form:"classifyId"`                    // 分类ID
	Task       json.RawMessage `json:"task" form:"task" swaggertype:"string"`           // 任务ID, array, 可执行多个任务，可以当成通知任务，每个节点都会去执行
	Notice     json.RawMessage `json:"notice" form:"notice" swaggertype:"string"`       // 绑定通知
	Remarks    string          `json:"remarks" form:"remarks"`                          // 流程备注
}

func (p *ProcessDefinitionRequest) ProcessDefinition() model.ProcessDefinition {
	var structure dto.Structure
	_ = json.Unmarshal(p.Structure, &structure)

	return model.ProcessDefinition{
		AuditableBase: model.AuditableBase{
			EntityBase: model.EntityBase{
				Id: p.Id,
			},
			CreateTime: time.Now().Local(),
			UpdateTime: time.Now().Local(),
		},
		Name:        p.Name,
		Structure:   structure,
		ClassifyId:  p.ClassifyId,
		Task:        datatypes.JSON(p.Task),
		Notice:      datatypes.JSON(p.Notice),
		Remarks:     p.Remarks,
		FormId:      p.FormId,
		SubmitCount: 0,
	}
}
