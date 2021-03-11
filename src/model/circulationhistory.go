/**
 * @Author: lzw5399
 * @Date: 2021/3/11 15:09
 * @Desc: 工单流转历史
 */
package model

// 工单流转历史
type CirculationHistory struct {
	EntityBase
	Title        string `json:"title" form:"title"`                 // 工单标题
	WorkOrder    int    `json:"work_order" form:"work_order"`       // 工单ID
	State        string `json:"state" form:"state"`                 // 工单状态
	Source       string `json:"source" form:"source"`               // 源节点ID
	Target       string `json:"target" form:"target"`               // 目标节点ID
	Circulation  string `json:"circulation" form:"circulation"`     // 流转ID
	Processor    string `json:"processor" form:"processor"`         // 处理人
	ProcessorId  int    `json:"processor_id" form:"processor_id"`   // 处理人ID
	CostDuration string `json:"cost_duration" form:"cost_duration"` // 处理时长
	Remarks      string `json:"remarks" form:"remarks"`             // 备注
}
