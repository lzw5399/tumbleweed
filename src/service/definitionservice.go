/**
 * @Author: lzw5399
 * @Date: 2021/1/15 23:35
 * @Desc:
 */
package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"

	"workflow/src/global"
	"workflow/src/global/constant"
	"workflow/src/global/shared"
	"workflow/src/model"
	"workflow/src/model/request"
	"workflow/src/model/response"
	"workflow/src/util"
)

type DefinitionService interface {
	CreateDefinition(*request.ProcessDefinitionRequest, uint, uint) (*model.ProcessDefinition, error)
	Validate(*request.ProcessDefinitionRequest, uint, uint) error
	UpdateDefinition(*request.ProcessDefinitionRequest, uint, uint) error
	DeleteDefinition(id uint, tenantId uint) error
	GetDefinition(id uint, tenantId uint) (*model.ProcessDefinition, error)
	List(r *request.DefinitionListRequest, currentUserId uint, tenantId uint) (interface{}, error)
}

func NewDefinitionService() *definitionService {
	return &definitionService{}
}

type definitionService struct {
}

func (d *definitionService) GetDefinition(id uint, tenantId uint) (*model.ProcessDefinition, error) {
	var definition model.ProcessDefinition

	err := global.BankDb.
		Where("id=?", id).
		Where("tenant_id=?", tenantId).
		First(&definition).Error
	if err != nil {
		log.Error(err)
		return nil, errors.New("查询流程详情失败")
	}

	return &definition, nil
}

// 验证
func (d *definitionService) Validate(r *request.ProcessDefinitionRequest, excludeId uint, tenantId uint) error {
	// 验证名称是否已存在
	var c int64
	global.BankDb.Model(&model.ProcessDefinition{}).
		Where("name=?", r.Name).
		Where("id!=?", excludeId).
		Where("tenant_id=?", tenantId).
		Count(&c)
	if c != 0 {
		return errors.New(fmt.Sprintf("当前名称为:\"%s\"的模板已存在", r.Name))
	}

	// 如果edge对象不存在id，则生成一个
	var definitionStructure map[string][]map[string]interface{}
	err := json.Unmarshal(r.Structure, &definitionStructure)
	if err != nil {
		return errors.New("当前structure不合法，请检查")
	}

	for _, edge := range definitionStructure["edges"] {
		if edge["id"] == nil {
			edge["id"] = fmt.Sprintf("flow_%s", util.GenUUID())
		}
	}
	r.Structure = util.MarshalToBytes(definitionStructure)

	// todo 校验structure的json

	return nil
}

// 创建新的process流程
func (d *definitionService) CreateDefinition(r *request.ProcessDefinitionRequest, currentUserId uint, tenantId uint) (*model.ProcessDefinition, error) {
	var (
		err error
	)

	processDefinition := r.ProcessDefinition()
	processDefinition.CreateBy = currentUserId
	processDefinition.UpdateBy = currentUserId
	processDefinition.TenantId = int(tenantId)

	if err = global.BankDb.Create(&processDefinition).Error; err != nil {
		log.Error(err)
		return nil, err
	}

	return &processDefinition, nil
}

// 更新流程定义
func (d *definitionService) UpdateDefinition(r *request.ProcessDefinitionRequest, currentUserId uint, tenantId uint) error {
	processDefinition := r.ProcessDefinition()

	// 先查询
	var count int64
	err := global.BankDb.Model(&model.ProcessDefinition{}).
		Where("id=?", processDefinition.Id).
		Where("tenant_id=?", tenantId).
		Count(&count).
		Error
	if err != nil || count == 0 {
		return errors.New("记录不存在")
	}

	err = global.BankDb.
		Model(&processDefinition).
		Updates(map[string]interface{}{
			"name":        processDefinition.Name,
			"form_id":     processDefinition.FormId,
			"structure":   processDefinition.Structure,
			"classify_id": processDefinition.ClassifyId,
			"task":        processDefinition.Task,
			"notice":      processDefinition.Notice,
			"remarks":     processDefinition.Remarks,
			"update_by":   currentUserId,
			"update_time": time.Now().Local(),
		}).Error

	return err
}

// 删除流程定义
func (d *definitionService) DeleteDefinition(id uint, tenantId uint) error {
	// 先查询
	var count int64
	err := global.BankDb.Model(&model.ProcessDefinition{}).
		Where("id=?", id).
		Where("tenant_id=?", tenantId).
		Count(&count).
		Error
	if err != nil || count == 0 {
		return errors.New("记录不存在")
	}

	err = global.BankDb.Delete(model.ProcessDefinition{}, "id=?", id).Error

	if err != nil {
		return errors.New("流程不存在")
	}

	return nil
}

func (d *definitionService) List(r *request.DefinitionListRequest, currentUserId uint, tenantId uint) (interface{}, error) {
	var definitions []model.ProcessDefinition
	db := global.BankDb.Model(&model.ProcessDefinition{}).Where("tenant_id = ?", tenantId)

	// 根据type的不同有不同的逻辑
	switch r.Type {
	case constant.D_ICreated:
		db = db.Where("create_by=?", currentUserId)
		break
	case constant.D_All:
		break
	default:
		return nil, errors.New("type不合法")
	}

	if r.Keyword != "" {
		db = db.Where("name ~ ?", r.Keyword)
	}

	var count int64
	db.Count(&count)

	db = shared.ApplyPaging(db, &r.PagingRequest)
	err := db.Find(&definitions).Error

	return &response.PagingResponse{
		TotalCount:   count,
		CurrentCount: int64(len(definitions)),
		Data:         &definitions,
	}, err
}
