/**
 * @Author: lzw5399
 * @Date: 2021/1/16 22:38
 * @Desc: 流程实例
 */
package model

import "gorm.io/datatypes"

// 流程实例
type ProcessInstance struct {
	EntityBase
	ProcessId  uint           `json:"processId" gorm:"index:idx_processId3"`
	Variables  datatypes.JSON `json:"variables"`                       // 流程变量，在流程生命周期中可用 [{"name":"",type:"",value: emm}]
	IsFinished bool           `json:"isFinished" gorm:"default:false"` // 是否已完成
}

// 对应ProcessInstance.Variables
type InstanceVariable struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
