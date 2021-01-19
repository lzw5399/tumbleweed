/**
 * @Author: lzw5399
 * @Date: 2021/1/17 16:48
 * @Desc:
 */
package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateDeployment(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
