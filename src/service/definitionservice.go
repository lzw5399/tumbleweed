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

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"workflow/src/global"
	"workflow/src/global/constant"
	"workflow/src/global/shared"
	"workflow/src/model"
	"workflow/src/model/request"
	"workflow/src/model/response"
	"workflow/src/util"
)

func GetDefinition(id int, tenantId int) (*model.ProcessDefinition, error) {
	var definition model.ProcessDefinition

	err := global.BankDb.
		Where("id=?", id).
		Where("tenant_id=?", tenantId).
		First(&definition).Error
	if err != nil {
		global.BankLogger.Error(err)
		return nil, util.NewError("查询流程详情失败")
	}

	return &definition, nil
}

// 验证
func ValidateDefinitionRequest(r *request.ProcessDefinitionRequest, excludeId int, tenantId int) error {
	// 验证名称是否已存在
	var c int64
	global.BankDb.Model(&model.ProcessDefinition{}).
		Where("name=?", r.Name).
		Where("id!=?", excludeId).
		Where("tenant_id=?", tenantId).
		Count(&c)
	if c != 0 {
		return util.BadRequest.Newf("当前名称为:\"%s\"的模板已存在", r.Name)
	}

	// 如果edge对象不存在id，则生成一个
	var definitionStructure map[string][]map[string]interface{}
	err := json.Unmarshal(r.Structure, &definitionStructure)
	if err != nil {
		return util.BadRequest.New("当前structure不合法，请检查")
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
func CreateDefinition(r *request.ProcessDefinitionRequest, c echo.Context) (*model.ProcessDefinition, error) {
	var (
		processDefinition        = r.ProcessDefinition()
		tenantId, userIdentifier = util.GetWorkContext(c)
	)
	processDefinition.CreateBy = userIdentifier
	processDefinition.UpdateBy = userIdentifier
	processDefinition.TenantId = tenantId

	err := global.BankDb.Create(&processDefinition).Error
	if err != nil {
		log.Error(err)
		return nil, util.NewError("创建失败")
	}

	return &processDefinition, nil
}

// 更新流程定义
func UpdateDefinition(r *request.ProcessDefinitionRequest, c echo.Context) error {
	var (
		processDefinition        = r.ProcessDefinition()
		tenantId, userIdentifier = util.GetWorkContext(c)
	)

	// 先查询
	var count int64
	err := global.BankDb.Model(&model.ProcessDefinition{}).
		Where("id=?", processDefinition.Id).
		Where("tenant_id=?", tenantId).
		Count(&count).
		Error
	if err != nil || count == 0 {
		return util.NotFound.New("记录不存在")
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
			"update_by":   userIdentifier,
			"update_time": time.Now().Local(),
		}).Error

	return err
}

// 删除流程定义
func DeleteDefinition(id int, tenantId int) error {
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

func GetDefinitionList(r *request.DefinitionListRequest, c echo.Context) (interface{}, error) {
	var (
		definitions              []model.ProcessDefinition
		tenantId, userIdentifier = util.GetWorkContext(c)
	)

	db := global.BankDb.Model(&model.ProcessDefinition{}).Where("tenant_id = ?", tenantId)

	// 根据type的不同有不同的逻辑
	switch r.Type {
	case constant.D_ICreated:
		db = db.Where("create_by=?", userIdentifier)
		break
	case constant.D_All:
		break
	default:
		return nil, util.BadRequest.New("type不合法")
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

func CloneDefinition(r *request.CloneDefinitionRequest, c echo.Context) (*model.ProcessDefinition, error) {
	var (
		tenantId, _ = util.GetWorkContext(c)
	)
	definition, err := GetDefinition(r.Id, tenantId)
	if err != nil {
		return nil, err
	}

	newD := *definition
	fmt.Printf("xinde: %p,jiude: %p,san:%p", &newD, definition,&*definition)

	return nil, nil
}