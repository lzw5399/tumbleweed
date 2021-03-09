/**
 * @Author: lzw5399
 * @Date: 2021/3/9 18:18
 * @Desc:
 */
package util

import "github.com/labstack/echo/v4"

func GetCurrentUserId(c echo.Context) uint {
	return c.Get("currentUser").(uint)
}
