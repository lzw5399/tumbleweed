/**
 * @Author: lzw5399
 * @Date: 2021/2/9 18:42
 * @Desc:
 */
package controller

import (
	"net/http"

	"workflow/src/global/response"
	"workflow/src/model/request"
	"workflow/src/service"

	"github.com/labstack/echo/v4"
)

var (
	taskSvc service.TaskService = service.NewTaskService()
)

func ListTasks(c echo.Context) error {
	var r request.PagingRequest
	if err := c.Bind(&r); err != nil {
		return response.Failed(c, http.StatusBadRequest)
	}

	list, err := taskSvc.List(&r)
	if err != nil {
		return response.Failed(c, http.StatusInternalServerError)
	}

	return response.OkWithData(c, list)
}

