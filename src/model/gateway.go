/**
 * @Author: lzw5399
 * @Date: 2021/1/16 13:15
 * @Desc:
 */
package model

import "github.com/lib/pq"

type ExclusiveGateway struct {
	DbBase
	Incoming pq.StringArray `json:"incoming" gorm:"type:text[] default:array[]::text[]"`
	Outgoing pq.StringArray `json:"outgoing" gorm:"type:text[] default:array[]::text[]"`
}
