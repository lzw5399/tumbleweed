/**
 * @Author: lzw5399
 * @Date: 2021/1/16 21:43
 * @Desc: 流程实例
 */
package controller

import (
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
// @param WF-TENANT-CODE header string true "WF-TENANT-CODE"
// @param WF-CURRENT-USER header string true "WF-CURRENT-USER"
// @Success 200 {object} response.HttpResponse
// @Router /api/wf/process-instances [post]
func CreateProcessInstance(c echo.Context) error {
	var r request.ProcessInstanceRequest
	if err := c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	currentUserId := util.GetCurrentUserId(c)
	tenantId := util.GetCurrentTenantId(c)
	processInstance, err := instanceService.CreateProcessInstance(&r, currentUserId, tenantId)
	if err != nil {
		return response.Failed(c, err)
	}

	return response.OkWithData(c, processInstance)
}

// @Tags process-instances
// @Summary 获取流程实例列表
// @Accept  json
// @Produce json
// @param request query request.InstanceListRequest true "request"
// @param WF-TENANT-CODE header string true "WF-TENANT-CODE"
// @param WF-CURRENT-USER header string true "WF-CURRENT-USER"
// @Success 200 {object} response.HttpResponse
// @Router /api/wf/process-instances [GET]
func ListProcessInstances(c echo.Context) error {
	// 从queryString获取分页参数
	var r request.InstanceListRequest
	if err := c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	tenantId := util.GetCurrentTenantId(c)
	instances, err := instanceService.ListProcessInstance(&r, util.GetCurrentUserId(c), tenantId)
	if err != nil {
		return response.Failed(c, err)
	}

	return response.OkWithData(c, instances)
}

// @Tags process-instances
// @Summary 处理/审批一个流程
// @Accept  json
// @Produce json
// @param request body request.HandleInstancesRequest true "request"
// @param WF-TENANT-CODE header string true "WF-TENANT-CODE"
// @param WF-CURRENT-USER header string true "WF-CURRENT-USER"
// @Success 200 {object} response.HttpResponse
// @Router /api/wf/process-instances/_handle [POST]
func HandleProcessInstance(c echo.Context) error {
	var r request.HandleInstancesRequest
	if err := c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	tenantId := util.GetCurrentTenantId(c)
	instance, err := instanceService.HandleProcessInstance(&r, util.GetCurrentUserId(c), tenantId)
	if err != nil {
		return response.Failed(c, err)
	}

	return response.OkWithData(c, instance)
}

// @Tags process-instances
// @Summary 否决流程流程
// @Accept  json
// @Produce json
// @param request body request.DenyInstanceRequest true "request"
// @param WF-TENANT-CODE header string true "WF-TENANT-CODE"
// @param WF-CURRENT-USER header string true "WF-CURRENT-USER"
// @Success 200 {object} response.HttpResponse
// @Router /api/wf/process-instances/_deny [POST]
func DenyProcessInstance(c echo.Context) error {
	var r request.DenyInstanceRequest
	if err := c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	tenantId := util.GetCurrentTenantId(c)
	currentUserId := util.GetCurrentUserId(c)
	instance, err := instanceService.DenyProcessInstance(&r, currentUserId, tenantId)
	if err != nil {
		return response.Failed(c, err)
	}

	return response.OkWithData(c, instance)
}

// @Tags process-instances
// @Summary 获取一个流程实例
// @Produce json
// @param id path int true "request"
// @param includeProcessTrain query bool false "request"
// @param WF-TENANT-CODE header string true "WF-TENANT-CODE"
// @param WF-CURRENT-USER header string true "WF-CURRENT-USER"
// @Success 200 {object} response.HttpResponse
// @Router /api/wf/process-instances/{id} [GET]
func GetProcessInstance(c echo.Context) error {
	var r request.GetInstanceRequest
	if err := c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	currentUserId := util.GetCurrentUserId(c)
	tenantId := util.GetCurrentTenantId(c)
	instance, err := instanceService.GetProcessInstance(&r, currentUserId, tenantId)
	if err != nil {
		return response.Failed(c, err)
	}

	return response.OkWithData(c, instance)
}

// @Tags process-instances
// @Summary 获取流程链路
// @Produce json
// @param id path int true "request"
// @param WF-TENANT-CODE header string true "WF-TENANT-CODE"
// @param WF-CURRENT-USER header string true "WF-CURRENT-USER"
// @Success 200 {object} response.HttpResponse
// @Router /api/wf/process-instances/{id}/train-nodes [GET]
func GetProcessTrain(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.BadRequest(c)
	}

	tenantId := util.GetCurrentTenantId(c)
	trainNodes, err := instanceService.GetProcessTrain(nil, uint(id), tenantId)
	if err != nil {
		return response.Failed(c, err)
	}

	return response.OkWithData(c, trainNodes)
}

// 获取流程实例中的变量
//func GetInstanceVariable(c echo.Context) error {
//	var r request.GetVariableRequest
//	if err := c.Bind(&r); err != nil {
//		return response.FailedOblete(c, http.StatusBadRequest)
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
//		return response.FailedOblete(c, http.StatusBadRequest)
//	}
//	variables, err := instanceService.ListVariables(&r)
//	if err != nil {
//		return response.FailWithMsg(c, http.StatusInternalServerError, err)
//	}
//
//	return response.OkWithData(c, variables)
//}
