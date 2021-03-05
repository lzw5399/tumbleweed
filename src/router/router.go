/**
 * @Author: lzw5399
 * @Date: 2020/9/30 13:44
 * @Desc: application main router
 */
package router

import (
	"workflow/src/controller"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Setup() *echo.Echo {
	r := echo.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.Use(middleware.CORS())

	// swagger
	//r.GET("/api/dq/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.GET("/api/dq/swagger", func(c *gin.Context) {
	//	c.Redirect(http.StatusMovedPermanently, "/api/dq/swagger/index.html")
	//})

	// APIs
	r.GET("/", controller.Index)
	r.GET("/api/info/ready", controller.Readiness)
	r.GET("/api/info/alive", controller.Liveliness)

	processGroup := r.Group("/api/process")
	{
		processGroup.POST("/create", controller.CreateProcessDefinition)
	}

	instanceGroup := r.Group("/api/instance")
	{
		instanceGroup.POST("/start", controller.StartProcessInstance)
		instanceGroup.GET("/get", controller.GetProcessInstance)
		instanceGroup.GET("/list", controller.ListProcessInstances)
		instanceGroup.GET("/variable/get", controller.GetInstanceVariable)
		instanceGroup.GET("/variable/list", controller.GetInstanceVariableList)
		instanceGroup.POST("/variable/set", controller.GetProcessInstance)
	}

	taskGroup := r.Group("/api/task")
	{
		taskGroup.GET("/list", controller.ListTasks)
	}

	return r
}
