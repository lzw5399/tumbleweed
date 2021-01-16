/**
 * @Author: lzw5399
 * @Date: 2021/1/14 23:14
 * @Desc: 流程任务表
 */
package model

// 流程任务
type UserTask struct {
	DbBase
	Name            string   `json:"name"`
	FormKey         string   `json:"formKey"`
	CandidateUsers  string   `json:"candidateUsers"`  // 候选人, 逗号分割
	CandidateGroups string   `json:"candidateGroups"` // 候选组， 逗号分割
	Assignee        string   `json:"assignee"`        // 指定人员
	Incoming        []string `json:"incoming" gorm:"type:text[] default:array[]::text[]"`
	Outgoing        []string `json:"outgoing" gorm:"type:text[] default:array[]::text[]"`
}
