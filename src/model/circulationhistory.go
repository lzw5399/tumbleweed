/**
 * @Author: lzw5399
 * @Date: 2021/3/11 15:09
 * @Desc: 工单流转历史
 */
package model

// 工单流转历史
type CirculationHistory struct {
	AuditableBase
	Title             string `json:"title" form:"title"`                          // 工单标题
	ProcessInstanceId int    `json:"processInstanceId" form:"processInstanceId"`  // 工单ID
	SourceState       string `json:"state" form:"state"`                          // 源节点label
	SourceId          string `json:"sourceId" form:"sourceId"`                    // 源节点ID
	TargetId          string `json:"targetId" form:"targetId"`                    // 目标节点ID
	Circulation       string `json:"circulation" form:"circulation"`              // 流转说明
	ProcessorId       string `gorm:"index" json:"processorId" form:"processorId"` // 处理人外部系统ID
	CostDuration      string `json:"costDuration" form:"costDuration"`            // 本条记录的处理时长(每次有新的一条的时候更新这个字段)
	Remarks           string `json:"remarks" form:"remarks"`                      // 备注
}
