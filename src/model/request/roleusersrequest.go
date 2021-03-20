/**
 * @Author: lzw5399
 * @Date: 2021/3/20 22:03
 * @Desc:
 */
package request

import (
	"time"

	"workflow/src/model"
)

type BatchSyncRoleUsersRequest struct {
	RoleUsersList []SyncRoleUsersRequest
}

type SyncRoleUsersRequest struct {
	RoleId  string   `json:"roleId"`
	UserIds []string `json:"userIds"`
}

func (s *SyncRoleUsersRequest) ToRoleUsers(tenantId uint) model.RoleUsers{
	return model.RoleUsers{
		RoleId:     s.RoleId,
		UserIds:    s.UserIds,
		TenantId:   int(tenantId),
		CreateTime: time.Now().Local(),
	}
}
