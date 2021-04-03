/**
 * @Author: lzw5399
 * @Date: 2021/3/9 18:18
 * @Desc:
 */
package util

import (
	"encoding/json"

	"github.com/labstack/echo/v4"

	"workflow/src/model"
)

func GetUserIdentifier(c echo.Context) string {
	u := c.Get("currentUser").(string)

	return u
}

func GetCurrentTenant(c echo.Context) (tenants model.Tenant) {
	u := c.Get("currentTenant")
	bytes := MarshalToBytes(u)
	_ = json.Unmarshal(bytes, &tenants)

	return
}

func GetCurrentTenantId(c echo.Context) int {
	return GetCurrentTenant(c).Id
}

func GetWorkContext(c echo.Context) (int, string) {
	return GetCurrentTenantId(c), GetUserIdentifier(c)
}
