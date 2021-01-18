/**
 * @Author: lzw5399
 * @Date: 2021/1/17 22:42
 * @Desc:
 */
package response

type PagingResponse struct {
	TotalCount   int64       `json:"totalCount"`   // 总数
	CurrentCount int64       `json:"currentCount"` // 当前数量
	Data         interface{} `json:"data"`         // 具体数据
}
