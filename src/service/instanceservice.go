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
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"

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
	Get(*request.GetInstanceRequest, uint, uint) (*response.ProcessInstanceResponse, error)
	List(*request.InstanceListRequest, uint, uint) (*response.PagingResponse, error)
	HandleProcessInstance(*request.HandleInstancesRequest, uint, uint) (*model.ProcessInstance, error)
	GetProcessTrain(pi *model.ProcessInstance, instanceId uint, tenantId uint) ([]response.ProcessChainNode, error)
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
		processInstance      = r.ToProcessInstance(currentUserId, tenantId)
		processEngine        *engine.ProcessEngine  // 流程定义引擎
		instanceEngine       *engine.InstanceEngine // 流程实例引擎
		condExprStatus       bool
		sourceEdges          []map[string]interface{}
		targetEdges          []map[string]interface{}
	)

	// 查询对应的流程模板
	err = global.BankDb.
		Where("id = ?", processInstance.ProcessDefinitionId).
		Where("tenant_id = ?", tenantId).
		First(&processDefinition).
		Error
	if err != nil {
		return nil, err
	}

	// 实例化流程引擎
	processEngine, err = engine.NewProcessEngine(processDefinition)
	if err != nil {
		return nil, err
	}

	instanceEngine, err = engine.NewInstanceEngine(processDefinition, currentUserId, tenantId)
	if err != nil {
		return nil, err
	}

	// 将初始状态赋值给当前的流程实例
	if currentInstanceState, err = instanceEngine.GetInstanceInitialState(); err != nil {
		return nil, err
	} else {
		processInstance.State = util.MarshalToBytes(currentInstanceState)
	}
	// TODO 省略了processInstance.State针对变量的预处理

	// 把对应的流程模板的structure单独反序列化出来处理
	comingNode, err := processEngine.GetNode(currentInstanceState[0]["id"].(string))
	if err != nil {
		return nil, err
	}

	switch comingNode["clazz"] {
	// 排他网关
	case "exclusiveGateway":
		var sourceEdges []map[string]interface{}
		sourceEdges, err = processEngine.GetEdge(comingNode["id"].(string), "source")
		if err != nil {
			return nil, err
		}
	breakTag:
		for _, edge := range sourceEdges {
			edgeCondExpr := make([]map[string]interface{}, 0)
			err = json.Unmarshal([]byte(edge["conditionExpression"].(string)), &edgeCondExpr)
			if err != nil {
				return nil, err
			}
			for _, condExpr := range edgeCondExpr {
				// 条件判断
				condExprStatus, err = processEngine.ConditionalJudgment(condExpr)
				if err != nil {
					return nil, err
				}
				if condExprStatus {
					// 进行节点跳转
					comingNode, err = processEngine.GetNode(edge["target"].(string))
					if err != nil {
						return nil, err
					}

					if comingNode["clazz"] == "userTask" || comingNode["clazz"] == "receiveTask" {
						if comingNode["assignValue"] == nil || comingNode["assignType"] == "" {
							err = errors.New("处理人不能为空")
							return nil, err
						}
					}
					currentInstanceState[0]["id"] = comingNode["id"].(string)
					currentInstanceState[0]["label"] = comingNode["label"]
					currentInstanceState[0]["processor"] = comingNode["assignValue"]
					currentInstanceState[0]["process_method"] = comingNode["assignType"]
					break breakTag
				}
			}
		}
		if !condExprStatus {
			return nil, errors.New("所有流转均不符合条件，请确认。")
		}
	// 并行网关
	case "parallelGateway":
		// 入口，判断
		sourceEdges, err = processEngine.GetEdge(comingNode["id"].(string), "source")
		if err != nil {
			return nil, fmt.Errorf("查询流转信息失败，%v", err.Error())
		}

		targetEdges, err = processEngine.GetEdge(comingNode["id"].(string), "target")
		if err != nil {
			return nil, fmt.Errorf("查询流转信息失败，%v", err.Error())
		}

		if len(sourceEdges) > 0 {
			comingNode, err = processEngine.GetNode(sourceEdges[0]["target"].(string))
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("并行网关流程不正确")
		}

		if len(sourceEdges) > 1 && len(targetEdges) == 1 {
			// 入口
			currentInstanceState = []map[string]interface{}{}
			for _, edge := range sourceEdges {
				targetStateValue, err := processEngine.GetNode(edge["target"].(string))
				if err != nil {
					return nil, err
				}
				currentInstanceState = append(currentInstanceState, map[string]interface{}{
					"id":             edge["target"].(string),
					"label":          targetStateValue["label"],
					"processor":      targetStateValue["assignValue"],
					"process_method": targetStateValue["assignType"],
				})
			}
		} else {
			return nil, errors.New("并行网关流程配置不正确")
		}
	}

	// 变量的预处理
	err = preprocessVariables(currentInstanceState, currentUserId)
	if err != nil {
		log.Error(err)
		return nil, errors.New("获取处理人变量值失败")
	}

	// 将最新的“变量/状态信息”赋值给processInstance
	processInstance.State, err = json.Marshal(currentInstanceState)
	if err != nil {
		return nil, err
	}

	// processInstance某些字段更新
	processInstance.RelatedPerson = append(processInstance.RelatedPerson, int64(currentUserId))

	// 开启事务
	err = global.BankDb.Transaction(func(tx *gorm.DB) error {

		// 创建
		err = tx.Create(&processInstance).Error
		if err != nil {
			return fmt.Errorf("创建工单失败，%v", err.Error())
		}

		// todo 省略了【创建工单模版关联数据】

		// todo 省略了【获取当前用户信息】

		// 创建历史记录
		initialNode, _ := processEngine.GetInitialNode()
		err = tx.Create(&model.CirculationHistory{
			AuditableBase: model.AuditableBase{
				CreateBy: currentUserId,
				UpdateBy: currentUserId,
			},
			Title:             processInstance.Title,
			ProcessInstanceId: processInstance.Id,
			SourceState:       initialNode["label"].(string),
			SourceId:          initialNode["id"].(string),
			TargetId:          currentInstanceState[0]["id"].(string),
			Circulation:       "新建",
			ProcessorId:       currentUserId,
			CostDuration:      "0",
			Remarks:           "",
		}).Error
		if err != nil {
			return fmt.Errorf("新建历史记录失败，%v", err.Error())
		}

		// 更新process_definition表的提交数量统计
		err = tx.Model(&model.ProcessDefinition{}).
			Where("id = ?", processInstance.ProcessDefinitionId).
			Update("submit_count", processDefinition.SubmitCount+1).Error
		if err != nil {
			return fmt.Errorf("更新流程提交数量统计失败，%v", err.Error())
		}

		return nil
	})

	// todo 暂时省略了发送通知

	// todo 暂时省略了执行脚本任务

	return &processInstance, err
}

