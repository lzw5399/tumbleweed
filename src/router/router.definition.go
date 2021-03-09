/**
 * @Author: lzw5399
 * @Date: 2021/3/9 15:52
 * @Desc:
 */
package router

import (
	"workflow/src/controller"

	"github.com/labstack/echo/v4"
)

func RegisterProcessDefinition(r *echo.Echo) {
	processGroup := r.Group("/api/process-definitions")
	{
		processGroup.POST("/create", controller.CreateProcessDefinition)
	}
}
