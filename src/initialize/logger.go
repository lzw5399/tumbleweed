/**
 * @Author: lzw5399
 * @Date: 2020/9/30 15:25
 * @Desc: auto load logger after app start
 */
package initialize

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"workflow/src/global"
)

func setupLogger() {
	logger := log.New()
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true) // 日志中添加调用方法

	// 配置格式化器
	if os.Getenv("APP_ENV") == "Production" {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetFormatter(&log.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
		log.SetLevel(log.DebugLevel)
	}

	global.BankLogger = logger
}
