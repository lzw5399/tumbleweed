/**
 * @Author: lzw5399
 * @Date: 2021/1/16 21:43
 * @Desc: 流程实例
 */
package controller

import (
	"net/http"
	"workflow/src/model/request"
	"workflow/src/service"

	"workflow/src/global/response"

	"github.com/gin-gonic/gin"
)

// 创建新的实例
func StartProcessInstance(c *gin.Context) {
	var r request.InstanceRequest
	if err := c.Bind(&r); err != nil {
		response.Failed(c, http.StatusBadRequest)
		return
	}

	if statusCode, err := service.StartProcessInstance(&r); err != nil {
		response.FailWithMsg(c, statusCode, err)
		return
	}

	response.OkWithData(c, r)
}
