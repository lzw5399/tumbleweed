/**
 * @Author: lzw5399
 * @Date: 2021/3/28 13:49
 * @Desc: 并行网关相关逻辑
 */
package engine

import (
	"errors"

	"workflow/src/model/dto"
)

// 处理并行网关
func (engine *ProcessEngine) ProcessParallelGateway() error {
	gatewayNode := engine.targetNode

	// 获取所有source为当前 网关id 的edge
	nextEdges := engine.GetEdges(gatewayNode.Id, "source")

	// 获取所有target为当前 网关id 的edge
	sourceEdges := engine.GetEdges(gatewayNode.Id, "target")

	// 判断当前是fork还是join
	switch {
	// fork
	case len(sourceEdges) == 1 && len(nextEdges) >= 1:
		return engine.ProcessParallelFork(sourceEdges[0], nextEdges)

	// join
	case len(sourceEdges) >= 1 && len(nextEdges) == 1:
		return engine.ProcessParallelJoin(sourceEdges, nextEdges[0])

	default:
		return errors.New("并行网关流程不正确")
	}

	// 更新当前的
	
}

// 处理并行网关的fork
func (engine *ProcessEngine) ProcessParallelFork(sourceEdge dto.Edge, nextEdges []dto.Edge) error {
	// 先获取并行网关之后，接着的节点
	targetNodes := make([]dto.Node, 0, 1)
	for _, edge := range nextEdges {
		node, err := engine.GetTargetNodeByEdgeId(edge.Id)
		if err != nil {
			return err
		}
		targetNodes = append(targetNodes, node)
	}

	// 根据节点获取state
	states, err := engine.GenStates(targetNodes)
	if err != nil {
		return err
	}

	// 跳转
	err = engine.Circulation(states)

	return err
}

// 处理并行网关的join
func (engine *ProcessEngine) ProcessParallelJoin(sourceEdges []dto.Edge, nextEdge dto.Edge) error {
	return nil
}
