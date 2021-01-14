/**
 * @Author: lzw5399
 * @Date: 2021/1/14 21:38
 * @Desc: 流程实例
 */
package model

import (
	"gorm.io/gorm"
)

// 流程实例
type ProcessInstance struct {
	gorm.Model
	ProcDefID     int    `json:"procDefId"`   // 流程定义ID
	ProcDefName   string `json:"procDefName"` // 流程定义名
	Title         string `json:"title"`       // title 标题
	Department    string `json:"department"`  // 用户部门
	Company       string `json:"company"`     // 公司
	NodeID        string `json:"nodeID"`      // 当前节点
	Candidate     string `json:"candidate"`   // 审批人
	TaskID        int    `json:"taskID"`      // 当前任务
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	Duration      int64  `json:"duration"`
	StartUserID   string `json:"startUserId"`
	StartUserName string `json:"startUserName"`
	IsFinished    bool   `gorm:"default:false" json:"isFinished"`
}
