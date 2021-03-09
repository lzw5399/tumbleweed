/**
 * @Author: lzw5399
 * @Date: 2020/9/30 15:13
 * @Desc: config model
 */
package config

type Config struct {
	Log   Log   `yaml:"log"`
	App   App   `yaml:"app"`
	Db    Db    `yaml:"db"`
	Redis Redis `yaml:"redis"`
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
	LogMode     bool   `yaml:"logmode"`
	AutoMigrate bool   `yaml:"auto_migrate"`
}

type Log struct {
	Prefix string `yaml:"prefix"`
	Stdout string `yaml:"stdout"`
}

type Redis struct {
	ConnStr string `yaml:"conn_str"`
	Enabled bool   `yaml:"enabled"`
}
