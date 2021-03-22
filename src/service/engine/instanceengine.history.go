/**
 * @Author: lzw5399
 * @Date: 2021/3/22 11:27
 * @Desc: 流转历史CirculationHistory的相关方法
 */
package engine

import (
	"time"

	"workflow/src/global/constant"
	"workflow/src/model"
	"workflow/src/util"
)

// 创建流转历史记录
func (i *InstanceEngine) CreateCirculationHistory(remark string) error {
	// 源节点不为【开始事件】的，获取上一条的流转历史的CreateTime来计算CostDuration
	duration := "0小时 0分钟"
	if i.sourceNode["clazz"].(string) != constant.START {
		var lastCirculation model.CirculationHistory
		err := i.tx.
			Where("process_instance_id = ?", i.ProcessInstance.Id).
			Order("create_time desc").
			Select("create_time").
			First(&lastCirculation).
			Error
		if err != nil {
			return err
		}
		duration = util.FmtDuration(time.Since(lastCirculation.CreateTime))
	}

	// 根据不同的类型取不同的值
	var sourceState, sourceId, targetId, circulation string
	switch {
	case i.sourceNode["clazz"].(string) == constant.START:
		sourceState = i.sourceNode["label"].(string)
		sourceId = i.sourceNode["id"].(string)
		targetId = i.targetNode["id"].(string)
		circulation = "新建"
	case i.sourceNode["clazz"].(string) == constant.End:
		sourceState = i.sourceNode["label"].(string)
		sourceId = i.sourceNode["id"].(string)
		targetId = ""
		circulation = "结束"
	default:
		sourceState = i.sourceNode["label"].(string)
		sourceId = i.sourceNode["id"].(string)
		targetId = i.targetNode["id"].(string)
		circulation = i.linkEdge["label"].(string)
	}

	// 创建新的一条流转历史
	cirHistory := model.CirculationHistory{
		AuditableBase: model.AuditableBase{
			CreateBy: i.currentUserId,
			UpdateBy: i.currentUserId,
		},
		Title:             i.ProcessInstance.Title,
		ProcessInstanceId: i.ProcessInstance.Id,
		SourceState:       sourceState,
		SourceId:          sourceId,
		TargetId:          targetId,
		Circulation:       circulation,
		ProcessorId:       i.currentUserId,
		CostDuration:      duration,
		Remarks:           remark,
	}

	err := i.tx.
		Model(&model.CirculationHistory{}).
		Create(&cirHistory).
		Error

	return err
}
