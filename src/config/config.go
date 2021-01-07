/**
 * @Author: lzw5399
 * @Date: 2020/9/30 15:13
 * @Desc: config model
 */
package config

type Config struct {
	Log Log `yaml:"log"`
	App App `yaml:"app"`
	Ocr Ocr `yaml:"ocr"`
}

type App struct {
	Name string `yaml:"name"`
}

type Ocr struct {
	DefaultLangs []string `yaml:"default-langs"`
}

type Log struct {
	Prefix  string `yaml:"prefix"`
	LogFile bool   `yaml:"log-file"`
	Stdout  string `yaml:"stdout"`
	File    string `yaml:"file"`
}
