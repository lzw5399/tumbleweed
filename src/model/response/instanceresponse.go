/**
 * @Author: lzw5399
 * @Date: 2021/1/17 21:34
 * @Desc:
 */
package response

import (
	"workflow/src/global/constant"
	"workflow/src/model"
)

type InstanceVariableResponse struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type ProcessInstanceResponse struct {
	model.ProcessInstance
	ProcessChainNodes []ProcessChainNode `json:"processChainNodes,omitempty"` // 流程链路【包括全部节点和当前节点】
}

type ProcessChainNode struct {
	Name       string                   `json:"name"`
	Id         string                   `json:"id"`
	Status     constant.ChainNodeStatus `json:"status"`     // 1: 已处理 2: 当前节点 3: 后续节点
	Sort       int                      `json:"sort"`       // 排序
	NodeType   int                      `json:"nodeType"`   // 1. 开始事件 2. 用户任务 3. 排他网关 4. 结束事件
	Obligatory bool                     `json:"obligatory"` // 是否必经节点
}
