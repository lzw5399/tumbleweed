/**
 * @Author: lzw5399
 * @Date: 2021/1/16 22:58
 * @Desc: 流程实例服务
 */
package service

import (
	"errors"
	"fmt"

	. "github.com/ahmetb/go-linq/v3"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"workflow/src/global"
	"workflow/src/global/constant"
	"workflow/src/global/shared"
	"workflow/src/model"
	"workflow/src/model/dto"
	"workflow/src/model/request"
	"workflow/src/model/response"
	"workflow/src/service/engine"
	"workflow/src/util"
)

// 创建实例
func CreateProcessInstance(r *request.ProcessInstanceRequest, c echo.Context) (*model.ProcessInstance, error) {
	var (
		processDefinition        model.ProcessDefinition // 流程模板
		tx                       = global.BankDb.Begin() // 开启事务
		tenantId, userIdentifier = util.GetWorkContext(c)
	)

	// 检查变量是否合法
	err := validateVariables(r.Variables)
	if err != nil {
		return nil, util.BadRequest.New(err)
	}

	// 查询对应的流程模板
	err = global.BankDb.
		Where("id = ?", r.ProcessDefinitionId).
		Where("tenant_id = ?", tenantId).
		First(&processDefinition).
		Error
	if err != nil {
		return nil, err
	}

	// 初始化流程引擎
	instanceEngine, err := engine.NewProcessEngine(processDefinition, r.ToProcessInstance(userIdentifier, tenantId), userIdentifier, tenantId, tx)
	if err != nil {
		return nil, err
	}

	// 将初始状态赋值给当前的流程实例
	if currentInstanceState, err := instanceEngine.GetInstanceInitialState(); err != nil {
		return nil, err
	} else {
		instanceEngine.ProcessInstance.State = currentInstanceState
	}

	// TODO 这里判断下一步是排他网关等情况

	// 更新instance的关联人
	instanceEngine.UpdateRelatedPerson()

	// 创建
	err = instanceEngine.CreateProcessInstance()
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return &instanceEngine.ProcessInstance, err
}

// 获取单个ProcessInstance
func GetProcessInstance(r *request.GetInstanceRequest, c echo.Context) (*response.ProcessInstanceResponse, error) {
	var (
		instance                 model.ProcessInstance
		tenantId, userIdentifier = util.GetWorkContext(c)
	)
	err := global.BankDb.
		Where("id=?", r.Id).
		Where("tenant_id = ?", tenantId).
		First(&instance).
		Error
	if err != nil {
		return nil, err
	}

	// 必须是相关的才能看到
	exist := false
	for _, state := range instance.State {
		for _, processor := range state.Processor {
			if processor == userIdentifier {
				exist = true
				break
			}
		}
	}
	if !exist {
		exist = From(instance.RelatedPerson).AnyWith(func(i interface{}) bool {
			return i.(string) == userIdentifier
		})

		if !exist {
			return nil, util.NotFound.New("记录不存在")
		}
	}

	resp := response.ProcessInstanceResponse{
		ProcessInstance: instance,
	}

	// 包括流程链路
	if r.IncludeProcessTrain {
		trainNodes, err := GetProcessTrain(&instance, instance.Id, c)
		if err != nil {
			return nil, err
		}
		resp.ProcessChainNodes = trainNodes
	}

	return &resp, nil
}

// 获取ProcessInstance列表
func ListProcessInstance(r *request.InstanceListRequest, c echo.Context) (*response.PagingResponse, error) {
	var (
		instances                []model.ProcessInstance
		db                       *gorm.DB
		tenantId, userIdentifier = util.GetWorkContext(c)
	)

	switch r.Type {
	case constant.I_MyToDo:
		return getTodoInstances(r, userIdentifier, tenantId)
	default:
		db = global.BankDb.Model(&model.ProcessInstance{}).
			Where("tenant_id = ?", tenantId)
	}

	// 根据type的不同有不同的逻辑
	switch r.Type {
	case constant.I_MyToDo:
		//db = db.Joins("cross join jsonb_array_elements(state) as elem").Where(fmt.Sprintf("elem -> 'processor' @> '[%v]'", userIdentifier))
		break
	case constant.I_ICreated:
		db = db.Where("create_by=?", userIdentifier)
		break
	case constant.I_IRelated:
		db = db.Joins("cross join jsonb_array_elements(state) as elem").
			Where(fmt.Sprintf("related_person @> Array[%d] or elem -> 'processor' @> '[%v]'", userIdentifier, userIdentifier))
		break
	case constant.I_All:
		break
	default:
		return nil, util.BadRequest.New("type不合法")
	}

	if r.Keyword != "" {
		db = db.Where("title ~ ?", r.Keyword)
	}

	var count int64
	db.Count(&count)

	db = shared.ApplyPaging(db, &r.PagingRequest)
	err := db.Find(&instances).Error

	return &response.PagingResponse{
		TotalCount:   count,
		CurrentCount: int64(len(instances)),
		Data:         &instances,
	}, err
}

