/**
 * @Author: lzw5399
 * @Date: 2021/3/18 22:00
 * @Desc: 审批/处理流程实例的接口的请求体
 */
package request

import "workflow/src/model"

// 审批/处理流程实例的接口的请求体
type HandleInstancesRequest struct {
	EdgeId            string                   `json:"edgeId" form:"edgeId"`                       // 走的流程的id
	ProcessInstanceId uint                     `json:"processInstanceId" form:"processInstanceId"` // 流程实例的id
	Remarks           string                   `json:"remarks" form:"remarks"`                     // 备注
	Variables         []model.InstanceVariable `json:"variables"`                                  // 变量
}

// 否决流程的请求体
type DenyInstanceRequest struct {
	ProcessInstanceId uint   `json:"processInstanceId" form:"processInstanceId"` // 流程实例的id
	NodeId            string `json:"nodeId" form:"nodeId"`                       // 所在节点的id
	Remarks           string `json:"remarks" form:"remarks"`                     // 备注
}