// 获取单个ProcessInstance
func (i *instanceService) Get(r *request.GetInstanceRequest, currentUserId uint, tenantId uint) (*response.ProcessInstanceResponse, error) {
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
func (i *instanceService) List(r *request.InstanceListRequest, currentUserId uint, tenantId uint) (*response.PagingResponse, error) {
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
	)

	// 流程实例引擎
	instanceEngine, err = engine.NewInstanceEngineByInstanceId(r.ProcessInstanceId, currentUserId, tenantId)
	if err != nil {
		return nil, err
	}

	// 验证合法性(1.edgeId是否合法 2.当前用户是否有权限处理)
	err = instanceEngine.ValidateHandleRequest(r, currentUserId)
	if err != nil {
		return nil, err
	}

	// 处理
	err = instanceEngine.Handle(r)

	return &instanceEngine.ProcessInstance, err
}

// 获取流程链(用于展示)
func (i *instanceService) GetProcessTrain(pi *model.ProcessInstance, instanceId uint, tenantId uint) ([]response.ProcessChainNode, error) {
	// 获取流程实例(如果为空)
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

	// 获取流程模板
	var definition model.ProcessDefinition
	err := global.BankDb.
		Where("id=?", instance.ProcessDefinitionId).
		Where("tenant_id = ?", tenantId).
		First(&definition).
		Error
	if err != nil {
		return nil, errors.New("当前流程对应的模板为空")
	}

	// 获取模板结构
	var definitionStructure engine.DefinitionStructure
	err = json.Unmarshal(definition.Structure, &definitionStructure)
	if err != nil {
		return nil, errors.New("流程模板的结构不合法，请检查")
	}

	// 获取实例的当前nodeId
	var currentInstanceState []map[string]interface{}
	err = json.Unmarshal(instance.State, &currentInstanceState)
	if err != nil {
		return nil, errors.New("当前流程实例的状态不合法, 请检查")
	}
	// todo 暂不支持并行网关，所以先用0
	currentNodeId := currentInstanceState[0]["id"].(string)

	shownNodes := make([]map[string]interface{}, 0) // 显示的节点
	currentNodeSortNumber := 0                      // 当前节点的顺序, 为了防止当前节点被隐藏的情况，抽出来
	for _, node := range definitionStructure["nodes"] {
		// 隐藏节点就跳过
		if node["isHideNode"] != nil && node["isHideNode"].(bool) == true {
			continue
		}
		// 获取当前节点的顺序
		if node["id"].(string) == currentNodeId {
			currentNodeSortNumber = util.StringToInt(node["sort"].(string))
		}
		shownNodes = append(shownNodes, node)
	}

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
			Name:     node["label"].(string),
			Id:       node["id"].(string),
			Status:   status,
			Sort:     currentNodeSort,
			NodeType: nodeType,
		}
	}).OrderBy(func(i interface{}) interface{} {
		return i.(response.ProcessChainNode).Sort
	}).ToSlice(&trainNodes)

	return trainNodes, nil
}

