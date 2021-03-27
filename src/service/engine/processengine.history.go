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
func (engine *ProcessEngine) CreateCirculationHistory(remark string) error {
	// 源节点不为【开始事件】的，获取上一条的流转历史的CreateTime来计算CostDuration
	duration := "0小时 0分钟"
	if engine.sourceNode.Clazz != constant.START {
		var lastCirculation model.CirculationHistory
		err := engine.tx.
			Where("process_instance_id = ?", engine.ProcessInstance.Id).
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
	case engine.sourceNode.Clazz == constant.START:
		sourceState = engine.sourceNode.Label
		sourceId = engine.sourceNode.Id
		targetId = engine.targetNode.Id
		circulation = "新建"
	case engine.sourceNode.Clazz == constant.End:
		sourceState = engine.sourceNode.Label
		sourceId = engine.sourceNode.Id
		targetId = ""
		circulation = "结束"
	default:
		sourceState = engine.sourceNode.Label
		sourceId = engine.sourceNode.Id
		targetId = engine.targetNode.Id
		circulation = engine.linkEdge.Label
	}

	// 创建新的一条流转历史
	cirHistory := model.CirculationHistory{
		AuditableBase: model.AuditableBase{
			CreateBy: engine.currentUserId,
			UpdateBy: engine.currentUserId,
		},
		Title:             engine.ProcessInstance.Title,
		ProcessInstanceId: engine.ProcessInstance.Id,
		SourceState:       sourceState,
		SourceId:          sourceId,
		TargetId:          targetId,
		Circulation:       circulation,
		ProcessorId:       engine.currentUserId,
		CostDuration:      duration,
		Remarks:           remark,
	}

	err := engine.tx.
		Model(&model.CirculationHistory{}).
		Create(&cirHistory).
		Error

	return err
}
