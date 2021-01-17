/**
 * @Author: lzw5399
 * @Date: 2021/1/12 17:45
 * @Desc: 初始化redis连接
 */
package initialize

import (
	"log"

	"workflow/src/global"

	"github.com/go-redis/redis"
)

func init() {
	log.Println("-------开始初始化redis连接--------")

	client := redis.NewClient(&redis.Options{
	})

	global.BankRedis = client
	log.Println("-------初始化redis连接成功--------")
}
