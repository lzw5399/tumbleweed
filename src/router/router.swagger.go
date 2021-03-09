/**
 * @Author: lzw5399
 * @Date: 2021/3/9 15:47
 * @Desc:
 */
package router

import "github.com/labstack/echo/v4"

func RegisterSwagger(r *echo.Echo) {
	// swagger
	//r.GET("/api/dq/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.GET("/api/dq/swagger", func(c *gin.Context) {
	//	c.Redirect(http.StatusMovedPermanently, "/api/dq/swagger/index.html")
	//})
}
