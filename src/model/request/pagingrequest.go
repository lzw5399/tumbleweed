/**
 * @Author: lzw5399
 * @Date: 2021/1/17 18:29
 * @Desc: 分页请求base model
 */
package request

type PagingRequest struct {
	Sort   string `json:"sort,omitempty" form:"sort,omitempty"`     // 排序键的名字，在各查询实现中默认值与可用值都不同
	Order  string `json:"order,omitempty" form:"order,omitempty"`   // asc或者是desc
	Offset int    `json:"offset,omitempty" form:"offset,omitempty"` // 跳过的条数
	Limit  int    `json:"limit,omitempty" form:"limit,omitempty"`   // 取的条数
}
