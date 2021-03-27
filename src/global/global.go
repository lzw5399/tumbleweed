/**
 * @Author: lzw5399
 * @Date: 2021/01/12 15:11:13
 * @Desc: 全局对象, 将在程序启动后初始化
 */
package global

import (
	golog "github.com/op/go-logging"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"

	"workflow/src/config"
)

var (
	BankConfig config.Config
	BankLogger *golog.Logger
	BankDb     *gorm.DB
	BankCache  *cache.Cache
)