// 处理/审批ProcessInstance
func HandleProcessInstance(r *request.HandleInstancesRequest, c echo.Context) (*model.ProcessInstance, error) {
	var (
		tx                       = global.BankDb.Begin() // 开启事务
		tenantId, userIdentifier = util.GetWorkContext(c)
	)

	// 验证变量是否符合要求
	err := validateVariables(r.Variables)
	if err != nil {
		return nil, err
	}

	// 流程实例引擎
	processEngine, err := engine.NewProcessEngineByInstanceId(r.ProcessInstanceId, userIdentifier, tenantId, tx)
	if err != nil {
		return nil, err
	}

	// 验证合法性(1.edgeId是否合法 2.当前用户是否有权限处理)
	err = processEngine.ValidateHandleRequest(r)
	if err != nil {
		return nil, err
	}

	// 合并最新的变量
	processEngine.MergeVariables(r.Variables)

	// 处理操作, 判断这里的原因是因为上面都不会进行数据库改动操作
	err = processEngine.Handle(r)
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return &processEngine.ProcessInstance, err
}

// 否决流程
func DenyProcessInstance(r *request.DenyInstanceRequest, c echo.Context) (*model.ProcessInstance, error) {
	var (
		tx                       = global.BankDb.Begin() // 开启事务
		tenantId, userIdentifier = util.GetWorkContext(c)
	)

	// 流程实例引擎
	instanceEngine, err := engine.NewProcessEngineByInstanceId(r.ProcessInstanceId, userIdentifier, tenantId, tx)
	if err != nil {
		return nil, err
	}

	// 验证当前用户是否有权限处理
	err = instanceEngine.ValidateDenyRequest(r)
	if err != nil {
		return nil, err
	}

	// 处理
	err = instanceEngine.Deny(r)
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return &instanceEngine.ProcessInstance, err
}

// 获取流程链(用于展示)
func GetProcessTrain(pi *model.ProcessInstance, instanceId int, c echo.Context) ([]response.ProcessChainNode, error) {
	var (
		instance                 model.ProcessInstance
		tenantId, _ = util.GetWorkContext(c)
	)
	
	// 1. 获取流程实例(如果为空)
	if pi == nil {
		err := global.BankDb.
			Where("id=?", instanceId).
			Where("tenant_id = ?", tenantId).
			First(&instance).
			Error
		if err != nil {
		}
	} else {
		instance = *pi
	}

	// 2. 获取流程模板
	var definition model.ProcessDefinition
	err := global.BankDb.
		Where("id=?", instance.ProcessDefinitionId).
		Where("tenant_id = ?", tenantId).
		First(&definition).
		Error
	if err != nil {
		return nil, errors.New("当前流程对应的模板为空")
	}

	// 3. 获取实例的当前nodeId列表
	currentNodeIds := make([]string, len(instance.State))
	for i, state := range instance.State {
		currentNodeIds[i] = state.Id
	}

	// 4. 获取所有的显示节点
	shownNodes := make([]dto.Node, 0)
	currentNodeSortRange := make([]int, 0) // 当前节点的顺序区间, 在这个区间内的顺序都当作当前节点
	initialNodeId := ""
	for _, node := range definition.Structure.Nodes {
		// 隐藏节点就跳过
		if node.IsHideNode {
			continue
		}
		// 获取当前节点的顺序
		if util.SliceAnyString(currentNodeIds, node.Id) {
			currentNodeSortRange = append(currentNodeSortRange, util.StringToInt(node.Sort))
		}
		// 找出开始节点的id
		if node.Clazz == constant.START {
			initialNodeId = node.Id
		}

		shownNodes = append(shownNodes, node)
	}

	// 5. 遍历出可能的流程链路
	possibleTrainNodesList := make([][]string, 0, util.Pow(len(definition.Structure.Nodes), 2))
	getPossibleTrainNode(definition.Structure, initialNodeId, []string{}, &possibleTrainNodesList)

	// 6. 遍历获取当前显示的节点是否必须显示的
	// 具体实现方法是遍历possibleTrainNodesList中每一个变量，然后看当前变量的hitCount是否等于len(possibleTrainNodesList)
	// 等于的话，说明在数组每个元素里面都出现了, 那么肯定是必须的
	hitCount := make(map[string]int, len(definition.Structure.Nodes))
	for _, possibleTrainNodes := range possibleTrainNodesList {
		for _, trainNode := range possibleTrainNodes {
			hitCount[trainNode] += 1
		}
	}

	// 7. 获取当前节点的排序
	// 由于当前节点可能有多个，排序也相应的有多个，多以会有一个当前节点排序的最大值和最小值
	// 这个范围内圈起来的都被当作当前节点
	currentNodeMinSort, currentNodeMaxSort := util.SliceMinMax(currentNodeSortRange)

	// 8. 最后将shownNodes映射成model返回
	var trainNodes []response.ProcessChainNode
	From(shownNodes).Select(func(i interface{}) interface{} {
		node := i.(dto.Node)
		currentNodeSort := util.StringToInt(node.Sort)

		var status constant.ChainNodeStatus
		switch {
		case currentNodeSort < currentNodeMinSort:
			status = 1 // 已处理
		case currentNodeSort > currentNodeMaxSort:
			status = 3 // 未处理的后续节点
		default:
			// 如果是结束节点，则不显示为当前节点，显示为已处理
			if node.Clazz == constant.End {
				status = 1
			} else { // 其他的等于情况显示为当前节点
				status = 2 // 当前节点
			}
		}

		var nodeType int
		switch node.Clazz {
		case constant.START:
			nodeType = 1
		case constant.UserTask:
			nodeType = 2
		case constant.ExclusiveGateway:
			nodeType = 3
		case constant.End:
			nodeType = 4
		}

		return response.ProcessChainNode{
			Name:       node.Label,
			Id:         node.Id,
			Obligatory: hitCount[node.Id] == len(possibleTrainNodesList),
			Status:     status,
			Sort:       currentNodeSort,
			NodeType:   nodeType,
		}
	}).OrderBy(func(i interface{}) interface{} {
		return i.(response.ProcessChainNode).Sort
	}).ToSlice(&trainNodes)

	return trainNodes, nil
}

