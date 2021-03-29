/**
 * @Author: lzw5399
 * @Date: 2020/9/30 15:13
 * @Desc: config model
 */
package config

type Config struct {
	App App `yaml:"app"`
	Db  Db  `yaml:"db"`
}

type App struct {
	Name          string `yaml:"name"`
	EnableSwagger bool   `yaml:"enable_swagger"`
}

type Db struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Database    string `yaml:"database"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	LogMode     bool   `yaml:"log_mode"`
	AutoMigrate bool   `yaml:"auto_migrate"`
}
