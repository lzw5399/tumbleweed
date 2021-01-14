/**
 * @Author: lzw5399
 * @Date: 2021/1/14 23:11
 * @Desc:
 */
package model

import "time"

type DbBase struct {
	Id         uint `gorm:"primarykey"`
	CreateTime time.Time
	UpdateTime time.Time
}
