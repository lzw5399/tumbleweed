/**
 * @Author: lzw5399
 * @Date: 2020/9/30 13:44
 * @Desc: application main router
 */
package router

import (
	"net/http"

	"workflow/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// default allow all origins
	r.Use(cors.Default())

	// swagger
	r.GET("/api/dq/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/api/dq/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/api/dq/swagger/index.html")
	})

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
		instanceGroup.GET("get", controller.GetProcessInstance)
		instanceGroup.GET("list", controller.ListProcessInstances)
		instanceGroup.GET("variable/get", controller.GetInstanceVariable)
		instanceGroup.GET("variable/list", controller.GetInstanceVariableList)
		instanceGroup.POST("variable/set", controller.GetProcessInstance)
	}

	return r
}
