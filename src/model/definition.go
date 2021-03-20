/**
 * @Author: lzw5399
 * @Date: 2021/1/14 22:13
 * @Desc: 流程的定义
 */
package model

import (
	"gorm.io/datatypes"
)

// 流程定义表
type ProcessDefinition struct {
	AuditableBase
	Name        string         `gorm:"column:name; type:varchar(128)" json:"name" form:"name"`                             // 流程名称
	FormId      int            `json:"formId" form:"formId"`                                                               // 对应的表单的id(表单不存在于当前系统中，仅对外部系统做一个标记)
	Structure   datatypes.JSON `gorm:"column:structure; type:jsonb" json:"structure" form:"structure"`                     // 流程的具体结构
	ClassifyId  int            `gorm:"column:classify_id; type:integer" json:"classifyId" form:"classifyId"`               // 分类ID
	Task        datatypes.JSON `gorm:"column:task; type:jsonb" jsonb:"task" form:"task"`                                   // 任务ID, array, 可执行多个任务，可以当成通知任务，每个节点都会去执行
	SubmitCount int            `gorm:"column:submit_count; type:integer; default:0" json:"submitCount" form:"submitCount"` // 提交统计
	Notice      datatypes.JSON `gorm:"column:notice; type:jsonb" json:"notice" form:"notice"`                              // 绑定通知
	TenantId    int            `gorm:"index" json:"tenantId" form:"tenantId"`                                              // 租户id
	Remarks     string         `gorm:"column:remarks; type:text" json:"remarks" form:"remarks"`                            // 流程备注
}
