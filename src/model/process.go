/**
 * @Author: lzw5399
 * @Date: 2021/1/14 22:13
 * @Desc: 流程的定义
 */
package model

import "github.com/lib/pq"

// 流程定义表
type Process struct {
	DbBase
	Name                string         `json:"name"`                                                           // 流程名字
	Category            string         `json:"category"`                                                       // 流程类别
	Version             int            `json:"version,omitempty"`                                              // 版本
	Resource            string         `json:"resource"`                                                       // 流程定义的完整bpmn xml字符串
	StartEventIds       pq.StringArray `json:"startEventIds" gorm:"type:text[] default array[]::text[]"`       // 开始事件的ids
	SequenceFlowIds     pq.StringArray `json:"sequenceFlowIds" gorm:"type:text[] default array[]::text[]"`     // 顺序流的ids
	UserTaskIds         pq.StringArray `json:"userTaskIds" gorm:"type:text[] default array[]::text[]"`         // 用户任务的ids
	ExclusiveGatewayIds pq.StringArray `json:"exclusiveGatewayIds" gorm:"type:text[] default array[]::text[]"` // 排他网关的ids
	EndEventIds         pq.StringArray `json:"endEventIds" gorm:"type:text[] default array[]::text[]"`         // 结束事件的ids
}
