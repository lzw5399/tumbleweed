/**
 * @Author: lzw5399
 * @Date: 2021/3/20 14:59
 * @Desc:
 */
package model

import "time"

type Tenant struct {
	EntityBase
	Name       string    `json:"name"`
	CreateTime time.Time `gorm:"default:now();type:timestamp" json:"createTime" form:"createTime"`
}
