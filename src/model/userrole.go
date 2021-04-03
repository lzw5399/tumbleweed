/**
 * @Author: lzw5399
 * @Date: 2021/3/20 21:51
 * @Desc: 外部系统的角色用户表
 */
package model

import (
	"time"
)

// 外部系统的用户表
type User struct {
	EntityBase
	Identifier string    `gorm:"index" json:"identifier"` // 外部系统用户id
	Name       string    `json:"name"`                    // 用户名称
	TenantId   int       `gorm:"index" json:"tenantId"`   // 租户id
	CreateTime time.Time `gorm:"default:now();type:timestamp" json:"createTime"`
}

// 外部系统的角色表
type Role struct {
	EntityBase
	Identifier string    `gorm:"index" json:"identifier"` // 外部系统角色id
	Name       string    `json:"name"`                    // 角色名字
	TenantId   int       `gorm:"index" json:"tenantId"`   // 租户id
	CreateTime time.Time `gorm:"default:now();type:timestamp" json:"createTime"`
}

// 用户角色关联表
type UserRole struct {
	EntityBase
	UserIdentifier string `gorm:"index" json:"userIdentifier"`
	RoleIdentifier string `gorm:"index" json:"roleIdentifier"`
}

//// 角色用户: 一条数据是一个角色
//type RoleUsers struct {
//	EntityBase
//	RoleId     int           `gorm:"index" json:"roleId" form:"roleId"`                                        // 外部系统的角色Id
//	UserIds    pq.Int64Array `gorm:"type:integer[]; default:array[]::integer[]" json:"userIds" form:"userIds"` // 外部系统的用户Id数组
//	TenantId   int           `gorm:"index" json:"tenantId" form:"tenantId"`                                    // 租户id
//	CreateTime time.Time     `gorm:"default:now();type:timestamp" json:"createTime" form:"createTime"`
//}
