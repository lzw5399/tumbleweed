/**
 * @Author: lzw5399
 * @Date: 2021/3/27 0:18
 * @Desc:
 */
package response

import "workflow/src/model/dto"

type TrainNodesResponse struct {
	dto.Node
	Obligatory bool `json:"obligatory"` // 是否是必经节点
}
