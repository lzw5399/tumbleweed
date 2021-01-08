/**
 * @Author: lzw5399
 * @Date: 2020/9/30 15:17
 * @Desc: auto load config setting after app start
 */
package initialize

import (
	"fmt"
	"os"
	"strings"

	"bank/workflow/engine/src/global"
	"bank/workflow/engine/src/util"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
)

// 同一个包内有多个init，则会按文件名顺序执行
// 部分其他的init会依赖当前这个config的init
// 所以如果需要依赖，不要让package内的其他文件排在config.go前面
func init() {
	envCode := getEnvCode()
	overrideConfigFileName := fmt.Sprintf("config/appsettings.%s.yaml", envCode)

	var err error
	if util.PathExists(overrideConfigFileName) {
		err = configor.Load(&global.BANK_CONFIG, "config/appsettings.yaml", overrideConfigFileName)
	} else {
		err = configor.Load(&global.BANK_CONFIG, "config/appsettings.yaml")
	}

	if err != nil {
		panic("resolve settings failed...")
	}
}

var envMap = map[string]string{
	"debug":   "Development",
	"release": "Production",
}

func getEnvCode() string {
	ginMode := os.Getenv(gin.EnvGinMode)

	for k, v := range envMap {
		if k == strings.ToLower(ginMode) {
			return v
		}
	}

	return "Development"
}

