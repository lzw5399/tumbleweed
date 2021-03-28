/**
 * @Author: lzw5399
 * @Date: 2021/3/28 13:49
 * @Desc: 并行网关相关逻辑
 */
package engine

import (
	"errors"
	"time"

	"workflow/src/model/dto"
)

// 处理并行网关
func (engine *ProcessEngine) ProcessParallelGateway() ([]dto.RelationInfo, error) {
	gatewayNode := engine.targetNode

	// 获取所有source为当前 网关id 的edge
	nextEdges := engine.GetEdges(gatewayNode.Id, "source")

	// 获取所有target为当前 网关id 的edge
	sourceEdges := engine.GetEdges(gatewayNode.Id, "target")

	// 判断当前是fork还是join
	switch {
	// fork
	case len(sourceEdges) == 1 && len(nextEdges) >= 1:
		return engine.ProcessParallelFork(*gatewayNode, nextEdges)

	// join
	case len(sourceEdges) >= 1 && len(nextEdges) == 1:
		return engine.ProcessParallelJoin(*gatewayNode, sourceEdges, nextEdges[0])

	default:
		return nil, errors.New("并行网关流程不正确")
	}
}

// 处理并行网关的fork
func (engine *ProcessEngine) ProcessParallelFork(gatewayNode dto.Node, nextEdges []dto.Edge) ([]dto.RelationInfo, error) {
	infos := make([]dto.RelationInfo, 0, 1)

	// 先获取并行网关之后，接着的节点列表
	targetNodes := make([]dto.Node, 0, 1)
	for _, edge := range nextEdges {
		targetNode, err := engine.GetTargetNodeByEdgeId(edge.Id)
		if err != nil {
			return nil, err
		}
		targetNodes = append(targetNodes, targetNode)

		infos = append(infos, dto.RelationInfo{
			SourceNode: gatewayNode,
			LinkedEdge: edge,
			TargetNode: targetNode,
		})
	}

	// 根据节点获取state
	newStates, err := engine.MergeStates(engine.sourceNode.Id, targetNodes)
	if err != nil {
		return nil, err
	}

	// 记录跳转
	err = engine.Circulation(newStates)
	if err != nil {
		return nil, err
	}

	return infos, err
}

// 处理并行网关的join
func (engine *ProcessEngine) ProcessParallelJoin(gatewayNode dto.Node, sourceEdges []dto.Edge, nextEdge dto.Edge) ([]dto.RelationInfo, error) {
	infos := make([]dto.RelationInfo, 0)

	// 获取当前ProcessInstance得state中，是gatewayNode前一个的个数
	gatewayPreviousNodes := engine.GetNodesByEdges(sourceEdges, "source")
	count := 0
	for _, state := range engine.ProcessInstance.State {
		for _, node := range gatewayPreviousNodes {
			if state.Id == node.Id {
				count++
			}
		}
	}

	switch {
	// 大于1，说明还有其他线没有处理完, 不跳转, state数组中去掉当前state即可
	case count > 1:
		// 获取合并后的states
		mergedStates, err := engine.MergeStates(engine.sourceNode.Id, []dto.Node{})
		if err != nil {
			return nil, err
		}
		// 更新到数据库
		err = engine.UpdateInstanceStateForParallel(mergedStates)
		if err != nil {
			return nil, err
		}

	// 等于1, 说明我其他的线都处理完了，可以跳到【当前并行网关】的【下一节点】
	case count == 1:
		// 获取最新的targetNode
		newTargetNode, err := engine.GetNode(nextEdge.Target)
		if err != nil {
			return nil, err
		}

		// 获取合并后的states
		mergedStates, err := engine.MergeStates(engine.sourceNode.Id, []dto.Node{newTargetNode})
		if err != nil {
			return nil, err
		}

		// 更新跳转信息
		err = engine.Circulation(mergedStates)
		if err != nil {
			return nil, err
		}

		// 添加信息(跳到外层进行递归)
		infos = append(infos, dto.RelationInfo{
			SourceNode: gatewayNode,
			LinkedEdge: nextEdge,
			TargetNode: newTargetNode,
		})

	default:
		return nil, errors.New("当前流程结构不合法, 请检查")
	}

	return infos, nil
}

func (engine *ProcessEngine) UpdateInstanceStateForParallel(mergedStates dto.StateArray) error {
	toUpdate := map[string]interface{}{
		"state":          mergedStates,
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
