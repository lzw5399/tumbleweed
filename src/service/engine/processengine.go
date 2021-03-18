/**
 * @Author: lzw5399
 * @Date: 2021/3/10 18:59
 * @Desc: 读取模板节点数据
 */
package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"gorm.io/gorm"

	"workflow/src/model"
)

type ProcessEngine struct {
	CirHistoryList      []model.CirculationHistory
	WorkOrderId         int
	UpdateValue         map[string]interface{}
	StateValue          map[string]interface{}
	TargetStateValue    map[string]interface{}
	WorkOrderData       [][]byte
	ProcessInstance     model.ProcessInstance // 流程实例
	EndHistory          bool
	FlowProperties      int
	CirculationValue    string
	DefinitionStructure DefinitionStructure // 流程模板结构
	tx                  *gorm.DB
}

type DefinitionStructure map[string][]map[string]interface{}

func NewProcessEngine(d model.ProcessDefinition) (*ProcessEngine, error) {
	var definitionStructure DefinitionStructure
	err := json.Unmarshal(d.Structure, &definitionStructure)
	if err != nil {
		return nil, err
	}

	return &ProcessEngine{
		DefinitionStructure: definitionStructure,
	}, nil
}

// 获取节点信息
func (p *ProcessEngine) GetNode(stateId string) (nodeValue map[string]interface{}, err error) {
	for _, node := range p.DefinitionStructure["nodes"] {
		if node["id"] == stateId {
			nodeValue = node
			return
		}
	}
	return
}

// 获取流转信息
func (p *ProcessEngine) GetEdge(stateId string, classify string) (edgeValue []map[string]interface{}, err error) {
	var (
		leftSort  int
		rightSort int
	)

	for _, edge := range p.DefinitionStructure["edges"] {
		if edge[classify] == stateId {
			edgeValue = append(edgeValue, edge)
		}
	}

	// 排序
	if len(edgeValue) > 1 {
		for i := 0; i < len(edgeValue)-1; i++ {
			for j := i + 1; j < len(edgeValue); j++ {
				if t, ok := edgeValue[i]["sort"]; ok {
					leftSort, _ = strconv.Atoi(t.(string))
				}
				if t, ok := edgeValue[j]["sort"]; ok {
					rightSort, _ = strconv.Atoi(t.(string))
				}
				if leftSort > rightSort {
					edgeValue[j], edgeValue[i] = edgeValue[i], edgeValue[j]
				}
			}
		}
	}

	return
}

// 时间格式化
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute

	return fmt.Sprintf("%02d小时 %02d分钟", h, m)
}

