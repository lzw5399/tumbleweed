/**
 * @Author: lzw5399
 * @Date: 2021/1/16 21:43
 * @Desc: 流程实例
 */
package controller

import (
	"net/http"
	"strconv"

	"workflow/global/response"
	"workflow/model/request"
	"workflow/service"

	"github.com/gin-gonic/gin"
)

var (
	instanceSvc service.InstanceService = service.NewInstanceService()
)

// 创建新的实例
func StartProcessInstance(c *gin.Context) {
	var r request.InstanceRequest
	if err := c.ShouldBind(&r); err != nil {
		response.Failed(c, http.StatusBadRequest)
		return
	}

	result, err := instanceSvc.Start(&r)
	if err != nil {
		response.FailWithMsg(c, int(result), err)
		return
	}

	response.OkWithData(c, result)
}

// 获取一个实例
func GetProcessInstance(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		response.Failed(c, http.StatusBadRequest)
		return
	}

	instance, err := instanceSvc.Get(uint(id))
	if err != nil {
		response.Failed(c, http.StatusNotFound)
		return
	}

	response.OkWithData(c, instance)
}

// process instance list
func ListProcessInstances(c *gin.Context) {
	// 从queryString获取分页参数
	var r request.PagingRequest
	if err := c.BindQuery(&r); err != nil {
		response.Failed(c, http.StatusBadRequest)
		return
	}

	instances, err := instanceSvc.List(&r)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError)
		return
	}

	response.OkWithData(c, instances)
}

// 获取流程实例中的变量
func GetInstanceVariable(c *gin.Context) {
	var r request.GetVariableRequest
	if err := c.BindQuery(&r); err != nil {
		response.Failed(c, http.StatusBadRequest)
		return
	}

	resp, err := instanceSvc.GetVariable(&r)
	if err != nil {
		response.FailWithMsg(c, http.StatusInternalServerError, err)
		return
	}

	response.OkWithData(c, resp)
}

func GetInstanceVariableList(c *gin.Context) {
	var r request.GetVariableListRequest
	if err := c.BindQuery(&r); err != nil {
		response.Failed(c, http.StatusBadRequest)
		return
	}
	variables, err := instanceSvc.ListVariables(&r)
	if err != nil {
		response.FailWithMsg(c, http.StatusInternalServerError, err)
		return
	}

	response.OkWithData(c, variables)
}
