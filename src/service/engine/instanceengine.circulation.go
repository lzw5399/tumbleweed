/**
 * @Author: lzw5399
 * @Date: 2021/3/19 17:10
 * @Desc: 工单的流转相关方法
 */
package engine

import (
	"time"

	"workflow/src/global/constant"
	"workflow/src/model"
	"workflow/src/model/dto"
	"workflow/src/model/request"
	"workflow/src/util"
)

// 一般流转处理，兼顾了会签的判断
func (i *InstanceEngine) CommonProcessing(newStates dto.StateArray) error {
	// 如果是拒绝的流程直接跳转
	if i.linkEdge.FlowProperties == "0" {
		return i.Circulation(newStates)
	}

	// TODO 暂不支持并行网关，这边先判断0
	state := i.ProcessInstance.State[0]
	// 不是会签
	if !state.IsCounterSign {
		return i.Circulation(newStates)
	}

	// 是会签的最后一个人
	if i.IsCounterSignLastPerson() {
		return i.Circulation(newStates)
	}

	// 不是会签的最后一个人
	//i.tx.Model(&model.ProcessInstance{}).
	return nil
}

// processInstance流转处理
func (i *InstanceEngine) Circulation(newStates dto.StateArray) error {
	// 获取最新的相关者RelatedPerson
	exist := false
	for _, person := range i.ProcessInstance.RelatedPerson {
		if uint(person) == i.currentUserId {
			exist = true
			break
		}
	}
	if !exist {
		i.ProcessInstance.RelatedPerson = append(i.ProcessInstance.RelatedPerson, int64(i.currentUserId))
	}

	toUpdate := map[string]interface{}{
		"state":          newStates,
		"related_person": i.ProcessInstance.RelatedPerson,
		"is_end":         false,
		"update_time":    time.Now().Local(),
		"update_by":      i.currentUserId,
		"variables":      i.ProcessInstance.Variables,
	}

	// 如果是跳转到结束节点，则需要修改节点状态
	if i.targetNode.Clazz == constant.End {
		toUpdate["is_end"] = true
	}

	err := i.tx.
		Model(&i.ProcessInstance).
		Updates(toUpdate).
		Error

	return err
}

// 否决
func (i *InstanceEngine) Deny(r *request.DenyInstanceRequest) error {
	// 获取最新的相关者RelatedPerson
	exist := false
	for _, person := range i.ProcessInstance.RelatedPerson {
		if uint(person) == i.currentUserId {
			exist = true
			break
		}
	}
	if !exist {
		i.ProcessInstance.RelatedPerson = append(i.ProcessInstance.RelatedPerson, int64(i.currentUserId))
	}

	// 更新instance字段
	toUpdate := map[string]interface{}{
		"related_person": i.ProcessInstance.RelatedPerson,
		"is_denied":      true,
		"update_time":    time.Now().Local(),
		"update_by":      i.currentUserId,
	}

	err := i.tx.
		Model(&i.ProcessInstance).
		Updates(toUpdate).
		Error

	// 获取上一条的流转历史的CreateTime来计算CostDuration
	var lastCirculation model.CirculationHistory
	err = i.tx.
		Where("process_instance_id = ?", i.ProcessInstance.Id).
		Order("create_time desc").
		Select("create_time").
		First(&lastCirculation).
		Error
	if err != nil {
		return err
	}
	duration := util.FmtDuration(time.Since(lastCirculation.CreateTime))

	// 创建新的一条流转历史
	// todo 这里先判断[0]
	state := i.ProcessInstance.State[0]
	cirHistory := model.CirculationHistory{
		AuditableBase: model.AuditableBase{
			CreateBy: i.currentUserId,
			UpdateBy: i.currentUserId,
		},
		Title:             i.ProcessInstance.Title,
		ProcessInstanceId: i.ProcessInstance.Id,
		SourceState:       state.Label,
		SourceId:          state.Id,
		TargetId:          "",
		Circulation:       "否决",
		ProcessorId:       i.currentUserId,
		CostDuration:      duration,
		Remarks:           r.Remarks,
	}

	err = i.tx.
		Model(&model.CirculationHistory{}).
		Create(&cirHistory).
		Error

	if err != nil {
		return err
	}

	return err
}

// 是否是会签的最后一个人
func (i *InstanceEngine) IsCounterSignLastPerson() bool {
	return false
}
