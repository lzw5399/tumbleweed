/**
 * @Author: lzw5399
 * @Date: 2021/1/15 23:35
 * @Desc:
 */
package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"

	"workflow/src/global"
	"workflow/src/model"
	"workflow/src/model/request"
)

type DefinitionService interface {
	CreateDefinition(*request.ProcessDefinitionRequest) (*model.ProcessDefinition, error)
	Validate(*request.ProcessDefinitionRequest, uint) error
	UpdateDefinition(r *request.ProcessDefinitionRequest) error
	DeleteDefinition(id uint) error
	GetDefinition(id uint) (*model.ProcessDefinition, error)
}

func NewDefinitionService() *definitionService {
	return &definitionService{}
}

type definitionService struct {
}

func (d *definitionService) GetDefinition(id uint) (*model.ProcessDefinition, error) {
	var definition model.ProcessDefinition

	err := global.BankDb.
		Where("id=?", id).
		Find(&definition).Error
	if err != nil {
		log.Error(err)
		return nil, errors.New("查询流程详情失败")
	}

	return &definition, nil
}

// 验证
func (d *definitionService) Validate(r *request.ProcessDefinitionRequest, excludeId uint) error {
	// 验证名称是否已存在
	var c int64
	global.BankDb.Model(&model.ProcessDefinition{}).
		Where("name=?", r.Name).
		Where("id!=?", excludeId).
		Count(&c)
	if c != 0 {
		return errors.New(fmt.Sprintf("当前名称为:\"%s\"的模板已存在", r.Name))
	}

	// todo 校验structure的json

	return nil
}

// 创建新的process流程
func (d *definitionService) CreateDefinition(r *request.ProcessDefinitionRequest) (*model.ProcessDefinition, error) {
	var (
		err error
	)

	processDefinition := r.ProcessDefinition()

	if err = global.BankDb.Create(&processDefinition).Error; err != nil {
		log.Error(err)
		return nil, err
	}

	return &processDefinition, nil
}

// 更新流程定义
func (d *definitionService) UpdateDefinition(r *request.ProcessDefinitionRequest) error {
	processDefinition := r.ProcessDefinition()

	err := global.BankDb.
		Where("id = ?", processDefinition.Id).
		Updates(map[string]interface{}{
			"name":        processDefinition.Name,
			"form_id":     processDefinition.FormId,
			"structure":   processDefinition.Structure,
			"classify_id": processDefinition.ClassifyId,
			"task":        processDefinition.Task,
			"notice":      processDefinition.Notice,
			"remarks":     processDefinition.Remarks,
			"update_by":   1, //todo currentid
			"update_time": time.Now(),
		}).Error

	return err
}

// 删除流程定义
func (d *definitionService) DeleteDefinition(id uint) error {
	err := global.BankDb.Delete(model.ProcessDefinition{}, "id=?", id).Error

	if err != nil {
		return errors.New("流程不存在")
	}

	return nil
}