// 获取实例的某一个变量
//func (i *instanceService) GetVariable(r *request.GetVariableRequest) (*response.InstanceVariableResponse, error) {
//	var str string
//	err := global.BankDb.Model(&model.ProcessInstance{}).Where("id=?", r.InstanceId).Select("variables").First(&str).Error
//	if err != nil {
//		return nil, errors.New("当前instance实例不存在")
//	}
//
//	var variables []model.InstanceVariable
//	err = json.Unmarshal([]byte(str), &variables)
//	if err != nil {
//		return nil, errors.New("获取失败")
//	}
//
//	for _, v := range variables {
//		if v.Name == r.VariableName {
//			return &response.InstanceVariableResponse{
//				Name:  v.Name,
//				Type:  v.Type,
//				Value: v.Value,
//			}, nil
//		}
//	}
//
//	return nil, errors.New("变量不存在")
//}

//func (i *instanceService) ListVariables(r *request.GetVariableListRequest) (*response.PagingResponse, error) {
//	var str string
//	err := global.BankDb.Model(&model.ProcessInstance{}).Where("id=?", r.InstanceId).Select("variables").First(&str).Error
//	if err != nil {
//		return nil, errors.New("当前instance实例不存在")
//	}
//
//	var variables []model.InstanceVariable
//	err = json.Unmarshal([]byte(str), &variables)
//	if err != nil {
//		return nil, errors.New("获取失败")
//	}
//
//	//var finalVariables []model.InstanceVariable
//	//util.NewPaging(variables).Offset(r.Offset).Limit(r.Limit).Get(&finalVariables)
//
//	return &response.PagingResponse{
//		TotalCount:   int64(len(variables)),
//		CurrentCount: int64(len(variables)),
//		Data:         &variables,
//	}, err
//}

// 创建流程实例的时候预先处理变量转成实际处理人
func preprocessVariables(stateList []map[string]interface{}, creator uint) (err error) {
	//var (
	//	userInfo system.SysUser
	//	deptInfo system.Dept
	//)
	//
	//// 变量转为实际的数据
	//for _, stateItem := range stateList {
	//	if stateItem.(map[string]interface{})["process_method"] == "variable" {
	//		for processorIndex, processor := range stateItem.(map[string]interface{})["processor"].([]interface{}) {
	//			if int(processor.(float64)) == 1 {
	//				// 创建者
	//				stateItem.(map[string]interface{})["processor"].([]interface{})[processorIndex] = creator
	//			} else if int(processor.(float64)) == 2 {
	//				// 1. 查询用户信息
	//				err = orm.Eloquent.Model(&userInfo).Where("user_id = ?", creator).Find(&userInfo).Error
	//				if err != nil {
	//					return
	//				}
	//				// 2. 查询部门信息
	//				err = orm.Eloquent.Model(&deptInfo).Where("dept_id = ?", userInfo.DeptId).Find(&deptInfo).Error
	//				if err != nil {
	//					return
	//				}
	//
	//				// 3. 替换处理人信息
	//				stateItem.(map[string]interface{})["processor"].([]interface{})[processorIndex] = deptInfo.Leader
	//			}
	//		}
	//		stateItem.(map[string]interface{})["process_method"] = "person"
	//	}
	//}

	return
}
