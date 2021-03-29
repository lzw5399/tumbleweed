/**
 * @Author: lzw5399
 * @Date: 2021/01/15
 * @Desc: process控制器
 */
package controller

import (
	"github.com/labstack/echo/v4"

	"workflow/src/global"
	"workflow/src/global/response"
	"workflow/src/model/request"
	"workflow/src/service"
	"workflow/src/util"
)

var (
	definitionService service.DefinitionService = service.NewDefinitionService()
)

// @Tags process-definitions
// @Summary 创建流程模板
// @Accept  json
// @Produce json
// @param request body request.ProcessDefinitionRequest true "request"
// @param WF-TENANT-CODE header string true "WF-TENANT-CODE"
// @param WF-CURRENT-USER header string true "WF-CURRENT-USER"
// @Success 200 {object} response.HttpResponse
// @Router /api/wf/process-definitions [POST]
func CreateProcessDefinition(c echo.Context) error {
	var (
		r   request.ProcessDefinitionRequest
		err error
	)

	if err = c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	// 验证
	tenantId := util.GetCurrentTenantId(c)
	err = definitionService.Validate(&r, 0, tenantId)
	if err != nil {
		return response.Failed(c, err)
	}

	// 创建
	currentUserId := util.GetCurrentUserId(c)
	processDefinition, err := definitionService.CreateDefinition(&r, currentUserId, tenantId)
	if err != nil {
		global.BankLogger.Error("CreateProcess错误", err)
		return response.Failed(c, err)
	}

	return response.OkWithData(c, processDefinition)
}

// @Tags process-definitions
// @Summary 更新流程模板
// @Accept  json
// @Produce json
// @param request body request.ProcessDefinitionRequest true "request"
// @param WF-TENANT-CODE header string true "WF-TENANT-CODE"
// @param WF-CURRENT-USER header string true "WF-CURRENT-USER"
// @Success 200 {object} response.HttpResponse
// @Router /api/wf/process-definitions [PUT]
func UpdateProcessDefinition(c echo.Context) error {
	var (
		r   request.ProcessDefinitionRequest
		err error
	)

	if err = c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	// 验证
	tenantId := util.GetCurrentTenantId(c)
	err = definitionService.Validate(&r, r.Id, tenantId)
	if err != nil {
		return response.Failed(c, err)
	}

	currentUserId := util.GetCurrentUserId(c)
	err = definitionService.UpdateDefinition(&r, currentUserId, tenantId)
	if err != nil {
		global.BankLogger.Error("UpdateProcessDefinition错误", err)
		return response.Failed(c, err)
	}

	return response.Ok(c)
}

// @Tags process-definitions
// @Summary 删除流程模板
// @Produce json
// @param id path string true "request"
// @param WF-TENANT-CODE header string true "WF-TENANT-CODE"
// @param WF-CURRENT-USER header string true "WF-CURRENT-USER"
// @Success 200 {object} response.HttpResponse
// @Router /api/wf/process-definitions/{id} [DELETE]
func DeleteProcessDefinition(c echo.Context) error {
	definitionId := c.Param("id")
	if definitionId == "" {
		return response.BadRequestWithMessage(c, "参数不正确，请确定参数processDefinitionId是否传递")
	}

	tenantId := util.GetCurrentTenantId(c)
	err := definitionService.DeleteDefinition(util.StringToUint(definitionId), tenantId)
	if err != nil {
		return response.Failed(c, err)
	}

	return response.Ok(c)
}

// @Tags process-definitions
// @Summary 获取流程模板详情
// @Produce json
// @param id path string true "request"
// @param WF-TENANT-CODE header string true "WF-TENANT-CODE"
// @param WF-CURRENT-USER header string true "WF-CURRENT-USER"
// @Success 200 {object} response.HttpResponse
// @Router /api/wf/process-definitions/{id} [GET]
func GetProcessDefinition(c echo.Context) error {
	definitionId := c.Param("id")
	if definitionId == "" {
		return response.BadRequestWithMessage(c, "参数不正确，请确定参数processDefinitionId是否传递")
	}

	tenantId := util.GetCurrentTenantId(c)
	definition, err := definitionService.GetDefinition(util.StringToUint(definitionId), tenantId)
	if err != nil {
		return response.Failed(c, err)
	}

	return response.OkWithData(c, definition)
}

// @Tags process-definitions
// @Summary 获取流程定义列表
// @Accept  json
// @Produce json
// @param request query request.DefinitionListRequest true "request"
// @param WF-TENANT-CODE header string true "WF-TENANT-CODE"
// @param WF-CURRENT-USER header string true "WF-CURRENT-USER"
// @Success 200 {object} response.HttpResponse
// @Router /api/wf/process-definitions [GET]
func ListProcessDefinition(c echo.Context) error {
	// 从queryString获取分页参数
	var r request.DefinitionListRequest
	if err := c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	tenantId := util.GetCurrentTenantId(c)
	instances, err := definitionService.List(&r, util.GetCurrentUserId(c), tenantId)
	if err != nil {
		return response.Failed(c, err)
	}

	return response.OkWithData(c, instances)
}

//// 分类流程列表
//func ClassifyProcessList(c echo.Context) error {
//	var (
//		err            error
//		classifyIdList []int
//		classifyList   []*struct {
//			process2.Classify
//			ProcessList []*process2.Info `json:"process_list"`
//		}
//	)
//
//	processName := c.DefaultQuery("name", "")
//	if processName == "" {
//		err = orm.Eloquent.Model(&process2.Classify{}).Find(&classifyList).Error
//		if err != nil {
//			app.Error(c, -1, err, fmt.Sprintf("获取分类列表失败，%v", err.Error()))
//			return
//		}
//	} else {
//		err = orm.Eloquent.Model(&process2.Info{}).
//			Where("name LIKE ?", fmt.Sprintf("%%%v%%", processName)).
//			Pluck("distinct classify", &classifyIdList).Error
//		if err != nil {
//			app.Error(c, -1, err, fmt.Sprintf("获取分类失败，%v", err.Error()))
//			return
//		}
//
//		err = orm.Eloquent.Model(&process2.Classify{}).
//			Where("id in (?)", classifyIdList).
//			Find(&classifyList).Error
//		if err != nil {
//			app.Error(c, -1, err, fmt.Sprintf("获取分类失败，%v", err.Error()))
//			return
//		}
//	}
//
//	for _, item := range classifyList {
//		err = orm.Eloquent.Model(&process2.Info{}).
//			Where("classify = ? and name LIKE ?", item.Id, fmt.Sprintf("%%%v%%", processName)).
//			Select("id, create_time, update_time, name, icon, remarks").
//			Find(&item.ProcessList).Error
//		if err != nil {
//			app.Error(c, -1, err, fmt.Sprintf("获取流程失败，%v", err.Error()))
//			return
//		}
//	}
//
//	app.OK(c, classifyList, "成功获取数据")
//}