// 会签
//func (h *ProcessEngine) Countersign(c echo.Context) (err error) {
//	var (
//		stateList       []map[string]interface{}
//		stateIdMap      map[string]interface{}
//		currentState    map[string]interface{}
//		cirHistoryCount int
//		//userInfoList      []system.SysUser
//		circulationStatus bool
//	)
//
//	err = json.Unmarshal(h.ProcessInstance.State, &stateList)
//	if err != nil {
//		return
//	}
//
//	stateIdMap = make(map[string]interface{})
//	for _, v := range stateList {
//		stateIdMap[v["id"].(string)] = v["label"]
//		if v["id"].(string) == h.StateValue["id"].(string) {
//			currentState = v
//		}
//	}
//	userStatusCount := 0
//	circulationStatus = false
//	for _, cirHistoryValue := range h.CirHistoryList {
//		if len(currentState["processor"].([]interface{})) > 1 {
//			if _, ok := stateIdMap[cirHistoryValue.Source]; !ok {
//				break
//			}
//		}
//
//		if currentState["process_method"].(string) == "person" {
//			// 用户会签
//			for _, processor := range currentState["processor"].([]interface{}) {
//				if cirHistoryValue.ProcessorId != int(util.GetCurrentUserId(c)) &&
//					cirHistoryValue.Source == currentState["id"].(string) &&
//					cirHistoryValue.ProcessorId == int(processor.(float64)) {
//					cirHistoryCount += 1
//				}
//			}
//			if cirHistoryCount == len(currentState["processor"].([]interface{}))-1 {
//				circulationStatus = true
//				break
//			}
//		} else if currentState["process_method"].(string) == "role" || currentState["process_method"].(string) == "department" {
//			// 全员处理
//			var tmpUserList []system.SysUser
//			if h.StateValue["fullHandle"].(bool) {
//				db := orm.Eloquent.Model(&system.SysUser{})
//				if currentState["process_method"].(string) == "role" {
//					db = db.Where("role_id in (?)", currentState["processor"].([]interface{}))
//				} else if currentState["process_method"].(string) == "department" {
//					db = db.Where("dept_id in (?)", currentState["processor"].([]interface{}))
//				}
//				err = db.Find(&userInfoList).Error
//				if err != nil {
//					return
//				}
//				temp := map[string]struct{}{}
//				for _, user := range userInfoList {
//					if _, ok := temp[user.Username]; !ok {
//						temp[user.Username] = struct{}{}
//						tmpUserList = append(tmpUserList, user)
//					}
//				}
//				for _, user := range tmpUserList {
//					if cirHistoryValue.Source == currentState["id"].(string) &&
//						cirHistoryValue.ProcessorId != tools.GetUserId(c) &&
//						cirHistoryValue.ProcessorId == user.UserId {
//						userStatusCount += 1
//						break
//					}
//				}
//			} else {
//				// 普通会签
//				for _, processor := range currentState["processor"].([]interface{}) {
//					db := orm.Eloquent.Model(&system.SysUser{})
//					if currentState["process_method"].(string) == "role" {
//						db = db.Where("role_id = ?", processor)
//					} else if currentState["process_method"].(string) == "department" {
//						db = db.Where("dept_id = ?", processor)
//					}
//					err = db.Find(&userInfoList).Error
//					if err != nil {
//						return
//					}
//					for _, user := range userInfoList {
//						if user.UserId != tools.GetUserId(c) &&
//							cirHistoryValue.Source == currentState["id"].(string) &&
//							cirHistoryValue.ProcessorId == user.UserId {
//							userStatusCount += 1
//							break
//						}
//					}
//				}
//			}
//			if h.StateValue["fullHandle"].(bool) {
//				if userStatusCount == len(tmpUserList)-1 {
//					circulationStatus = true
//				}
//			} else {
//				if userStatusCount == len(currentState["processor"].([]interface{}))-1 {
//					circulationStatus = true
//				}
//			}
//		}
//	}
//	if circulationStatus {
//		h.EndHistory = true
//		err = h.circulation()
//		if err != nil {
//			return
//		}
//	}
//	return
//}

// 工单跳转
//func (h *ProcessEngine) circulation() (err error) {
//	var (
//		StateValue []byte
//	)
//
//	stateList := make([]interface{}, 0)
//	for _, v := range h.UpdateValue["state"].([]map[string]interface{}) {
//		stateList = append(stateList, v)
//	}
//	err = GetVariableValue(stateList, h.ProcessInstance.Creator)
//	if err != nil {
//		return
//	}
//
//	StateValue, err = json.Marshal(h.UpdateValue["state"])
//	if err != nil {
//		return
//	}
//
//	err = h.tx.Model(&process.WorkOrderInfo{}).
//		Where("id = ?", h.WorkOrderId).
//		Updates(map[string]interface{}{
//			"state":          StateValue,
//			"related_person": h.UpdateValue["related_person"],
//		}).Error
//	if err != nil {
//		h.tx.Rollback()
//		return
//	}
//
//	// 如果是跳转到结束节点，则需要修改节点状态
//	if h.TargetStateValue["clazz"] == "end" {
//		err = h.tx.Model(&process.WorkOrderInfo{}).
//			Where("id = ?", h.WorkOrderId).
//			Update("is_end", 1).Error
//		if err != nil {
//			h.tx.Rollback()
//			return
//		}
//	}
//
//	return
//}