// 检查变量是否合法
func validateVariables(variables []model.InstanceVariable) error {
	checkedVariables := make(map[string]model.InstanceVariable, 0)
	for _, v := range variables {
		// 检查是否重名
		if _, present := checkedVariables[v.Name]; present {
			return fmt.Errorf("当前变量名:%s 重复, 请检查", v.Name)
		}
		checkedVariables[v.Name] = v
	}

	return nil
}

// 获取所有的可能的流程链路
// definitionStructure: 流程模板的结构
// currentNodes: 当前需要遍历的节点
// dependencies: 依赖项
// possibleTrainNodes: 最终返回的可能的流程链路
func getPossibleTrainNode(definitionStructure dto.Structure, currentNodeId string, dependencies []string, possibleTrainNodes *[][]string) {
	targetNodeIds := make([]string, 0)
	// 当前节点添加到依赖中
	dependencies = append(dependencies, currentNodeId)
	for _, edge := range definitionStructure.Edges {
		// 找到edge的source是当前nodeId的edge
		if edge.Source == currentNodeId && edge.FlowProperties != "0" {
			targetNodeIds = append(targetNodeIds, edge.Target)
		}
	}

	// 已经是最终节点了
	if len(targetNodeIds) == 0 {
		*possibleTrainNodes = append(*possibleTrainNodes, dependencies)
	} else {
		// 不是最终节点，继续递归遍历
		for _, targetNodeId := range targetNodeIds {
			getPossibleTrainNode(definitionStructure, targetNodeId, dependencies, possibleTrainNodes)
		}
	}
}

func getTodoInstances(r *request.InstanceListRequest, userIdentifier string, tenantId int) (*response.PagingResponse, error) {
	countSql := fmt.Sprintf(`with base as (
		select *,
			jsonb_array_elements(state) as singleState
			from wf.process_instance
			where tenant_id = %d
			AND is_end = false
			AND is_denied = false
			)
			select count(1)
				from base
				where singleState -> 'processor' @> '[%d]'
				and singleState -> 'unCompletedProcessor' @> '[%d]'`,
		tenantId, userIdentifier, userIdentifier)
	if r.Keyword != "" {
		countSql += fmt.Sprintf(" AND title ~ '%s'", r.Keyword)
	}
	var c int64
	err := global.BankDb.Raw(countSql).Scan(&c).Error
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf(`with base as (
		select *,
			jsonb_array_elements(state) as singleState
			from wf.process_instance
			where tenant_id = %d
			AND is_end = false
			AND is_denied = false
			)
			select id, create_time, update_time, create_by, update_by, title, priority,
				process_definition_id, classify_id, is_end, is_denied, state, related_person, tenant_id, variables
				from base
				where singleState -> 'processor' @> '[%d]'
				and singleState -> 'unCompletedProcessor' @> '[%d]' `,
		tenantId, userIdentifier, userIdentifier)
	if r.Keyword != "" {
		sql += fmt.Sprintf(" AND title ~ '%s'", r.Keyword)
	}
	sql = shared.ApplyRawPaging(sql, &r.PagingRequest)
	var instances []model.ProcessInstance
	err = global.BankDb.Raw(sql).Scan(&instances).Error

	return &response.PagingResponse{
		TotalCount:   c,
		CurrentCount: int64(len(instances)),
		Data:         &instances,
	}, err
}
