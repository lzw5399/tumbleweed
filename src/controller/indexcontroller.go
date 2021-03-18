/**
 * @Author: lzw5399
 * @Date: 2021/01/07 14:22
 * @Desc: home page controller
 */
package controller

import (
	"time"

	"workflow/src/global/response"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	return response.OkWithData(c, time.Now().Local())
}

func Liveliness(c echo.Context) error {
	return response.Ok(c)
}

func Readiness(c echo.Context) error {
	return response.Ok(c)
}
