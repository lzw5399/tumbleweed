/**
 * @Author: lzw5399
 * @Date: 2021/1/16 22:58
 * @Desc: 流程实例服务
 */
package service

import (
	"encoding/json"
	"errors"
	"fmt"

	. "github.com/ahmetb/go-linq/v3"

	"workflow/src/global"
	"workflow/src/global/constant"
	"workflow/src/global/shared"
	"workflow/src/model"
	"workflow/src/model/request"
	"workflow/src/model/response"
	"workflow/src/service/engine"
	"workflow/src/util"
)

type InstanceService interface {
	CreateProcessInstance(*request.ProcessInstanceRequest, uint, uint) (*model.ProcessInstance, error)
	GetProcessInstance(*request.GetInstanceRequest, uint, uint) (*response.ProcessInstanceResponse, error)
	ListProcessInstance(*request.InstanceListRequest, uint, uint) (*response.PagingResponse, error)
	HandleProcessInstance(*request.HandleInstancesRequest, uint, uint) (*model.ProcessInstance, error)
	GetProcessTrain(pi *model.ProcessInstance, instanceId uint, tenantId uint) ([]response.ProcessChainNode, error)
	DenyProcessInstance(*request.DenyInstanceRequest, uint, uint) (*model.ProcessInstance, error)
}

type instanceService struct {
}

func NewInstanceService() *instanceService {
	return &instanceService{}
}

