/**
 * @Author: lzw5399
 * @Date: 2021/3/9 15:52
 * @Desc: 流程定义路由
 */
package router

import (
	"workflow/src/controller"
	customMiddleware "workflow/src/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterProcessDefinition(r *echo.Echo) {
	processGroup := r.Group("/api/process-definitions", customMiddleware.Auth)
	{
		processGroup.POST("", controller.CreateProcessDefinition)       // 新建
		processGroup.PUT("", controller.UpdateProcessDefinition)        // 修改
		processGroup.DELETE("/:id", controller.DeleteProcessDefinition) // 删除
		processGroup.GET("/:id", controller.GetProcessDefinition)       // 获取流程
	}
}
