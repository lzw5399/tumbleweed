/**
 * @Author: lzw5399
 * @Date: 2021/01/15
 * @Desc: process控制器
 */
package controller

import (
	"encoding/xml"
	"net/http"

	"workflow/global/response"
	"workflow/model/request"
	"workflow/service"

	"github.com/labstack/echo/v4"
)

// 创建流程process
func CreateProcess(c echo.Context) error {
	var r request.BpmnRequest
	if err := c.Bind(&r); err != nil {
		return response.Failed(c, http.StatusBadRequest)
	}

	var bpmnDefinitions request.Definitions
	err := xml.Unmarshal([]byte(r.Data), &bpmnDefinitions)
	if err != nil {
		return response.FailWithMsg(c, http.StatusBadRequest, "不是标准的bpmn2.0定义的流程，请使用工作流设计器创建流程")
	}

	if err := service.CreateProcess(&bpmnDefinitions.Process, r.Data); err != nil {
		return response.Failed(c, http.StatusInternalServerError)
	}

	return response.Ok(c)
}
