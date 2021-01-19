/**
 * @Author: lzw5399
 * @Date: 2020/9/30 13:44
 * @Desc: application main router
 */
package router

import (
	"workflow/controller"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	//swaggerFiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup() *echo.Echo {
	r := echo.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())

	// default allow all origins
	r.Use(middleware.CORS())

	// swagger
	//r.GET("/api/wf/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.GET("/api/wf/swagger", func(c *gin.Context) {
	//	c.Redirect(http.StatusMovedPermanently, "/api/dq/swagger/index.html")
	//})

	// APIs
	r.GET("/", controller.Index)
	r.GET("/api/info/ready", controller.Readiness)
	r.GET("/api/info/alive", controller.Liveliness)

	processGroup := r.Group("/api/process")
	{
		processGroup.POST("create", controller.CreateProcess)
	}

	instanceGroup := r.Group("/api/instance")
	{
		instanceGroup.POST("start", controller.StartProcessInstance)
		instanceGroup.GET("list", controller.ListProcessInstances)
		instanceGroup.GET(":id", controller.GetProcessInstance)
		//instanceGroup.GET("variable/get", controller.GetInstanceVariable)
		//instanceGroup.GET("variable/list", controller.GetInstanceVariableList)
		//instanceGroup.POST("variable/set", controller.GetProcessInstance)
	}

	return r
}