// 条件判断
func (p *ProcessEngine) ConditionalJudgment(condExpr map[string]interface{}) (result bool, err error) {
	var (
		condExprOk    bool
		condExprValue interface{}
	)

	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case string:
				err = errors.New(e)
			case error:
				err = e
			default:
				err = errors.New("未知错误")
			}
			return
		}
	}()

	for _, data := range p.WorkOrderData {
		var formData map[string]interface{}
		err = json.Unmarshal(data, &formData)
		if err != nil {
			return
		}
		if condExprValue, condExprOk = formData[condExpr["key"].(string)]; condExprOk {
			break
		}
	}

	if condExprValue == nil {
		err = errors.New("未查询到对应的表单数据。")
		return
	}

	// todo 待优化
	switch reflect.TypeOf(condExprValue).String() {
	case "string":
		switch condExpr["sign"] {
		case "==":
			if condExprValue.(string) == condExpr["value"].(string) {
				result = true
			}
		case "!=":
			if condExprValue.(string) != condExpr["value"].(string) {
				result = true
			}
		case ">":
			if condExprValue.(string) > condExpr["value"].(string) {
				result = true
			}
		case ">=":
			if condExprValue.(string) >= condExpr["value"].(string) {
				result = true
			}
		case "<":
			if condExprValue.(string) < condExpr["value"].(string) {
				result = true
			}
		case "<=":
			if condExprValue.(string) <= condExpr["value"].(string) {
				result = true
			}
		default:
			err = errors.New("目前仅支持6种常规判断类型，包括（等于、不等于、大于、大于等于、小于、小于等于）")
		}
	case "float64":
		switch condExpr["sign"] {
		case "==":
			if condExprValue.(float64) == condExpr["value"].(float64) {
				result = true
			}
		case "!=":
			if condExprValue.(float64) != condExpr["value"].(float64) {
				result = true
			}
		case ">":
			if condExprValue.(float64) > condExpr["value"].(float64) {
				result = true
			}
		case ">=":
			if condExprValue.(float64) >= condExpr["value"].(float64) {
				result = true
			}
		case "<":
			if condExprValue.(float64) < condExpr["value"].(float64) {
				result = true
			}
		case "<=":
			if condExprValue.(float64) <= condExpr["value"].(float64) {
				result = true
			}
		default:
			err = errors.New("目前仅支持6种常规判断类型，包括（等于、不等于、大于、大于等于、小于、小于等于）")
		}
	default:
		err = errors.New("条件判断目前仅支持字符串、整型。")
	}

	return
}

// 并行网关，确认其他节点是否完成
func (p *ProcessEngine) completeAllParallel(target string) (statusOk bool, err error) {
	var (
		stateList []map[string]interface{}
	)

	err = json.Unmarshal(p.ProcessInstance.State, &stateList)
	if err != nil {
		err = fmt.Errorf("反序列化失败，%v", err.Error())
		return
	}

continueHistoryTag:
	for _, v := range p.CirHistoryList {
		status := false
		for i, s := range stateList {
			if v.Source == s["id"].(string) && v.Target == target {
				status = true
				stateList = append(stateList[:i], stateList[i+1:]...)
				continue continueHistoryTag
			}
		}
		if !status {
			break
		}
	}

	if len(stateList) == 1 && stateList[0]["id"].(string) == p.StateValue["id"] {
		statusOk = true
	}

	return
}

//func (h *ProcessEngine) commonProcessing(c echo.Context) (err error) {
//	// 如果是拒绝的流转则直接跳转
//	if h.FlowProperties == 0 {
//		err = h.circulation()
//		if err != nil {
//			err = fmt.Errorf("工单跳转失败，%v", err.Error())
//		}
//		return
//	}
//
//	// 会签
//	if h.StateValue["assignValue"] != nil && len(h.StateValue["assignValue"].([]interface{})) > 0 {
//		if isCounterSign, ok := h.StateValue["isCounterSign"]; ok {
//			if isCounterSign.(bool) {
//				h.EndHistory = false
//				err = h.Countersign(c)
//				if err != nil {
//					return
//				}
//			} else {
//				err = h.circulation()
//				if err != nil {
//					return
//				}
//			}
//		} else {
//			err = h.circulation()
//			if err != nil {
//				return
//			}
//		}
//	} else {
//		err = h.circulation()
//		if err != nil {
//			return
//		}
//	}
//	return
//}