// 创建实例
func (i *instanceService) CreateProcessInstance(r *request.ProcessInstanceRequest, currentUserId uint, tenantId uint) (*model.ProcessInstance, error) {
	var (
		currentInstanceState []map[string]interface{} // 变量值
		err                  error
		processDefinition    model.ProcessDefinition // 流程模板
		tx                   = global.BankDb.Begin() // 开启事务
	)

	// 检查变量是否合法
	err = validateVariables(r.Variables)
	if err != nil {
		return nil, err
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
	instanceEngine, err := engine.NewInstanceEngine(processDefinition, r.ToProcessInstance(currentUserId, tenantId), currentUserId, tenantId, tx)
	if err != nil {
		return nil, err
	}

	// 将初始状态赋值给当前的流程实例
	if currentInstanceState, err = instanceEngine.GetInstanceInitialState(); err != nil {
		return nil, err
	} else {
		instanceEngine.ProcessInstance.State = util.MarshalToDbJson(currentInstanceState)
	}

	// TODO 这里判断下一步是排他网关等情况

	// 将最新的“变量/状态信息”赋值给processInstance
	instanceEngine.ProcessInstance.State, err = json.Marshal(currentInstanceState)
	if err != nil {
		return nil, err
	}

	// processInstance某些字段更新
	instanceEngine.ProcessInstance.RelatedPerson = append(instanceEngine.ProcessInstance.RelatedPerson, int64(currentUserId))

	// 创建
	err = instanceEngine.CreateProcessInstance(currentInstanceState)
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return &instanceEngine.ProcessInstance, err
}

// 获取单个ProcessInstance
func (i *instanceService) GetProcessInstance(r *request.GetInstanceRequest, currentUserId uint, tenantId uint) (*response.ProcessInstanceResponse, error) {
	var instance model.ProcessInstance
	err := global.BankDb.
		Where("id=?", r.Id).
		Where("tenant_id = ?", tenantId).
		First(&instance).
		Error
	if err != nil {
		return nil, err
	}

	// 必须是相关的才能看到
	exist := From(instance.RelatedPerson).AnyWith(func(i interface{}) bool {
		return i.(int64) == int64(currentUserId)
	})
	if !exist {
		return nil, errors.New("记录不存在")
	}

	resp := response.ProcessInstanceResponse{
		ProcessInstance: instance,
	}

	// 包括流程链路
	if r.IncludeProcessTrain {
		trainNodes, err := i.GetProcessTrain(&instance, instance.Id, tenantId)
		if err != nil {
			return nil, err
		}
		resp.ProcessChainNodes = trainNodes
	}

	return &resp, nil
}

// 获取ProcessInstance列表
func (i *instanceService) ListProcessInstance(r *request.InstanceListRequest, currentUserId uint, tenantId uint) (*response.PagingResponse, error) {
	var instances []model.ProcessInstance
	db := global.BankDb.Model(&model.ProcessInstance{}).Where("tenant_id = ?", tenantId)

	// 根据type的不同有不同的逻辑
	switch r.Type {
	case constant.I_MyToDo:
		db = db.Joins("cross join jsonb_array_elements(state) as elem").Where(fmt.Sprintf("elem -> 'processor' @> '%v'", currentUserId))
		break
	case constant.I_ICreated:
		db = db.Where("create_by=?", currentUserId)
		break
	case constant.I_IRelated:
		db = db.Where(fmt.Sprintf("related_person @> '%v'", currentUserId))
		break
	case constant.I_All:
		break
	default:
		return nil, errors.New("type不合法")
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
func (i *instanceService) HandleProcessInstance(r *request.HandleInstancesRequest, currentUserId uint, tenantId uint) (*model.ProcessInstance, error) {
	var (
		instanceEngine *engine.InstanceEngine
		err            error
		tx             = global.BankDb.Begin() // 开启事务
	)

	// 验证变量是否符合要求
	err = validateVariables(r.Variables)
	if err != nil {
		return nil, err
	}

	// 流程实例引擎
	instanceEngine, err = engine.NewInstanceEngineByInstanceId(r.ProcessInstanceId, currentUserId, tenantId, tx)
	if err != nil {
		return nil, err
	}

	// 验证合法性(1.edgeId是否合法 2.当前用户是否有权限处理)
	err = instanceEngine.ValidateHandleRequest(r)
	if err != nil {
		return nil, err
	}

	// 合并最新的变量
	instanceEngine.UpdateVariables(r.Variables)

	// 处理操作, 判断这里的原因是因为上面都不会进行数据库改动操作
	err = instanceEngine.Handle(r)
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return &instanceEngine.ProcessInstance, err
}

// 否决流程
func (i *instanceService) DenyProcessInstance(r *request.DenyInstanceRequest, currentUserId uint, tenantId uint) (*model.ProcessInstance, error) {
	var (
		instanceEngine *engine.InstanceEngine
		err            error
		tx             = global.BankDb.Begin() // 开启事务
	)

	// 流程实例引擎
	instanceEngine, err = engine.NewInstanceEngineByInstanceId(r.ProcessInstanceId, currentUserId, tenantId, tx)
	if err != nil {
		return nil, err
	}

	// 验证当前用户是否有权限处理
	err = instanceEngine.ValidateDenyRequest()
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
func (i *instanceService) GetProcessTrain(pi *model.ProcessInstance, instanceId uint, tenantId uint) ([]response.ProcessChainNode, error) {
	// 1. 获取流程实例(如果为空)
	var instance model.ProcessInstance
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

	// 3. 获取模板结构
	var definitionStructure engine.DefinitionStructure
	err = json.Unmarshal(definition.Structure, &definitionStructure)
	if err != nil {
		return nil, errors.New("流程模板的结构不合法，请检查")
	}

	// 4. 获取实例的当前nodeId
	var currentInstanceState []map[string]interface{}
	err = json.Unmarshal(instance.State, &currentInstanceState)
	if err != nil {
		return nil, errors.New("当前流程实例的状态不合法, 请检查")
	}
	// todo 暂不支持并行网关，所以先用0
	currentNodeId := currentInstanceState[0]["id"].(string)

	// 5. 获取所有的显示节点
	shownNodes := make([]map[string]interface{}, 0)
	currentNodeSortNumber := 0 // 当前节点的顺序, 为了防止当前节点被隐藏的情况，抽出来
	initialNodeId := ""
	for _, node := range definitionStructure["nodes"] {
		// 隐藏节点就跳过
		if node["isHideNode"] != nil && node["isHideNode"].(bool) == true {
			continue
		}
		// 获取当前节点的顺序
		if node["id"].(string) == currentNodeId {
			currentNodeSortNumber = util.StringToInt(node["sort"].(string))
		}
		// 找出开始节点的id
		if node["clazz"].(string) == constant.START {
			initialNodeId = node["id"].(string)
		}
		shownNodes = append(shownNodes, node)
	}

	// 6. 遍历出可能的流程链路
	possibleTrainNodesList := make([][]string, 0, util.Pow(len(definitionStructure["nodes"]), 2))
	getPossibleTrainNode(definitionStructure, initialNodeId, []string{}, &possibleTrainNodesList)

	// 7. 遍历获取当前显示的节点是否必须显示的
	// 具体实现方法是遍历possibleTrainNodesList中每一个变量，然后看当前变量的hitCount是否等于len(possibleTrainNodesList)
	// 等于的话，说明在数组每个元素里面都出现了, 那么肯定是必须的
	hitCount := make(map[string]int, len(definitionStructure["nodes"]))
	for _, possibleTrainNodes := range possibleTrainNodesList {
		for _, trainNode := range possibleTrainNodes {
			hitCount[trainNode] = hitCount[trainNode] + 1
		}
	}
	for _, shownNode := range shownNodes {
		shownNode["obligatory"] = hitCount[shownNode["id"].(string)] == len(possibleTrainNodesList)
	}

	// 8. 最后将shownNodes映射成model返回
	var trainNodes []response.ProcessChainNode
	From(shownNodes).Select(func(i interface{}) interface{} {
		node := i.(map[string]interface{})
		currentNodeSort := util.StringToInt(node["sort"].(string))

		var status constant.ChainNodeStatus
		switch {
		case currentNodeSort < currentNodeSortNumber:
			status = 1 // 已处理
		case currentNodeSort > currentNodeSortNumber:
			status = 3 // 未处理的后续节点
		default:
			// 如果是结束节点，则不显示为当前节点，显示为已处理
			if node["clazz"].(string) == constant.End {
				status = 1
			} else { // 其他的等于情况显示为当前节点
				status = 2 // 当前节点
			}
		}

		var nodeType int
		switch node["clazz"].(string) {
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
			Name:       node["label"].(string),
			Id:         node["id"].(string),
			Obligatory: node["obligatory"].(bool),
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
		illegalValueError := fmt.Errorf("当前变量:%s 的类型对应的值不合法，请检查", v.Name)
		// 检查类型
		switch v.Type {
		case constant.VariableNumber:
			_, succeed := v.Value.(float64)
			if !succeed {
				return illegalValueError
			}
		case constant.VariableString:
			_, succeed := v.Value.(string)
			if !succeed {
				return illegalValueError
			}
		case constant.VariableBool:
			_, succeed := v.Value.(bool)
			if !succeed {
				return illegalValueError
			}
		default:
			return fmt.Errorf("当前变量:%s 的类型不合法，请检查", v.Name)
		}

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
func getPossibleTrainNode(definitionStructure engine.DefinitionStructure, currentNodeId string, dependencies []string, possibleTrainNodes *[][]string) {
	targetNodeIds := make([]string, 0)
	// 当前节点添加到依赖中
	dependencies = append(dependencies, currentNodeId)
	for _, edge := range definitionStructure["edges"] {
		// 找到edge的source是当前nodeId的edge
		if edge["source"].(string) == currentNodeId {
			targetNodeIds = append(targetNodeIds, edge["target"].(string))
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
