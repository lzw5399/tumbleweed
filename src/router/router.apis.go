/**
 * @Author: lzw5399
 * @Date: 2021/3/9 15:52
 * @Desc: 流程定义路由
 */
package router

import (
	"github.com/labstack/echo/v4"

	"workflow/src/controller"
)

// 流程定义
func RegisterProcessDefinition(r *echo.Group) {
	processGroup := r.Group("/process-definitions")
	{
		processGroup.POST("", controller.CreateProcessDefinition)       // 新建
		processGroup.PUT("", controller.UpdateProcessDefinition)        // 修改
		processGroup.DELETE("/:id", controller.DeleteProcessDefinition) // 删除
		processGroup.GET("/:id", controller.GetProcessDefinition)       // 获取流程
		processGroup.GET("", controller.ListProcessDefinition)          // 获取列表
	}
}

// 流程实例
func RegisterProcessInstance(r *echo.Group) {
	instanceGroup := r.Group("/process-instances")
	{
		instanceGroup.POST("", controller.CreateProcessInstance)          // 新建流程
		instanceGroup.GET("/:id", controller.GetProcessInstance)          // 获取
		instanceGroup.GET("", controller.ListProcessInstances)            // 获取列表
		instanceGroup.POST("/_handle", controller.HandleProcessInstance)  // 流程审批
		instanceGroup.POST("/_deny", controller.DenyProcessInstance)      // 流程否决
		instanceGroup.GET("/:id/train-nodes", controller.GetProcessTrain) // 获取流程链路
	}
}

// 外部系统 角色和用户对应关系
func RegisterRoleUsers(r *echo.Group) {
	instanceGroup := r.Group("/role-users")
	{
		instanceGroup.POST("/_batch", controller.BatchSyncRoleUsers) // 批量更新
		instanceGroup.POST("", controller.SyncRoleUsers)             // 单条更新
	}
}
