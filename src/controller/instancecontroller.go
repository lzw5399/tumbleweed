/**
 * @Author: lzw5399
 * @Date: 2021/1/16 21:43
 * @Desc: 流程实例
 */
package controller

import (
	"net/http"
	"strconv"

	"workflow/src/global/response"
	"workflow/src/model/request"
	"workflow/src/service"
	"workflow/src/util"

	"github.com/labstack/echo/v4"
)

var (
	instanceService service.InstanceService = service.NewInstanceService()
)

// @Tags process-instances
// @Summary 创建新的流程实例
// @Accept  json
// @Produce json
// @param request body request.ProcessInstanceRequest true "request"
// @param wf-tenant-code header string true "wf-tenant-code"
// @param wf-current-user header string true "wf-current-user"
// @Success 200 {object} response.HttpResponse
// @Router /api/process-instances [post]
func CreateProcessInstance(c echo.Context) error {
	var r request.ProcessInstanceRequest
	if err := c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	currentUserId := util.GetCurrentUserId(c)
	tenantId := util.GetCurrentTenantId(c)
	processInstance, err := instanceService.CreateProcessInstance(&r, currentUserId, tenantId)
	if err != nil {
		return response.FailWithMsg(c, http.StatusInternalServerError, err)
	}

	return response.OkWithData(c, processInstance)
}

// @Tags process-instances
// @Summary 获取流程实例列表
// @Accept  json
// @Produce json
// @param request query request.InstanceListRequest true "request"
// @param wf-tenant-code header string true "wf-tenant-code"
// @param wf-current-user header string true "wf-current-user"
// @Success 200 {object} response.HttpResponse
// @Router /api/process-instances [GET]
func ListProcessInstances(c echo.Context) error {
	// 从queryString获取分页参数
	var r request.InstanceListRequest
	if err := c.Bind(&r); err != nil {
		return response.Failed(c, http.StatusBadRequest)
	}

	tenantId := util.GetCurrentTenantId(c)
	instances, err := instanceService.List(&r, util.GetCurrentUserId(c), tenantId)
	if err != nil {
		return response.FailWithMsg(c, http.StatusInternalServerError, err)
	}

	return response.OkWithData(c, instances)
}

// @Tags process-instances
// @Summary 处理/审批一个流程
// @Accept  json
// @Produce json
// @param request body request.HandleInstancesRequest true "request"
// @param wf-tenant-code header string true "wf-tenant-code"
// @param wf-current-user header string true "wf-current-user"
// @Success 200 {object} response.HttpResponse
// @Router /api/process-instances/_handle [POST]
func HandleProcessInstance(c echo.Context) error {
	var r request.HandleInstancesRequest
	if err := c.Bind(&r); err != nil {
		return response.Failed(c, http.StatusBadRequest)
	}

	tenantId := util.GetCurrentTenantId(c)
	instance, err := instanceService.HandleProcessInstance(&r, util.GetCurrentUserId(c), tenantId)
	if err != nil {
		return response.FailWithMsg(c, http.StatusBadRequest, err)
	}

	return response.OkWithData(c, instance)
}

// @Tags process-instances
// @Summary 获取一个流程实例
// @Produce json
// @param id path int true "request"
// @param includeProcessTrain query bool false "request"
// @param wf-tenant-code header string true "wf-tenant-code"
// @param wf-current-user header string true "wf-current-user"
// @Success 200 {object} response.HttpResponse
// @Router /api/process-instances/{id} [GET]
func GetProcessInstance(c echo.Context) error {
	var r request.GetInstanceRequest
	if err := c.Bind(&r); err != nil {
		return response.Failed(c, http.StatusBadRequest)
	}

	currentUserId := util.GetCurrentUserId(c)
	tenantId := util.GetCurrentTenantId(c)
	instance, err := instanceService.Get(&r, currentUserId, tenantId)
	if err != nil {
		return response.FailWithMsg(c, http.StatusNotFound, "记录不存在")
	}

	return response.OkWithData(c, instance)
}

// @Tags process-instances
// @Summary 获取流程链路
// @Produce json
// @param id path int true "request"
// @param wf-tenant-code header string true "wf-tenant-code"
// @param wf-current-user header string true "wf-current-user"
// @Success 200 {object} response.HttpResponse
// @Router /api/process-instances/{id}/train-nodes [GET]
func GetProcessTrain(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.Failed(c, http.StatusBadRequest)
	}

	tenantId := util.GetCurrentTenantId(c)
	trainNodes, err := instanceService.GetProcessTrain(nil, uint(id), tenantId)
	if err != nil {
		return response.InternalServerErrorWithMessage(c, err)
	}

	return response.OkWithData(c, trainNodes)
}

// 获取流程实例中的变量
//func GetInstanceVariable(c echo.Context) error {
//	var r request.GetVariableRequest
//	if err := c.Bind(&r); err != nil {
//		return response.Failed(c, http.StatusBadRequest)
//	}
//
//	resp, err := instanceService.GetVariable(&r)
//	if err != nil {
//		return response.FailWithMsg(c, http.StatusInternalServerError, err)
//	}
//
//	return response.OkWithData(c, resp)
//}
//
//func GetInstanceVariableList(c echo.Context) error {
//	var r request.GetVariableListRequest
//	if err := c.Bind(&r); err != nil {
//		return response.Failed(c, http.StatusBadRequest)
//	}
//	variables, err := instanceService.ListVariables(&r)
//	if err != nil {
//		return response.FailWithMsg(c, http.StatusInternalServerError, err)
//	}
//
//	return response.OkWithData(c, variables)
//}
