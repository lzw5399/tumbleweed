/**
 * @Author: lzw5399
 * @Date: 2021/01/15
 * @Desc: ocr related functionality
 */
package controller

import (
	"encoding/xml"
	"net/http"
	"workflow/src/service"

	"workflow/src/global/response"
	"workflow/src/model/request"

	"github.com/gin-gonic/gin"
)

// 创建流程process
func CreateProcess(c *gin.Context) {
	var r request.BpmnRequest
	if err := c.ShouldBind(&r); err != nil {
		response.Failed(c, http.StatusBadRequest)
		return
	}

	var bpmnDefinitions request.Definitions
	err := xml.Unmarshal([]byte(r.Data), &bpmnDefinitions)
	if err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, "不是标准的bpmn2.0定义的流程，请使用工作流设计器创建流程")
		return
	}

	if err := service.CreateProcess(&bpmnDefinitions); err != nil {
		response.Failed(c, http.StatusInternalServerError)
	}

	response.Ok(c)
}
