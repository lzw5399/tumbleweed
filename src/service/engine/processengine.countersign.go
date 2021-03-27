/**
 * @Author: lzw5399
 * @Date: 2021/3/27 21:59
 * @Desc: 会签相关方法
 */
package engine

import (
	"errors"
	"time"

	"workflow/src/util"
)

// 判断是否是会签，如果是就更新相关状态
func (engine *ProcessEngine) JudgeCounterSign() (isCounterSign bool, isLastProcessor bool, err error) {
	// 判断当前节点是否会签
	isCounterSign = engine.IsCounterSign()
	isLastProcessor = true

	// 不是会签 或者 流程为拒绝的情况下 直接退出
	if !isCounterSign || engine.linkEdge.FlowProperties == "0" {
		return
	}

	// 是最后一个人也退出
	isLastProcessor, err = engine.IsCounterSignLastProcessor()
	if err != nil || isLastProcessor {
		return
	}

	// 不是最后一个审批者则更新相关信息
	err = engine.UpdateInstanceForCounterSign()
	return
}

// 判断是否会签
func (engine *ProcessEngine) IsCounterSign() bool {
	isCounterSign := false
	for _, state := range engine.ProcessInstance.State {
		if state.Id == engine.sourceNode.Id {
			isCounterSign = state.IsCounterSign
			break
		}
	}

	return isCounterSign
}

// 判断当前用户是否是最后一个未审批的
// 并且更新CompletedProcessor字段
func (engine *ProcessEngine) IsCounterSignLastProcessor() (bool, error) {
	isLastPerson := false
	matched := false
	for index, state := range engine.ProcessInstance.State {
		if state.Id == engine.sourceNode.Id {
			// 取差集获取未审批的人
			unCompletedProcessor := util.SliceDiff(state.Processor, state.CompletedProcessor)

			// 判断是否是最后一个未审批的人
			if len(unCompletedProcessor) == 1 && unCompletedProcessor[0] == int(engine.currentUserId) {
				isLastPerson = true
			}

			// 更新CompletedProcessor字段
			engine.ProcessInstance.State[index].CompletedProcessor = append(engine.ProcessInstance.State[index].CompletedProcessor, int(engine.currentUserId))
			matched = true
			break
		}
	}

	if !matched {
		return isLastPerson, errors.New("未找到当前的state，请检查")
	}

	return isLastPerson, nil
}

// 更新会签的流程状态
// 必须是当前会签的角色不为最后一个人，且走的流程不为拒绝流程
func (engine *ProcessEngine) UpdateInstanceForCounterSign() error {
	toUpdate := map[string]interface{}{
		"state":          engine.ProcessInstance.State,
		"update_time":    time.Now().Local(),
		"update_by":      engine.currentUserId,
		"related_person": engine.ProcessInstance.RelatedPerson,
		"variables":      engine.ProcessInstance.Variables,
	}

	err := engine.tx.Model(&engine.ProcessInstance).
		Updates(toUpdate).
		Error

	return err
}
