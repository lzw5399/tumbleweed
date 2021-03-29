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

// @Tags health
// @Accept  json
// @Produce json
// @Success 200 {object} response.HttpResponse
// @Router /health/alive [GET]
func Liveliness(c echo.Context) error {
	return response.Ok(c)
}

// @Tags health
// @Accept  json
// @Produce json
// @Success 200 {object} response.HttpResponse
// @Router /health/ready [GET]
func Readiness(c echo.Context) error {
	return response.Ok(c)
}
