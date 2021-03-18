/**
 * @Author: lzw5399
 * @Date: 2021/1/17 18:29
 * @Desc: 分页请求base model
 */
package request

type PagingRequest struct {
	Sort   string `json:"sort,omitempty" form:"sort,omitempty" query:"sort"`       // 排序键的名字，在各查询实现中默认值与可用值都不同
	Order  string `json:"order,omitempty" form:"order,omitempty" query:"order"`    // asc或者是desc
	Offset int    `json:"offset,omitempty" form:"offset,omitempty" query:"offset"` // 跳过的条数
	Limit  int    `json:"limit,omitempty" form:"limit,omitempty" query:"limit"`    // 取的条数
}

type InstanceListRequest struct {
	PagingRequest
	Type int `json:"type,omitempty" form:"type" query:"type"` // 类别
}
