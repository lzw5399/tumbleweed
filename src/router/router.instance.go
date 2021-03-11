/**
 * @Author: lzw5399
 * @Date: 2021/3/9 15:54
 * @Desc:
 */
package router

import (
	"workflow/src/controller"

	"github.com/labstack/echo/v4"
)

func RegisterProcessInstance(r *echo.Echo) {
	instanceGroup := r.Group("/api/process-instances")
	{
		instanceGroup.POST("", controller.CreateProcessInstance)
		instanceGroup.GET("/get", controller.GetProcessInstance)
		instanceGroup.GET("/list", controller.ListProcessInstances)
		instanceGroup.POST("/variable/set", controller.GetProcessInstance)
	}
}
