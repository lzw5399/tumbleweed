/**
 * @Author: lzw5399
 * @Date: 2021/1/15 23:35
 * @Desc:
 */
package service

import (
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"time"

	"workflow/src/global"
	"workflow/src/model"
	"workflow/src/model/request"
)

type DefinitionService interface {
	CreateProcess(*request.ProcessDefinitionRequest) (*model.ProcessDefinition, error)
	Validate(*request.ProcessDefinitionRequest, uint) error
	UpdateProcess(r *request.ProcessDefinitionRequest) error
}

type definitionService struct {
}

func NewDefinitionService() *definitionService {
	return &definitionService{}
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
func (d *definitionService) CreateProcess(r *request.ProcessDefinitionRequest) (*model.ProcessDefinition, error) {
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

func (d *definitionService) UpdateProcess(r *request.ProcessDefinitionRequest) error {
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
