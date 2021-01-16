/**
 * @Author: lzw5399
 * @Date: 2021/1/16 22:38
 * @Desc: 流程实例
 */
package model

import "gorm.io/datatypes"

type ProcessInstance struct {
	TableBase
	ProcessId  uint           `json:"processId" gorm:"index:idx_processId3"`
	Variables  datatypes.JSON `json:"variables"`                       // 流程变量，在流程生命周期中可用
	IsFinished bool           `json:"isFinished" gorm:"default:false"` // 是否已完成
}
