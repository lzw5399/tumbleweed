/**
 * @Author: lzw5399
 * @Date: 2021/3/21 21:34
 * @Desc: 排他网关的相关方法
 */
package engine

import (
	"errors"
	"fmt"
	"strings"

	"workflow/src/global"
	"workflow/src/global/constant"
	"workflow/src/model/dto"
	"workflow/src/model/request"
	"workflow/src/util"
)

// 处理排他网关的跳转
func (engine *ProcessEngine) ProcessingExclusiveGateway(gatewayNode dto.Node, r *request.HandleInstancesRequest) error {
	// 1. 找到所有source为当前网关节点的edges, 并按照sort排序
	edges := engine.GetEdges(gatewayNode.Id, "source")

	// 2. 遍历edges, 获取当前第一个符合条件的edge
	hitEdge := dto.Edge{}
	for _, edge := range edges {
		if edge.ConditionExpression == "" {
			return errors.New("处理失败, 排他网关的后续流程的条件表达式不能为空, 请检查")
		}

		// 进行条件判断
		condExprStatus, err := engine.ConditionJudgment(edge.ConditionExpression)
		if err != nil {
			return err
		}
		// 获取成功的节点
		if condExprStatus {
			hitEdge = edge
			break
		}
	}

	if hitEdge.Id == "" {
		return errors.New("没有符合条件的流向，请检查")
	}

	// 3. 获取必要的信息
	newTargetNode, err := engine.GetTargetNodeByEdgeId(hitEdge.Id)
	if err != nil {
		return errors.New("模板结构错误")
	}

	// 4. 合并获得最新的states
	var removeStateId string
	if engine.sourceNode.Clazz == constant.ExclusiveGateway || engine.sourceNode.Clazz == constant.ParallelGateway {
		removeStateId = engine.targetNode.Id
	} else {
		removeStateId = engine.sourceNode.Id
	}

	newStates, err := engine.MergeStates(removeStateId, []dto.Node{newTargetNode})
	if err != nil {
		return err
	}

	// 5. 更新最新的node edge等信息
	engine.SetCurrentNodeEdgeInfo(&gatewayNode, &hitEdge, &newTargetNode)

	// 6. 根据edge进行跳转
	err = engine.Circulation(newStates)
	if err != nil {
		return err
	}

	return nil
}

// 条件表达式判断
func (engine *ProcessEngine) ConditionJudgment(condExpr string) (bool, error) {
	// 先获取变量列表
	variables := util.UnmarshalToInstanceVariables(engine.ProcessInstance.Variables)

	envMap := make(map[string]interface{}, len(variables))
	for _, variable := range variables {
		envMap[variable.Name] = variable.Value
	}

	// 替换变量表达式符
	condExpr = strings.Replace(condExpr, "{{", "", -1)
	condExpr = strings.Replace(condExpr, "}}", "", -1)
	condExpr = strings.Replace(condExpr, "&gt;", ">", -1)
	condExpr = strings.Replace(condExpr, "&lt;", "<", -1)

	result, err := util.CalculateExpression(condExpr, envMap)
	if err != nil {
		err = fmt.Errorf("计算表达式发生错误, 当前表达式：%s ,当前变量:%v, 错误原因：%s", condExpr, envMap, err.Error())
		global.BankLogger.Error(err)
		return false, err
	}

	return result, nil
}
