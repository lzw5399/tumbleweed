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
func (i *InstanceEngine) CreateProcessInstance(currentInstanceState []map[string]interface{}) error {
	// 创建
	err := i.tx.Create(&i.ProcessInstance).Error
	if err != nil {
		return fmt.Errorf("创建工单失败，%v", err.Error())
	}

	// 创建历史记录
	initialNode, _ := i.GetInitialNode()
	i.SetNodeEdgeInfo(initialNode, nil, currentInstanceState[0])
	err = i.CreateCirculationHistory("")
	if err != nil {
		return fmt.Errorf("新建历史记录失败，%v", err.Error())
	}

	// 更新process_definition表的提交数量统计
	err = i.tx.Model(&model.ProcessDefinition{}).
		Where("id = ?", i.ProcessInstance.ProcessDefinitionId).
		Update("submit_count", i.ProcessDefinition.SubmitCount+1).Error
	if err != nil {
		return fmt.Errorf("更新流程提交数量统计失败，%v", err.Error())
	}

	return nil
}
