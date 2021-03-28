/**
 * @Author: lzw5399
 * @Date: 2021/3/22 16:41
 * @Desc: 新建流程实例相关的方法
 */
package engine

import (
	"fmt"

	"workflow/src/model"
)

// 创建实例化相关信息
func (engine *ProcessEngine) CreateProcessInstance() error {
	// 创建
	err := engine.tx.Create(&engine.ProcessInstance).Error
	if err != nil {
		return fmt.Errorf("创建工单失败，%v", err.Error())
	}

	// 创建历史记录
	initialNode, _ := engine.GetInitialNode()
	nextNodes, _ := engine.GetTargetNodes(initialNode)
	nextNode := nextNodes[0] // 开始节点后面只会直连一个节点
	engine.SetCurrentNodeEdgeInfo(&initialNode, nil, &nextNode)
	err = engine.CreateCirculationHistory("")
	if err != nil {
		return fmt.Errorf("新建历史记录失败，%v", err.Error())
	}

	// 更新process_definition表的提交数量统计
	err = engine.tx.Model(&model.ProcessDefinition{}).
		Where("id = ?", engine.ProcessInstance.ProcessDefinitionId).
		Update("submit_count", engine.ProcessDefinition.SubmitCount+1).Error
	if err != nil {
		return fmt.Errorf("更新流程提交数量统计失败，%v", err.Error())
	}

	return nil
}
