/**
 * @Author: lzw5399
 * @Date: 2021/01/15
 * @Desc: process控制器
 */
package controller

import (
	"log"
	"net/http"

	"workflow/src/global/response"
	"workflow/src/model/request"
	"workflow/src/service"

	"github.com/labstack/echo/v4"
)

var (
	definitionService service.DefinitionService = service.NewDefinitionService()
)

// 创建流程process
func CreateProcessDefinition(c echo.Context) error {
	var (
		r   request.ProcessDefinitionRequest
		err error
	)

	if err = c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	// 验证
	err = definitionService.Validate(&r, -1)
	if err != nil {
		return response.BadRequestWithMessage(c, err)
	}

	// 创建
	processDefinition, err := definitionService.CreateProcess(&r)
	if err != nil {
		log.Printf("CreateProcess错误，原因: %s", err.Error())
		return response.Failed(c, http.StatusInternalServerError)
	}

	return response.OkWithData(c, processDefinition)
}

// 更新模板
func UpdateProcessDefinition(c echo.Context) error {
	var (
		r   request.ProcessDefinitionRequest
		err error
	)

	if err = c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	// 验证
	err = definitionService.Validate(&r, r.Id)
	if err != nil {
		return response.BadRequestWithMessage(c, err)
	}

	err = definitionService.UpdateProcess(&r)
	if err != nil {
		log.Printf("UpdateProcessDefinition错误，原因: %s", err.Error())
		return response.Failed(c, http.StatusInternalServerError)
	}

	return response.Ok(c)
}

// 删除流程
func DeleteProcess(c echo.Context) error {
	processId := c.DefaultQuery("processId", "")
	if processId == "" {
		app.Error(c, -1, errors.New("参数不正确，请确定参数processId是否传递"), "")
		return
	}

	err := orm.Eloquent.Delete(process2.Info{}, "id = ?", processId).Error
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("删除流程失败, %v", err.Error()))
		return
	}
	app.OK(c, "", "删除流程成功")
}
