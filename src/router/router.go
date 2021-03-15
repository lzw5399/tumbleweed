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
)

func Setup() *echo.Echo {
	r := echo.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.Use(middleware.CORS())

	// probe
	r.GET("/", controller.Index)
	r.GET("/api/info/ready", controller.Readiness)
	r.GET("/api/info/alive", controller.Liveliness)

	// swagger
	if global.BankConfig.App.EnableSwagger {
		RegisterSwagger(r)
	}

	// apis
	RegisterProcessDefinition(r)
	RegisterProcessInstance(r)

	return r
}
