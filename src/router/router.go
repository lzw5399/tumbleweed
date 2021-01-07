/**
 * @Author: lzw5399
 * @Date: 2020/9/30 13:44
 * @Desc: application main router
 */
package router

import (
	"net/http"

	"bank/distributedquery/src/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// default allow all origins
	r.Use(cors.Default())

	// swagger
	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/api/swagger/index.html")
	})

	// static
	r.LoadHTMLGlob("./src/app/views/*")
	r.Static("/assets", "./src/app/assets")
	r.StaticFile("/favicon.ico", "./src/app/assets/favicon.ico")

	// APIs
	r.GET("/", controller.Index)
	r.GET("/api/info", controller.Info)
	r.GET("/api/info/ready", controller.Readiness)
	r.GET("/api/info/alive", controller.Liveliness)

	ocrGroup := r.Group("/api/ocr")
	{
		ocrGroup.POST("file", controller.ScanFile)
		ocrGroup.POST("scan-crop-file", controller.ScanCropFile)
		ocrGroup.POST("base64", controller.Base64)
		ocrGroup.POST("scan-crop-base64", controller.ScanCropBase64)
	}

	return r
}
