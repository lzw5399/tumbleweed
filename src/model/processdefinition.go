/**
 * @Author: lzw5399
 * @Date: 2021/1/14 22:13
 * @Desc: 流程的定义
 */
package model

import "gorm.io/gorm"

// 流程定义表
type ProcessDefinition struct {
	gorm.Model
	Name       string `json:"name,omitempty"`
	Version    int    `json:"version,omitempty"`
	Resource   string `gorm:"size:10000" json:"resource,omitempty"` // 流程定义json字符串
	Userid     string `json:"userid,omitempty"`                     // 用户id
	Username   string `json:"username,omitempty"`
	Company    string `json:"company,omitempty"` // 用户所在公司
	DeployTime string `json:"deployTime,omitempty"`
}
