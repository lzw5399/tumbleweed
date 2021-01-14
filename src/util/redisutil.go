/**
 * @Author: lzw5399
 * @Date: 2021/1/14 23:40
 * @Desc:
 */
package util

import (
	"log"
	"time"

	"workflow/src/global"
)

// 优先从缓存中取
func GetOrCreate(key string, getFromDb func() interface{}) interface{} {
	if !global.BankConfig.Redis.Enabled {
		return getFromDb()
	}

	val, err := global.BankRedis.Get(key).Result()
	if err != nil {
		log.Printf("从redis获取失败，原因：%s", err.Error())
		return getFromDbAndSetToRedis(key, getFromDb)
	}
	if val != "" {
		return val
	}
	return getFromDbAndSetToRedis(key, getFromDb)
}

func getFromDbAndSetToRedis(key string, getFromDb func() interface{}) interface{} {
	val := getFromDb()
	global.BankRedis.Set(key, val, time.Hour*24)
	return val
}
