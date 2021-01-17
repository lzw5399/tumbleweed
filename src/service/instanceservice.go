/**
 * @Author: lzw5399
 * @Date: 2021/1/16 22:58
 * @Desc: 流程实例服务
 */
package service

import (
	"errors"
	"net/http"

	"workflow/src/global"
	"workflow/src/model"
	"workflow/src/model/request"

	"gorm.io/gorm"
)

func StartProcessInstance(r *request.InstanceRequest) (uint, error) {
	// 检查流程是否存在
	var process model.Process
	err := global.BankDb.Where("code=?", r.ProcessCode).First(&process).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusBadRequest, errors.New("当前流程不存在，请检查后重试")
	}

	// 创建流程
	instance := r.ProcessInstance(process.Id)
	err = global.BankDb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&instance).Error; err != nil {
			return err
		}

		// 返回nil提交事务
		return nil
	})

	if err != nil {
		return http.StatusInternalServerError, errors.New("流程实例创建失败")
	}

	return instance.Id, nil
}
