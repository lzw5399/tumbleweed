/**
 * @Author: lzw5399
 * @Date: 2021/3/9 15:54
 * @Desc: 流程实例路由
 */
package router

import (
	"workflow/src/controller"
	customMiddleware "workflow/src/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterProcessInstance(r *echo.Echo) {
	instanceGroup := r.Group("/api/process-instances", customMiddleware.Auth)
	{
		instanceGroup.POST("", controller.CreateProcessInstance)         // 新建流程
		instanceGroup.GET("/:id", controller.GetProcessInstance)         // 获取
		instanceGroup.GET("", controller.ListProcessInstances)           // 获取列表
		instanceGroup.POST("/_handle", controller.HandleProcessInstance) // 流程审批
	}
}
