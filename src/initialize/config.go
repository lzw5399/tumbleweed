/**
 * @Author: lzw5399
 * @Date: 2020/9/30 15:17
 * @Desc: auto load config setting after app start
 */
package initialize

import (
	"fmt"
	"log"
	"os"

	"workflow/src/global"
	"workflow/src/util"

	"github.com/jinzhu/configor"
)

func setupConfig() {
	envCode := getEnvCode()
	overrideConfigFileName := fmt.Sprintf("src/config/appsettings.%s.yaml", envCode)

	var err error
	if util.PathExists(overrideConfigFileName) {
		err = configor.Load(&global.BankConfig, "src/config/appsettings.yaml", overrideConfigFileName)
	} else {
		err = configor.Load(&global.BankConfig, "src/config/appsettings.yaml")
	}

	if err != nil {
		log.Fatalf("配置初始化失败, 原因:%s", err.Error())
	}

	// 加载环境变量
	util.LoadEnv(&global.BankConfig)
}

func getEnvCode() string {
	envMode := os.Getenv("APP_ENV")

	if envMode == "Production" {
		return "Production"
	}

	return "Development"
}
