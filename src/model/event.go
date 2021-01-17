/**
 * @Author: lzw5399
 * @Date: 2021/1/16 11:20
 * @Desc:
 */
package model

import "github.com/lib/pq"

type Event struct {
	EntityBase
	Code      string         `json:"code" gorm:"uniqueIndex"`
	Name      string         `json:"name"`
	Incoming  pq.StringArray `json:"incoming" gorm:"type:text[];default:array[]::text[]"`
	Outgoing  pq.StringArray `json:"outgoing" gorm:"type:text[];default:array[]::text[]"`
	Type      int            `json:"type" gorm:"index:idx_type"` // constant.StartEvent 或 constant.EndEvent
	ProcessId uint           `json:"processId" gorm:"index:idx_processId1"`
}
