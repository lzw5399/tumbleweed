/**
 * @Author: lzw5399
 * @Date: 2021/1/16 11:37
 * @Desc: 顺序流
 */
package model

type SequenceFlow struct {
	TableBase
	Code                string `json:"code" gorm:"uniqueIndex"`
	SourceRef           string `json:"sourceRef"`
	TargetRef           string `json:"targetRef"`
	ConditionExpression string `json:"conditionExpression"`
	ProcessId           uint   `json:"processId" gorm:"index:idx_processId4"`
}
