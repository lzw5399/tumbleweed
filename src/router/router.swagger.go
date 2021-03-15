/**
 * @Author: lzw5399
 * @Date: 2021/3/9 15:47
 * @Desc:
 */
package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
)

func RegisterSwagger(r *echo.Echo) {
	r.GET("/api/wf/swagger/*", echoSwagger.WrapHandler)
	r.GET("/api/wf/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/api/wf/swagger/index.html")
	})
}
