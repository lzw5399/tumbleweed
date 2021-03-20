/**
 * @Author: lzw5399
 * @Date: 2021/3/19 17:10
 * @Desc: 工单的流转相关方法
 */
package engine

import (
	"time"

	"workflow/src/global"
	"workflow/src/global/constant"
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
