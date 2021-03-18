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
	CreateTime time.Time `gorm:"default:now()" json:"createTime" form:"createTime"`
	UpdateTime time.Time `gorm:"default:now()"  json:"updateTime" form:"updateTime"`
	CreateBy   uint      `json:"createBy" form:"createBy"`
	UpdateBy   uint      `json:"updateBy" form:"updateBy"`
}
