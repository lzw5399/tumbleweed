/**
 * @Author: lzw5399
 * @Date: 2021/3/20 22:03
 * @Desc:
 */
package request

import (
	"time"

	"workflow/src/model"
	"workflow/src/util"
)

type BatchSyncRoleUsersRequest struct {
	RoleUsersList []SyncRoleUsersRequest
}

type SyncRoleUsersRequest struct {
	RoleId  int   `json:"roleId"`
	UserIds []int `json:"userIds"`
}

func (s *SyncRoleUsersRequest) ToRoleUsers(tenantId uint) model.RoleUsers {
	return model.RoleUsers{
		RoleId:     s.RoleId,
		UserIds:    util.ParseToInt64Array(s.UserIds),
		TenantId:   int(tenantId),
		CreateTime: time.Now().Local(),
	}
}
