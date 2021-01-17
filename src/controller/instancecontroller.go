/**
 * @Author: lzw5399
 * @Date: 2021/1/16 21:43
 * @Desc: 流程实例
 */
package controller

import (
	"net/http"

	"workflow/src/global/response"
	"workflow/src/model/request"
	"workflow/src/service"

	"github.com/gin-gonic/gin"
)

// 创建新的实例
func StartProcessInstance(c *gin.Context) {
	var r request.InstanceRequest
	if err := c.ShouldBind(&r); err != nil {
		response.Failed(c, http.StatusBadRequest)
		return
	}

	result, err := service.StartProcessInstance(&r)
	if err != nil {
		response.FailWithMsg(c, int(result), err)
		return
	}

	response.OkWithData(c, result)
}
