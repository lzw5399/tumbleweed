/**
 * @Author: lzw5399
 * @Date: 2021/1/12 15:12
 * @Desc: 初始化数据库连接
 */
package initialize

import (
	"fmt"
	"log"

	"workflow/src/global"
	"workflow/src/model"

	_ "github.com/jackc/pgx/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func init() {
	log.Println("-------开始初始化postgres数据库连接--------")

	// 获取数据库配置
	dbCfg := global.BankConfig.Db

	// 初始化数据库连接
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Shanghai",
		dbCfg.Host, dbCfg.Port, dbCfg.Username, dbCfg.Database, dbCfg.Password)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  connStr,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true, // 忽略迁移添加外键
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "wf.", // 指定schema
			SingularTable: true,  // 表名不加s
		},
	})

	if err != nil {
		log.Fatalf("数据库连接初始化失败, 原因: %s", err.Error())
	}

	global.BankDb = db
	log.Println("-------初始化postgres数据库连接成功--------")

	if dbCfg.AutoMigrate {
		log.Println("-------启动了数据库自动迁移, 开始迁移--------")
		doMigration()
		log.Println("-------表结构迁移完成--------")
	}
}

// 自动迁移
func doMigration() {
	global.BankDb.Exec("create schema if not exists emm;")
	err := global.BankDb.AutoMigrate(
		&model.Process{}, &model.Event{},
		&model.ExclusiveGateway{}, &model.SequenceFlow{},
		&model.UserTask{}, &model.ProcessInstance{})
	if err != nil {
		log.Fatalf("迁移发生错误，错误信息为:%s", err.Error())
	}
}
