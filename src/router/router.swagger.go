/**
 * @Author: lzw5399
 * @Date: 2021/3/9 15:47
 * @Desc:
 */
package router

import (
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
)

const (
	_SWAGGER_BASE_PATH = "/api/wf/swagger"
)

func RegisterSwagger(r *echo.Echo) {
	r.GET(path.Join(_SWAGGER_BASE_PATH, "/*"), echoSwagger.WrapHandler)
	r.GET(_SWAGGER_BASE_PATH, func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, path.Join(_SWAGGER_BASE_PATH, "index.html"))
	})
}
