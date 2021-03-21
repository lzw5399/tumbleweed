/**
 * @Author: lzw5399
 * @Date: 2021/3/19 17:10
 * @Desc: 工单的流转相关方法
 */
package engine

import (
	"encoding/json"
	"errors"
	"time"

	"workflow/src/global"
	"workflow/src/global/constant"
	"workflow/src/model"
	"workflow/src/model/request"
	"workflow/src/util"
)

// 一般流转处理，兼顾了会签的判断
func (i *InstanceEngine) CommonProcessing(edge map[string]interface{}, targetNode map[string]interface{}, newStates []map[string]interface{}) error {
	// 如果是拒绝的流程直接跳转
	if edge["flowProperties"] == 0 {
		return i.Circulation(targetNode, newStates)
	}

	// TODO 同意的流程需要判断是否会签

	return i.Circulation(targetNode, newStates)
}

// processInstance流转处理
func (i *InstanceEngine) Circulation(targetNode map[string]interface{}, newStates []map[string]interface{}) error {
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

	state := util.MarshalToDbJson(newStates)

	toUpdate := map[string]interface{}{
		"state":          state,
		"related_person": i.ProcessInstance.RelatedPerson,
		"is_end":         false,
		"update_time":    time.Now().Local(),
		"update_by":      i.currentUserId,
		"variables":      i.ProcessInstance.Variables,
	}

	// 如果是跳转到结束节点，则需要修改节点状态
	if targetNode["clazz"] == constant.End {
		toUpdate["is_end"] = true
	}

	err := global.BankDb.
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

	err := global.BankDb.
		Model(&i.ProcessInstance).
		Updates(toUpdate).
		Error

	// 获取上一条的流转历史的CreateTime来计算CostDuration
	var lastCirculation model.CirculationHistory
	err = global.BankDb.
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
	var currentInstanceState []map[string]interface{}
	err = json.Unmarshal(i.ProcessInstance.State, &currentInstanceState)
	if err != nil {
		return errors.New("当前processInstance的state状态不合法, 请检查")
	}
	// todo 这里先判断[0]
	state := currentInstanceState[0]
	cirHistory := model.CirculationHistory{
		AuditableBase: model.AuditableBase{
			CreateBy: i.currentUserId,
			UpdateBy: i.currentUserId,
		},
		Title:             i.ProcessInstance.Title,
		ProcessInstanceId: i.ProcessInstance.Id,
		SourceState:       state["label"].(string),
		SourceId:          state["id"].(string),
		TargetId:          "",
		Circulation:       "否决",
		ProcessorId:       i.currentUserId,
		CostDuration:      duration,
		Remarks:           r.Remarks,
	}

	err = global.BankDb.
		Model(&model.CirculationHistory{}).
		Create(&cirHistory).
		Error

	if err != nil {
		return err
	}

	return err
}
