/**
 * @Author: lzw5399
 * @Date: 2021/3/20 22:08
 * @Desc: 外部系统的角色用户对应关系 服务
 */
package service

import (
	"log"

	. "github.com/ahmetb/go-linq/v3"

	"workflow/src/global"
	"workflow/src/model"
	"workflow/src/model/request"
)

type RoleUsersService interface {
	BatchSyncRoleUsers(*request.BatchSyncRoleUsersRequest, uint) error
	SyncRoleUsers(*request.SyncRoleUsersRequest, uint) error
}

func NewRoleUsersService() *roleUsersService {
	return &roleUsersService{}
}

type roleUsersService struct {
}

// 异步批量同步外部系统的角色用户对应关系
func (u *roleUsersService) BatchSyncRoleUsers(r *request.BatchSyncRoleUsersRequest, tenantId uint) error {
	go u.BatchSyncRoleUsersAsync(r, tenantId)

	return nil
}

// 批量同步外部系统的角色用户对应关系
func (u *roleUsersService) BatchSyncRoleUsersAsync(r *request.BatchSyncRoleUsersRequest, tenantId uint) {
	var roleIds []int
	From(r.RoleUsersList).Select(func(i interface{}) interface{} {
		return i.(request.SyncRoleUsersRequest).RoleId
	}).ToSlice(&roleIds)

	// 删除当前租户下当前请求的所有用户信息
	err := global.BankDb.
		Where("tenant_id = ?", tenantId).
		Where("role_id in ?", roleIds).
		Delete(model.RoleUsers{}).
		Error
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 映射成数据库实体
	var roleUsersList []model.RoleUsers
	From(r.RoleUsersList).Select(func(i interface{}) interface{} {
		re := i.(request.SyncRoleUsersRequest)
		return re.ToRoleUsers(tenantId)
	}).ToSlice(&roleUsersList)

	// 批量创建数据
	err = global.BankDb.
		Model(&model.RoleUsers{}).
		Create(&roleUsersList).
		Error
	if err != nil {
		log.Printf("批量更新失败, 错误：%s", err.Error())
	}
	log.Printf("批量更新角色用户对应关系成功，更新条数:%d\n", len(roleUsersList))
}

// 同步外部系统的角色用户对应关系
func (u *roleUsersService) SyncRoleUsers(r *request.SyncRoleUsersRequest, tenantId uint) error {
	// 删除当前租户下当前请求的角色用户
	err := global.BankDb.
		Where("tenant_id = ?", tenantId).
		Where("role_id = ?", r.RoleId).
		Delete(model.RoleUsers{}).
		Error
	if err != nil {
		log.Printf("同步失败, 错误：%s", err.Error())
		return err
	}

	roleUsers := r.ToRoleUsers(tenantId)
	err = global.BankDb.
		Model(&model.RoleUsers{}).
		Create(&roleUsers).
		Error
	if err != nil {
		log.Printf("同步失败, 错误：%s", err.Error())
		return err
	}

	return nil
}
