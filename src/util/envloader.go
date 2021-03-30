/**
 * @Author: lzw5399
 * @Date: 2021/1/22 15:18
 * @Desc: 获取环境变量
 */
package util

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func LoadEnv(val interface{}) {
	v := reflect.ValueOf(val)

	// 获取到指针指向的值
	v = reflect.Indirect(v)

	if v.Kind() != reflect.Struct {
		fmt.Println("请使用结构体指针")
		return
	}

	loadEnvToStruct(v, []string{}, 1)
}

func loadEnvToStruct(v reflect.Value, dependencies []string, deepLevel int) {
	num := v.NumField()
	for i := 0; i < num; i++ {
		if deepLevel == 1 {
			dependencies = []string{}
		}
		f := v.Field(i)
		fieldName := v.Type().Field(i).Name

		envKey := genDependenciesString(dependencies, fieldName)
		envValue := os.Getenv(envKey)

		if f.Kind() != reflect.Struct && (envValue == "" || !f.CanSet()) {
			log.Printf("【环境变量配置加载】当前envKey为: %s 的环境变量为空, 将跳过", envKey)
			continue
		}

		switch f.Kind() {
		case reflect.String:
			f.SetString(envValue)
		case reflect.Bool:
			if value, err := strconv.ParseBool(envValue); err == nil {
				log.Printf("【环境变量配置加载】当前envKey为: %s 的环境变量已成功替换", envKey)
				f.SetBool(value)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value, err := strconv.Atoi(envValue); err == nil {
				log.Printf("【环境变量配置加载】当前envKey为: %s 的环境变量已替换", envKey)
				f.SetInt(int64(value))
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if value, err := strconv.Atoi(envValue); err == nil {
				log.Printf("【环境变量配置加载】当前envKey为: %s 的环境变量已替换", envKey)
				f.SetUint(uint64(value))
			}
		case reflect.Struct:
			dependencies = append(dependencies, fieldName)
			loadEnvToStruct(f, dependencies, deepLevel+1)
		}
	}
}

func genDependenciesString(dependencies []string, fieldName string) string {
	prefix := strings.Join(dependencies, "__")
	if prefix == "" {
		return fieldName
	}

	return fmt.Sprintf("%s__%s", prefix, fieldName)
}
