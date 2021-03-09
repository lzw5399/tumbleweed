/**
 * @Author: lzw5399
 * @Date: 2021/3/9 13:12
 * @Desc: 流程分类
 */
package model


// 流程分类
type Classify struct {
	AuditableBase
	Name    string `json:"name" form:"name"`     // 分类名称
}