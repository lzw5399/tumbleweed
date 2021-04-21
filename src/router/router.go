/**
 * @Author: lzw5399
 * @Date: 2020/9/30 13:44
 * @Desc: application main router
 */
package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"workflow/src/controller"
	"workflow/src/global"
	customMiddleware "workflow/src/middleware"
)

func Setup() *echo.Echo {
	r := echo.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.Use(middleware.CORS())

	// probe
	r.GET("/", controller.Index)
	r.GET("/health/ready", controller.Readiness)
	r.GET("/health/alive", controller.Liveliness)

	// swagger
	if global.BankConfig.App.EnableSwagger {
		RegisterSwagger(r)
	}

	// apis
	g := r.Group("/api/wf", customMiddleware.MultiTenant, customMiddleware.Auth)
	{
		RegisterProcessDefinition(g) // 流程定义
		RegisterProcessInstance(g)   // 流程实例
		RegisterRoleUsers(g)         // 外部系统的角色用户映射
	}

	return r
}
