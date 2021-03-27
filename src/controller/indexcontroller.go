/**
 * @Author: lzw5399
 * @Date: 2021/01/07 14:22
 * @Desc: home page controller
 */
package controller

import (
	"errors"
	"time"

	"workflow/src/global"
	"workflow/src/global/response"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	err := errors.New("错误")
	global.BankLogger.Info("测试info")
	global.BankLogger.Debug("测试debug")
	global.BankLogger.Warning("测试warning")
	global.BankLogger.Error("测试error",err)
	return response.OkWithData(c, time.Now().Local())
}

func Liveliness(c echo.Context) error {
	return response.Ok(c)
}

func Readiness(c echo.Context) error {
	return response.Ok(c)
}
