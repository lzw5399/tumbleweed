/**
 * @Author: lzw5399
 * @Date: 2021/3/20 22:03
 * @Desc:
 */
package request

import (
	"github.com/ahmetb/go-linq/v3"

	"workflow/src/model"
)

// 同步
type BatchSyncUserRoleRequest struct {
	Users     []UserRequest     `json:"users"`
	Roles     []RoleRequest     `json:"roles"`
	UserRoles []UserRoleRequest `json:"userRoles"`
}

type UserRequest struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

type RoleRequest struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

type UserRoleRequest struct {
	UserIdentifier string `json:"userIdentifier"`
	RoleIdentifier string `json:"roleIdentifier"`
}

func (s *BatchSyncUserRoleRequest) ToDbEntities(tenantId int) (users []model.User, roles []model.Role, userRole []model.UserRole) {
	if s.Users != nil {
		linq.From(s.Users).SelectT(func(i UserRequest) interface{} {
			return model.User{
				Identifier: i.Identifier,
				Name:       i.Name,
				TenantId:   tenantId,
			}
		}).ToSlice(&users)
	}

	if s.Roles != nil {
		linq.From(s.Roles).SelectT(func(i RoleRequest) interface{} {
			return model.Role{
				Identifier: i.Identifier,
				Name:       i.Name,
				TenantId:   tenantId,
			}
		}).ToSlice(&roles)
	}

	if s.UserRoles != nil {
		linq.From(s.UserRoles).SelectT(func(i UserRoleRequest) interface{} {
			return model.UserRole{
				UserIdentifier: i.UserIdentifier,
				RoleIdentifier: i.RoleIdentifier,
			}
		}).ToSlice(&userRole)
	}

	return
}
