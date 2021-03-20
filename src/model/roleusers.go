/**
 * @Author: lzw5399
 * @Date: 2021/3/20 21:51
 * @Desc: 外部系统的角色用户表
 */
package model

import (
	"time"

	"github.com/lib/pq"
)

// 角色用户: 一条数据是一个角色
type RoleUsers struct {
	EntityBase
	RoleId     string         `gorm:"index" json:"roleId" form:"roleId"`                                  // 外部系统的角色Id
	UserIds    pq.StringArray `gorm:"type:text[]; default:array[]::text[]" json:"userIds" form:"userIds"` // 外部系统的用户Id数组
	TenantId   int            `gorm:"index" json:"tenantId" form:"tenantId"`                              // 租户id
	CreateTime time.Time      `gorm:"default:now();type:timestamp" json:"createTime" form:"createTime"`
}
