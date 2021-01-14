/**
 * @Author: lzw5399
 * @Date: 2021/1/14 23:09
 * @Desc: 流程实例（执行流）表
 */
package model

// 流程实例（执行流）
type Execution struct {
	DbBase
	Rev         int    `json:"rev"`
	ProcInstID  int    `json:"procInstID"`
	ProcDefID   int    `json:"procDefID"`
	ProcDefName string `json:"procDefName"`
	NodeInfos   string `json:"nodeInfos"` // 执行流经过的所有节点
	IsActive    int8   `json:"isActive"`
	StartTime   string `json:"startTime"`
}
