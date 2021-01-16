/**
 * @Author: lzw5399
 * @Date: 2021/1/14 22:13
 * @Desc: 流程的定义
 */
package model

// 流程定义表
type Process struct {
	TableBase
	Code                string         `json:"code" gorm:"uniqueIndex"`
	Name                string         `json:"name"`                                                           // 流程名字
	Category            string         `json:"category"`                                                       // 流程类别
	Version             int            `json:"version,omitempty"`                                              // 版本
	Resource            string         `json:"resource"`                                                       // 流程定义的完整bpmn xml字符串
}
