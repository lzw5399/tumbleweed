/**
 * @Author: lzw5399
 * @Date: 2021/1/12 15:12
 * @Desc: 初始化数据库连接
 */
package initialize

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"

	"workflow/src/global"

	_ "github.com/jackc/pgx/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	log.Print("-------开始初始化postgres数据库连接--------")

	// 获取数据库配置
	dbCfg := global.BankConfig.Db

	// 初始化数据库连接
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Shanghai",
		dbCfg.Host, dbCfg.Port, dbCfg.Username, dbCfg.InitialDb, dbCfg.Password)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  connStr,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true, // 忽略迁移添加外键
	})

	if err != nil {
		log.Fatalf("数据库连接初始化失败, 原因: %s", err.Error())
	}

	global.BankDb = db
	log.Print("-------初始化postgres数据库连接成功--------")
}