//func (h *ProcessEngine) HandleWorkOrder(
//	c echo.Context,
//	WorkOrderId int,
//	tasks []string,
//	targetState string,
//	sourceState string,
//	CirculationValue string,
//	FlowProperties int,
//	remarks string,
//	tpls []map[string]interface{},
//) (err error) {
//	h.WorkOrderId = WorkOrderId
//	h.FlowProperties = FlowProperties
//	h.EndHistory = true
//
//	var (
//		execTasks          []string
//		relatedPersonList  []int
//		cirHistoryValue    []process.CirculationHistory
//		cirHistoryData     process.CirculationHistory
//		costDurationValue  string
//		sourceEdges        []map[string]interface{}
//		targetEdges        []map[string]interface{}
//		condExprStatus     bool
//		relatedPersonValue []byte
//		parallelStatusOk   bool
//		processInfo        process.Info
//		currentUserInfo    system.SysUser
//		sendToUserList     []system.SysUser
//		noticeList         []int
//		sendSubject        string = "您有一条待办工单，请及时处理"
//		sendDescription    string = "您有一条待办工单请及时处理，工单描述如下"
//		paramsValue        struct {
//			Id       int           `json:"id"`
//			Title    string        `json:"title"`
//			Priority int           `json:"priority"`
//			FormData []interface{} `json:"form_data"`
//		}
//	)
//
//	defer func() {
//		if r := recover(); r != nil {
//			switch e := r.(type) {
//			case string:
//				err = errors.New(e)
//			case error:
//				err = e
//			default:
//				err = errors.New("未知错误")
//			}
//			return
//		}
//	}()
//
//	// 获取工单信息
//	err = orm.Eloquent.Model(&process.WorkOrderInfo{}).Where("id = ?", WorkOrderId).Find(&h.workOrderDetails).Error
//	if err != nil {
//		return
//	}
//
//	// 获取流程信息
//	err = orm.Eloquent.Model(&process.Info{}).Where("id = ?", h.workOrderDetails.Process).Find(&processInfo).Error
//	if err != nil {
//		return
//	}
//	err = json.Unmarshal(processInfo.Structure, &h.DefinitionStructure.Structure)
//	if err != nil {
//		return
//	}
//
//	// 获取当前节点
//	h.StateValue, err = h.DefinitionStructure.GetNode(sourceState)
//	if err != nil {
//		return
//	}
//
//	// 目标状态
//	h.TargetStateValue, err = h.DefinitionStructure.GetNode(targetState)
//	if err != nil {
//		return
//	}
//
//	// 获取工单数据
//	err = orm.Eloquent.Model(&process.TplData{}).
//		Where("work_order = ?", WorkOrderId).
//		Pluck("form_data", &h.WorkOrderData).Error
//	if err != nil {
//		return
//	}
//
//	// 根据处理人查询出需要会签的条数
//	err = orm.Eloquent.Model(&process.CirculationHistory{}).
//		Where("work_order = ?", WorkOrderId).
//		Order("id desc").
//		Find(&h.CirHistoryList).Error
//	if err != nil {
//		return
//	}
//
//	err = json.Unmarshal(h.workOrderDetails.RelatedPerson, &relatedPersonList)
//	if err != nil {
//		return
//	}
//	relatedPersonStatus := false
//	for _, r := range relatedPersonList {
//		if r == tools.GetUserId(c) {
//			relatedPersonStatus = true
//			break
//		}
//	}
//	if !relatedPersonStatus {
//		relatedPersonList = append(relatedPersonList, tools.GetUserId(c))
//	}
//
//	relatedPersonValue, err = json.Marshal(relatedPersonList)
//	if err != nil {
//		return
//	}
//
//	h.UpdateValue = map[string]interface{}{
//		"related_person": relatedPersonValue,
//	}
//
//	// 开启事务
//	h.tx = orm.Eloquent.Begin()
//
//	StateValue := map[string]interface{}{
//		"label": h.TargetStateValue["label"].(string),
//		"id":    h.TargetStateValue["id"].(string),
//	}
//
//	switch h.TargetStateValue["clazz"] {
//	// 排他网关
//	case "exclusiveGateway":
//		sourceEdges, err = h.DefinitionStructure.GetEdge(h.TargetStateValue["id"].(string), "source")
//		if err != nil {
//			return
//		}
//	breakTag:
//		for _, edge := range sourceEdges {
//			edgeCondExpr := make([]map[string]interface{}, 0)
//			err = json.Unmarshal([]byte(edge["conditionExpression"].(string)), &edgeCondExpr)
//			if err != nil {
//				return
//			}
//			for _, condExpr := range edgeCondExpr {
//				// 条件判断
//				condExprStatus, err = h.ConditionalJudgment(condExpr)
//				if err != nil {
//					return
//				}
//				if condExprStatus {
//					// 进行节点跳转
//					h.TargetStateValue, err = h.DefinitionStructure.GetNode(edge["target"].(string))
//					if err != nil {
//						return
//					}
//
//					if h.TargetStateValue["clazz"] == "userTask" || h.TargetStateValue["clazz"] == "receiveTask" {
//						if h.TargetStateValue["assignValue"] == nil || h.TargetStateValue["assignType"] == "" {
//							err = errors.New("处理人不能为空")
//							return
//						}
//					}
//
//					h.UpdateValue["state"] = []map[string]interface{}{{
//						"id":             h.TargetStateValue["id"].(string),
//						"label":          h.TargetStateValue["label"],
//						"processor":      h.TargetStateValue["assignValue"],
//						"process_method": h.TargetStateValue["assignType"],
//					}}
//					err = h.commonProcessing(c)
//					if err != nil {
//						err = fmt.Errorf("流程流程跳转失败，%v", err.Error())
//						return
//					}
//
//					break breakTag
//				}
//			}
//		}
//		if !condExprStatus {
//			err = errors.New("所有流转均不符合条件，请确认。")
//			return
//		}
//	// 并行/聚合网关
//	case "parallelGateway":
//		// 入口，判断
//		sourceEdges, err = h.DefinitionStructure.GetEdge(h.TargetStateValue["id"].(string), "source")
//		if err != nil {
//			err = fmt.Errorf("查询流转信息失败，%v", err.Error())
//			return
//		}
//
//		targetEdges, err = h.DefinitionStructure.GetEdge(h.TargetStateValue["id"].(string), "target")
//		if err != nil {
//			err = fmt.Errorf("查询流转信息失败，%v", err.Error())
//			return
//		}
//
//		if len(sourceEdges) > 0 {
//			h.TargetStateValue, err = h.DefinitionStructure.GetNode(sourceEdges[0]["target"].(string))
//			if err != nil {
//				return
//			}
//		} else {
//			err = errors.New("并行网关流程不正确")
//			return
//		}
//
//		if len(sourceEdges) > 1 && len(targetEdges) == 1 {
//			// 入口
//			h.UpdateValue["state"] = make([]map[string]interface{}, 0)
//			for _, edge := range sourceEdges {
//				TargetStateValue, err := h.DefinitionStructure.GetNode(edge["target"].(string))
//				if err != nil {
//					return err
//				}
//				h.UpdateValue["state"] = append(h.UpdateValue["state"].([]map[string]interface{}), map[string]interface{}{
//					"id":             edge["target"].(string),
//					"label":          TargetStateValue["label"],
//					"processor":      TargetStateValue["assignValue"],
//					"process_method": TargetStateValue["assignType"],
//				})
//			}
//			err = h.circulation()
//			if err != nil {
//				err = fmt.Errorf("工单跳转失败，%v", err.Error())
//				return
//			}
//		} else if len(sourceEdges) == 1 && len(targetEdges) > 1 {
//			// 出口
//			parallelStatusOk, err = h.completeAllParallel(sourceEdges[0]["target"].(string))
//			if err != nil {
//				err = fmt.Errorf("并行检测失败，%v", err.Error())
//				return
//			}
//			if parallelStatusOk {
//				h.EndHistory = true
//				endAssignValue, ok := h.TargetStateValue["assignValue"]
//				if !ok {
//					endAssignValue = []int{}
//				}
//
//				endAssignType, ok := h.TargetStateValue["assignType"]
//				if !ok {
//					endAssignType = ""
//				}
//
//				h.UpdateValue["state"] = []map[string]interface{}{{
//					"id":             h.TargetStateValue["id"].(string),
//					"label":          h.TargetStateValue["label"],
//					"processor":      endAssignValue,
//					"process_method": endAssignType,
//				}}
//				err = h.circulation()
//				if err != nil {
//					err = fmt.Errorf("工单跳转失败，%v", err.Error())
//					return
//				}
//			} else {
//				h.EndHistory = false
//			}
//
//		} else {
//			err = errors.New("并行网关流程不正确")
//			return
//		}
//	// 包容网关
//	case "inclusiveGateway":
//		return
//	case "start":
//		StateValue["processor"] = []int{h.workOrderDetails.Creator}
//		StateValue["process_method"] = "person"
//		h.UpdateValue["state"] = []map[string]interface{}{StateValue}
//		err = h.circulation()
//		if err != nil {
//			return
//		}
//	case "userTask":
//		StateValue["processor"] = h.TargetStateValue["assignValue"].([]interface{})
//		StateValue["process_method"] = h.TargetStateValue["assignType"].(string)
//		h.UpdateValue["state"] = []map[string]interface{}{StateValue}
//		err = h.commonProcessing(c)
//		if err != nil {
//			return
//		}
//	case "receiveTask":
//		StateValue["processor"] = h.TargetStateValue["assignValue"].([]interface{})
//		StateValue["process_method"] = h.TargetStateValue["assignType"].(string)
//		h.UpdateValue["state"] = []map[string]interface{}{StateValue}
//		err = h.commonProcessing(c)
//		if err != nil {
//			return
//		}
//	case "scriptTask":
//		StateValue["processor"] = []int{}
//		StateValue["process_method"] = ""
//		h.UpdateValue["state"] = []map[string]interface{}{StateValue}
//	case "end":
//		StateValue["processor"] = []int{}
//		StateValue["process_method"] = ""
//		h.UpdateValue["state"] = []map[string]interface{}{StateValue}
//		err = h.commonProcessing(c)
//		if err != nil {
//			return
//		}
//	}
//
//	// 更新表单数据
//	for _, t := range tpls {
//		var (
//			tplValue []byte
//		)
//		tplValue, err = json.Marshal(t["tplValue"])
//		if err != nil {
//			h.tx.Rollback()
//			return
//		}
//
//		paramsValue.FormData = append(paramsValue.FormData, t["tplValue"])
//
//		// 是否可写，只有可写的模版可以更新数据
//		updateStatus := false
//		if h.StateValue["clazz"].(string) == "start" {
//			updateStatus = true
//		} else if writeTplList, writeOK := h.StateValue["writeTpls"]; writeOK {
//		tplListTag:
//			for _, writeTplId := range writeTplList.([]interface{}) {
//				if writeTplId == t["tplId"] { // 可写
//					// 是否隐藏，隐藏的模版无法修改数据
//					if hideTplList, hideOK := h.StateValue["hideTpls"]; hideOK {
//						if hideTplList != nil && len(hideTplList.([]interface{})) > 0 {
//							for _, hideTplId := range hideTplList.([]interface{}) {
//								if hideTplId == t["tplId"] { // 隐藏的
//									updateStatus = false
//									break tplListTag
//								} else {
//									updateStatus = true
//								}
//							}
//						} else {
//							updateStatus = true
//						}
//					} else {
//						updateStatus = true
//					}
//				}
//			}
//		} else {
//			// 不可写
//			updateStatus = false
//		}
//		if updateStatus {
//			err = h.tx.Model(&process.TplData{}).Where("id = ?", t["tplDataId"]).Update("form_data", tplValue).Error
//			if err != nil {
//				h.tx.Rollback()
//				return
//			}
//		}
//	}
//
//	// 流转历史写入
//	err = orm.Eloquent.Model(&cirHistoryValue).
//		Where("work_order = ?", WorkOrderId).
//		Find(&cirHistoryValue).
//		Order("create_time desc").Error
//	if err != nil {
//		h.tx.Rollback()
//		return
//	}
//	for _, t := range cirHistoryValue {
//		if t.Source != h.StateValue["id"] {
//			costDuration := time.Since(t.CreatedAt.Time)
//			costDurationValue = fmtDuration(costDuration)
//		}
//	}
//
//	// 获取当前用户信息
//	err = orm.Eloquent.Model(&currentUserInfo).
//		Where("user_id = ?", tools.GetUserId(c)).
//		Find(&currentUserInfo).Error
//	if err != nil {
//		return
//	}
//
//	cirHistoryData = process.CirculationHistory{
//		Model:        base.Model{},
//		Title:        h.workOrderDetails.Title,
//		WorkOrder:    h.workOrderDetails.Id,
//		State:        h.StateValue["label"].(string),
//		Source:       h.StateValue["id"].(string),
//		Target:       h.TargetStateValue["id"].(string),
//		Circulation:  CirculationValue,
//		Processor:    currentUserInfo.NickName,
//		ProcessorId:  tools.GetUserId(c),
//		CostDuration: costDurationValue,
//		Remarks:      remarks,
//	}
//	err = h.tx.Create(&cirHistoryData).Error
//	if err != nil {
//		h.tx.Rollback()
//		return
//	}
//
//	// 获取流程通知类型列表
//	err = json.Unmarshal(processInfo.Notice, &noticeList)
//	if err != nil {
//		return
//	}
//
//	bodyData := notify.BodyData{
//		SendTo: map[string]interface{}{
//			"userList": sendToUserList,
//		},
//		Subject:     sendSubject,
//		Description: sendDescription,
//		Classify:    noticeList,
//		ProcessId:   h.workOrderDetails.Process,
//		Id:          h.workOrderDetails.Id,
//		Title:       h.workOrderDetails.Title,
//		Creator:     currentUserInfo.NickName,
//		Priority:    h.workOrderDetails.Priority,
//		CreatedAt:   h.workOrderDetails.CreatedAt.Format("2006-01-02 15:04:05"),
//	}
//
//	// 判断目标是否是结束节点
//	if h.TargetStateValue["clazz"] == "end" && h.EndHistory == true {
//		sendSubject = "您的工单已处理完成"
//		sendDescription = "您的工单已处理完成，工单描述如下"
//		err = h.tx.Create(&process.CirculationHistory{
//			Model:       base.Model{},
//			Title:       h.workOrderDetails.Title,
//			WorkOrder:   h.workOrderDetails.Id,
//			State:       h.TargetStateValue["label"].(string),
//			Source:      h.TargetStateValue["id"].(string),
//			Processor:   currentUserInfo.NickName,
//			ProcessorId: tools.GetUserId(c),
//			Circulation: "结束",
//			Remarks:     "工单已结束",
//		}).Error
//		if err != nil {
//			h.tx.Rollback()
//			return
//		}
//		if len(noticeList) > 0 {
//			// 查询工单创建人信息
//			err = h.tx.Model(&system.SysUser{}).
//				Where("user_id = ?", h.workOrderDetails.Creator).
//				Find(&sendToUserList).Error
//			if err != nil {
//				return
//			}
//
//			bodyData.SendTo = map[string]interface{}{
//				"userList": sendToUserList,
//			}
//			bodyData.Subject = sendSubject
//			bodyData.Description = sendDescription
//
//			// 发送通知
//			go func(bodyData notify.BodyData) {
//				err = bodyData.SendNotify()
//				if err != nil {
//					return
//				}
//			}(bodyData)
//		}
//	}
//
//	h.tx.Commit() // 提交事务
//
//	// 发送通知
//	if len(noticeList) > 0 {
//		stateList := make([]interface{}, 0)
//		for _, v := range h.UpdateValue["state"].([]map[string]interface{}) {
//			stateList = append(stateList, v)
//		}
//		sendToUserList, err = GetPrincipalUserInfo(stateList, h.workOrderDetails.Creator)
//		if err != nil {
//			return
//		}
//
//		bodyData.SendTo = map[string]interface{}{
//			"userList": sendToUserList,
//		}
//		bodyData.Subject = sendSubject
//		bodyData.Description = sendDescription
//
//		// 发送通知
//		go func(bodyData notify.BodyData) {
//			err = bodyData.SendNotify()
//			if err != nil {
//				return
//			}
//		}(bodyData)
//	}
//
//	// 执行流程公共任务及节点任务
//	if h.StateValue["task"] != nil {
//		for _, task := range h.StateValue["task"].([]interface{}) {
//			tasks = append(tasks, task.(string))
//		}
//	}
//continueTag:
//	for _, task := range tasks {
//		for _, t := range execTasks {
//			if t == task {
//				continue continueTag
//			}
//		}
//		execTasks = append(execTasks, task)
//	}
//
//	paramsValue.Id = h.workOrderDetails.Id
//	paramsValue.Title = h.workOrderDetails.Title
//	paramsValue.Priority = h.workOrderDetails.Priority
//	params, err := json.Marshal(paramsValue)
//	if err != nil {
//		return err
//	}
//	go ExecTask(execTasks, string(params))
//
//	return
//}
