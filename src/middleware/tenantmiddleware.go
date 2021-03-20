/**
 * @Author: lzw5399
 * @Date: 2021/3/20 14:57
 * @Desc: 多租户中间件
 */
package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	. "github.com/ahmetb/go-linq/v3"
	"github.com/labstack/echo/v4"

	"workflow/src/global"
	"workflow/src/global/response"
	"workflow/src/model"
	"workflow/src/util"
)

func MultiTenant(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tenantCode := c.Request().Header.Get("wf-tenant-code")
		if tenantCode == "" {
			return response.FailWithMsg(c, http.StatusUnauthorized, "未指定当前租户")
		}

		// 从内存缓存中获取全部的tenant
		tenants := GetTenants()

		tenant := From(tenants).WhereT(func(i model.Tenant) bool {
			return i.Name == tenantCode
		}).First()

		// 不存在新增并更新全部租户缓存
		if tenant == nil {
			t := model.Tenant{
				Name:       tenantCode,
				CreateTime: time.Now().Local(),
			}
			err := global.BankDb.Model(&model.Tenant{}).Create(&t).Error
			if err != nil {
				return response.FailWithMsg(c, http.StatusUnauthorized, "指定租户失败")
			}
			tenant = t

			// 更新缓存
			go UpdateTenantCache()
		}

		// 当前租户放进缓存
		c.Set("currentTenant", tenant)

		return next(c)
	}
}

// 从缓存中获取所有的租户信息
func GetTenants() []model.Tenant {
	var tenants []model.Tenant

	t, succeed := global.BankCache.Get("tenants")
	if !succeed {
		tenants = *(UpdateTenantCache())
	}
	bytes := util.MarshalToBytes(t)
	_ = json.Unmarshal(bytes, &tenants)

	return tenants
}

// 更新新的租户信息到缓存中
func UpdateTenantCache() *[]model.Tenant {
	var tenants []model.Tenant
	global.BankDb.
		Model(&model.Tenant{}).
		Find(&tenants)

	global.BankCache.SetDefault("tenants", tenants)
	log.Print("租户缓存更新成功")

	return &tenants
}
