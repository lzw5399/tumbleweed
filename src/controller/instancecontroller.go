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
// @param current-user header string true "current-user"
// @Success 200 {object} response.HttpResponse
// @Router /api/process-instances [post]
func CreateProcessInstance(c echo.Context) error {
	var r request.ProcessInstanceRequest
	if err := c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	currentUserId := util.GetCurrentUserId(c)
	result, err := instanceService.CreateProcessInstance(&r, currentUserId)
	if err != nil {
		return response.FailWithMsg(c, int(result), err)
	}

	return response.OkWithData(c, result)
}

// process instance list
func ListProcessInstances(c echo.Context) error {
	// 从queryString获取分页参数
	var r request.PagingRequest
	if err := c.Bind(&r); err != nil {
		return response.Failed(c, http.StatusBadRequest)
	}

	instances, err := instanceService.List(&r)
	if err != nil {
		return response.Failed(c, http.StatusInternalServerError)
	}

	return response.OkWithData(c, instances)
}

// 获取一个实例
func GetProcessInstance(c echo.Context) error {
	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return response.Failed(c, http.StatusBadRequest)
	}

	instance, err := instanceService.Get(uint(id))
	if err != nil {
		return response.Failed(c, http.StatusNotFound)
	}

	return response.OkWithData(c, instance)
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
