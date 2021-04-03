/**
 * @Author: lzw5399
 * @Date: 2021/1/14 23:11
 * @Desc:
 */
package model

import "time"

type EntityBase struct {
	Id int `gorm:"primarykey" json:"id" form:"id"`
}

type AuditableBase struct {
	EntityBase
	CreateTime time.Time `gorm:"default:now();type:timestamp" json:"createTime" form:"createTime"`
	UpdateTime time.Time `gorm:"default:now();type:timestamp"  json:"updateTime" form:"updateTime"`
	CreateBy   string    `json:"createBy" form:"createBy"`
	UpdateBy   string    `json:"updateBy" form:"updateBy"`
}
