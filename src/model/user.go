/**
 * @Author: lzw5399
 * @Date: 2021/1/18 21:12
 * @Desc: candidateUsers & assignee
 */
package model

type User struct {
	EntityBase
	Identifier string `json:"identifier"` // 外部系统的标识
	Name       string `json:"name"`
	TenantId   string `json:"tenantId"`
}
