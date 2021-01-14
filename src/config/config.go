/**
 * @Author: lzw5399
 * @Date: 2020/9/30 15:13
 * @Desc: config model
 */
package config

type Config struct {
	Log Log `yaml:"log"`
	App App `yaml:"app"`
	Db  Db  `yaml:"db"`
}

type App struct {
	Name string `yaml:"name"`
}

type Db struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	InitialDb   string `yaml:"initial_db"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	MaxIdleConn int    `yaml:"max_idle_conn"`
	MaxOpenConn int    `yaml:"max_open_conn"`
	LogMode     bool   `yaml:"logmode"`
}

type Log struct {
	Prefix string `yaml:"prefix"`
	Stdout string `yaml:"stdout"`
}
