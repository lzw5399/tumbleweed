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

// 同一个包内有多个init，则会按文件名顺序执行
// 部分其他的init会依赖当前这个config的init
// 所以如果需要依赖，不要让package内的其他文件排在config.go前面
func init() {
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
