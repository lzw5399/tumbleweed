/**
 * @Author: lzw5399
 * @Date: 2020/9/30 13:44
 * @Desc: application main router
 */
package router

import (
	"net/http"

	"workflow/src/controller"

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

	// static
	r.LoadHTMLGlob("./src/app/views/*")
	r.Static("/assets", "./src/app/assets")
	r.StaticFile("/favicon.ico", "./src/app/assets/favicon.ico")

	// APIs
	r.GET("/", controller.Index)
	r.GET("/api/info/ready", controller.Readiness)
	r.GET("/api/info/alive", controller.Liveliness)

	ocrGroup := r.Group("/api/workflow")
	{
		ocrGroup.POST("file", controller.ScanFile)
	}

	return r
}
