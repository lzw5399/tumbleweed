/**
 * @Author: lzw5399
 * @Date: 2021/3/9 18:18
 * @Desc:
 */
package util

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetCurrentUserId(c echo.Context) uint {
	u := c.Get("currentUser").(string)
	i, _ := strconv.Atoi(u)

	return uint(i)
}