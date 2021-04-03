/**
 * @Author: lzw5399
 * @Date: 2021/3/20 22:08
 * @Desc: 外部系统的角色用户对应关系 服务
 */
package service

import (
	. "github.com/ahmetb/go-linq/v3"

	"workflow/src/global"
	"workflow/src/model"
	"workflow/src/model/request"
)

// 异步批量同步外部系统的角色用户对应关系
func BatchSyncRoleUsers(r *request.BatchSyncUserRoleRequest, tenantId int) error {
	go BatchSyncRoleUsersAsync(r, tenantId)

	return nil
}

func BatchSyncRoleUsersAsync(r *request.BatchSyncUserRoleRequest, tenantId int) {
	users, roles, userRoles := r.ToDbEntities(tenantId)

	// 同步user
	err := DeleteOriginUsers(users, tenantId)
	if err != nil {
		global.BankLogger.Error("删除用户信息失败", err)
	}
	err = DeleteOriginRoles(roles, tenantId)
	if err != nil {
		global.BankLogger.Error("删除角色信息失败", err)
	}
	err = DeleteOriginUserRoles(userRoles, tenantId)
	if err != nil {
		global.BankLogger.Error("删除角色用户关联信息失败", err)
	}

	// 批量创建数据
	err = global.BankDb.
		Model(&model.User{}).
		Create(&users).
		Error
	if err != nil {
		global.BankLogger.Error("批量更新用户失败", err)
	}
	err = global.BankDb.
		Model(&model.Role{}).
		Create(&roles).
		Error
	if err != nil {
		global.BankLogger.Error("批量更新角色失败", err)
	}
	err = global.BankDb.
		Model(&model.UserRole{}).
		Create(&userRoles).
		Error
	if err != nil {
		global.BankLogger.Error("批量更新用户角色关联关系失败", err)
	}

	global.BankLogger.Infof("批量更新角色用户对应关系成功，更新用户条数:%d，更新角色条数:%d，更新关联关系条数:%d\n", len(users), len(roles), len(userRoles))
}

func DeleteOriginUsers(users []model.User, tenantId int) error {
	var userIds []string
	From(users).Select(func(i interface{}) interface{} {
		return i.(model.User).Identifier
	}).ToSlice(&userIds)

	err := global.BankDb.Model(&model.User{}).
		Where("identifier in ?", userIds).
		Where("tenant_id = ?", tenantId).
		Delete(model.User{}).
		Error

	return err
}

func DeleteOriginRoles(roles []model.Role, tenantId int) error {
	var roleIds []string
	From(roles).Select(func(i interface{}) interface{} {
		return i.(model.Role).Identifier
	}).ToSlice(&roleIds)

	err := global.BankDb.Model(&model.Role{}).
		Where("identifier in ?", roleIds).
		Where("tenant_id = ?", tenantId).
		Delete(model.User{}).
		Error

	return err
}

func DeleteOriginUserRoles(userRoles []model.UserRole, tenantId int) error {
	var userIds []string
	var roleIds []string
	for _, userRole := range userRoles {
		userIds = append(userIds, userRole.UserIdentifier)
		roleIds = append(roleIds, userRole.RoleIdentifier)
	}

	err := global.BankDb.Model(&model.UserRole{}).
		Where("tenant_id = ?", tenantId).
		Where("user_identifier in ?", userIds).
		Or("role_identifier in ?", roleIds).
		Delete(model.UserRole{}).
		Error

	return err
}

//
//// 批量同步外部系统的角色用户对应关系
//func BatchSyncRoleUsersAsync(r *request.BatchSyncUserRoleRequest, tenantId int) {
//	var roleIds []int
//	From(r.RoleUsersList).Select(func(i interface{}) interface{} {
//		return i.(request.SyncRoleUsersRequest).RoleId
//	}).ToSlice(&roleIds)
//
//	// 删除当前租户下当前请求的所有用户信息
//	err := global.BankDb.
//		Where("tenant_id = ?", tenantId).
//		Where("role_id in ?", roleIds).
//		Delete(model.RoleUsers{}).
//		Error
//	if err != nil {
//		log.Println(err.Error())
//		return
//	}
//
//	// 映射成数据库实体
//	var roleUsersList []model.RoleUsers
//	From(r.RoleUsersList).Select(func(i interface{}) interface{} {
//		re := i.(request.SyncRoleUsersRequest)
//		return re.ToRoleUsers(tenantId)
//	}).ToSlice(&roleUsersList)
//
//	// 批量创建数据
//	err = global.BankDb.
//		Model(&model.RoleUsers{}).
//		Create(&roleUsersList).
//		Error
//	if err != nil {
//		global.BankLogger.Error("批量更新用户角色关系失败", err)
//	}
//	global.BankLogger.Infof("批量更新角色用户对应关系成功，更新条数:%d\n", len(roleUsersList))
//}

//
//// 同步外部系统的角色用户对应关系
//func SyncRoleUsers(r *request.SyncRoleUsersRequest, tenantId int) error {
//	// 删除当前租户下当前请求的角色用户
//	err := global.BankDb.
//		Where("tenant_id = ?", tenantId).
//		Where("role_id = ?", r.RoleId).
//		Delete(model.RoleUsers{}).
//		Error
//	if err != nil {
//		global.BankLogger.Error("同步用户角色关系失败", err)
//		return err
//	}
//
//	roleUsers := r.ToRoleUsers(tenantId)
//	err = global.BankDb.
//		Model(&model.RoleUsers{}).
//		Create(&roleUsers).
//		Error
//	if err != nil {
//		global.BankLogger.Error("同步用户角色关系失败", err)
//		return err
//	}
//
//	global.BankLogger.Info("同步用户角色关系成功", r)
//
//	return nil
//}
