/**
 * @Author: lzw5399
 * @Date: 2021/1/16 11:37
 * @Desc: 顺序流
 */
package model

type SequenceFlow struct {
	DbBase
	SourceRef           string `json:"sourceRef"`
	TargetRef           string `json:"targetRef"`
	ConditionExpression string `json:"conditionExpression"`
}
