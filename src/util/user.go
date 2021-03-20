/**
 * @Author: lzw5399
 * @Date: 2021/3/9 18:18
 * @Desc:
 */
package util

import (
	"encoding/json"
	"strconv"

	"github.com/labstack/echo/v4"

	"workflow/src/model"
)

func GetCurrentUserId(c echo.Context) uint {
	u := c.Get("currentUser").(string)
	i, _ := strconv.Atoi(u)

	return uint(i)
}

func GetCurrentTenant(c echo.Context) (tenants model.Tenant) {
	u := c.Get("currentTenant")
	bytes := MarshalToBytes(u)
	_ = json.Unmarshal(bytes, &tenants)

	return
}

func GetCurrentTenantId(c echo.Context) uint {
	return GetCurrentTenant(c).Id
}
