/**
 * @Author: lzw5399
 * @Date: 2021/01/12 15:11:13
 * @Desc: 全局对象, 将在程序启动后初始化
 */
package global

import (
	"workflow/src/config"

	golog "github.com/op/go-logging"
	"gorm.io/gorm"
)

var (
	BankConfig config.Config
	BankLogger *golog.Logger
	BankDb     *gorm.DB
)
