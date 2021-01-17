/**
 * @Author: lzw5399
 * @Date: 2021/1/14 23:14
 * @Desc: 流程任务表
 */
package model

import "github.com/lib/pq"

// 流程任务
type UserTask struct {
	EntityBase
	Code            string         `json:"code" gorm:"uniqueIndex"`
	Name            string         `json:"name"`
	FormKey         string         `json:"formKey"`
	Assignee        string         `json:"assignee"`                                                   // 指定人员
	CandidateUsers  pq.StringArray `json:"candidateUsers" gorm:"type:text[];default:array[]::text[]"`  // 候选人, 逗号分割
	CandidateGroups pq.StringArray `json:"candidateGroups" gorm:"type:text[];default:array[]::text[]"` // 候选组， 逗号分割
	Incoming        pq.StringArray `json:"incoming" gorm:"type:text[];default:array[]::text[]"`
	Outgoing        pq.StringArray `json:"outgoing" gorm:"type:text[];default:array[]::text[]"`
	ProcessId       uint           `json:"processId" gorm:"index:idx_processId5"`
}
