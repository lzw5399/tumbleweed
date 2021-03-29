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

	loadEnvToStruct(v, []string{})
}

func loadEnvToStruct(v reflect.Value, dependencies []string) {
	num := v.NumField()
	for i := 0; i < num; i++ {
		f := v.Field(i)
		fieldName := v.Type().Field(i).Name

		envValue := os.Getenv(genDependenciesString(dependencies, fieldName))

		if envValue == "" || !f.CanSet() {
			continue
		}

		log.Printf("当前环境变量: %s, 已加载", envValue)

		switch f.Kind() {
		case reflect.String:
			f.SetString(envValue)
		case reflect.Bool:
			if value, err := strconv.ParseBool(envValue); err == nil {
				f.SetBool(value)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value, err := strconv.Atoi(envValue); err == nil {
				f.SetInt(int64(value))
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if value, err := strconv.Atoi(envValue); err == nil {
				f.SetUint(uint64(value))
			}
		case reflect.Struct:
			dependencies = append(dependencies, fieldName)
			loadEnvToStruct(f, dependencies)
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
