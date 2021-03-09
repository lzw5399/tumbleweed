/**
 * @Author: lzw5399
 * @Date: 2021/1/14 23:11
 * @Desc:
 */
package model

import "time"

type EntityBase struct {
	Id uint `gorm:"primarykey" json:"id" form:"id"`
}

type AuditableBase struct {
	EntityBase
	CreateTime time.Time `gorm:"default:now()" json:"create_time" form:"create_time"`
	UpdateTime time.Time `gorm:"default:now()"  json:"update_time" form:"update_time"`
	CreateBy   uint      `json:"create_by" form:"create_by"`
	UpdateBy   uint      `json:"update_by" form:"update_by"`
}
