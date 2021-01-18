/**
 * @Author: lzw5399
 * @Date: 2021/1/18 21:05
 * @Desc: candidateGroups，可以认为是【角色role】
 */
package model

import "gorm.io/gorm"

// 用户组
type Group struct {
	gorm.Model
	EntityBase
	Identifier string `json:"identifier"` // 外部系统的标识
	Name       string `json:"name"`
	//TenantId   string `json:"tenantId"`
}
