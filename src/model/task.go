/**
 * @Author: lzw5399
 * @Date: 2021/1/14 23:14
 * @Desc: 流程任务表
 */
package model

// 流程任务
type Task struct {
	DbBase
	NodeID        string `json:"nodeId"` // 当前执行流所在的节点
	Step          int    `json:"step"`
	ProcInstID    int    `json:"procInstID"` // 流程实例id
	Assignee      string `json:"assignee"`
	ClaimTime     string `json:"claimTime"`
	MemberCount   int8   `json:"memberCount" gorm:"default:1"` // 还未审批的用户数，等于0代表会签已经全部审批结束，默认值为1
	UnCompleteNum int8   `json:"unCompleteNum" gorm:"default:1"`
	AgreeNum      int8   `json:"agreeNum"`                    //审批通过数
	ActType       string `json:"actType" gorm:"default:'or'"` // and 为会签，or为或签，默认为or
	IsFinished    bool   `gorm:"default:false" json:"isFinished"`
}
