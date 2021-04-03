/**
 * @Author: lzw5399
 * @Date: 2021/3/27 22:45
 * @Desc: 验证的相关方法
 */
package engine

import (
	. "github.com/ahmetb/go-linq/v3"

	"workflow/src/model/dto"
	"workflow/src/model/request"
	"workflow/src/util"
)

// 验证入参合法性
func (engine *ProcessEngine) ValidateHandleRequest(r *request.HandleInstancesRequest) error {
	currentEdge, err := engine.GetEdge(r.EdgeId)
	if err != nil {
		return util.BadRequest.New(err)
	}

	state, err := engine.GetStateByEdgeId(currentEdge)
	if err != nil {
		return err
	}

	if currentEdge.Source != state.Id {
		return util.BadRequest.New("当前审批不合法, 请检查")
	}

	// 判断当前流程实例状态是否已结束或者被否决
	if engine.ProcessInstance.IsEnd {
		return util.BadRequest.New("当前流程已结束, 不能进行审批操作")
	}

	if engine.ProcessInstance.IsDenied {
		return util.BadRequest.New("当前流程已被否决, 不能进行审批操作")
	}

	// 判断当前用户是否有权限
	return engine.EnsurePermission(state)
}

// 验证否决请求的入参
func (engine *ProcessEngine) ValidateDenyRequest(r *request.DenyInstanceRequest) error {
	state, err := engine.GetStateByNodeId(r.NodeId)
	if err != nil {
		return util.BadRequest.New(err)
	}

	// 判断当前流程实例状态是否已结束或者被否决
	if engine.ProcessInstance.IsEnd {
		return util.BadRequest.New("当前流程已结束, 不能进行审批操作")
	}

	if engine.ProcessInstance.IsDenied {
		return util.BadRequest.New("当前流程已被否决, 不能进行审批操作")
	}

	// 判断是否有权限
	return engine.EnsurePermission(state)
}

// 判断当前用户是否有权限
func (engine *ProcessEngine) EnsurePermission(state dto.State) error {
	// 判断当前角色是否有权限
	hasPermission := false
	for _, processor := range state.Processor {
		if processor == engine.userIdentifier {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return util.Forbidden.New("当前用户无权限进行当前操作")
	}

	alreadyCompleted := From(state.CompletedProcessor).AnyWith(func(it interface{}) bool {
		return it.(string) == engine.userIdentifier
	})
	if alreadyCompleted {
		return util.Forbidden.New("当前用户针对目前节点已审核, 无法重复审核")
	}

	return nil
}
