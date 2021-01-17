/**
 * @Author: lzw5399
 * @Date: 2021/1/17 16:29
 * @Desc: 已部署(deploy)的process definition
 */
package model

type Deployment struct {
	EntityBase
	ProcessId uint           `json:"processId" gorm:"index:idx_processId6"`
}
