/**
 * @Author: lzw5399
 * @Date: 2021/3/9 15:42
 * @Desc:
 */
package middleware

import (
	"net/http"

	"workflow/src/global/response"

	"github.com/labstack/echo/v4"
)

// 先使用此种方式传递当前用户的标识id
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUserId := c.Request().Header.Get("current-user")
		if currentUserId == "" {
			return response.FailWithMsg(c, http.StatusUnauthorized, "未指定当前用户")
		}
		c.Set("currentUser", currentUserId)

		return next(c)
	}
}
