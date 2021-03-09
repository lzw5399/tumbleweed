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
	"github.com/labstack/gommon/log"

	"workflow/src/global"
	"workflow/src/global/shared"
	"workflow/src/model"
	"workflow/src/model/request"
	"workflow/src/model/response"

	"gorm.io/gorm"
)

type InstanceService interface {
	CreateProcessInstance(*request.ProcessInstanceRequest, uint) (uint, error)
	Get(uint) (*model.ProcessInstance, error)
	List(*request.PagingRequest) (*response.PagingResponse, error)
	GetVariable(*request.GetVariableRequest) (*response.InstanceVariableResponse, error)
	ListVariables(r *request.GetVariableListRequest) (*response.PagingResponse, error)
}

type instanceService struct {
}

func NewInstanceService() *instanceService {
	return &instanceService{}
}

// 创建实例
func (i *instanceService) CreateProcessInstance(r *request.ProcessInstanceRequest, currentUserId uint) (uint, error) {
	var (
		variableValue []interface{}
	)

	// 获取变量值
	err := json.Unmarshal(r.State, &variableValue)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	PrefillVariables(variableValue, currentUserId)
	err = GetVariableValue(variableValue, currentUserId)
	if err != nil {
		err = fmt.Errorf("获取处理人变量值失败，%v", err.Error())
		return
	}
}

// 获取单个ProcessInstance
func (i *instanceService) Get(instanceId uint) (*model.ProcessInstance, error) {
	var instance model.ProcessInstance
	err := global.BankDb.Where("id=?", instanceId).First(&instance).Error

	return &instance, err
}

// 获取ProcessInstance列表
func (i *instanceService) List(r *request.PagingRequest) (*response.PagingResponse, error) {
	var instances []model.ProcessInstance
	db := shared.ApplyPaging(global.BankDb, r)
	err := db.Find(&instances).Error

	var totalCount int64
	global.BankDb.Model(&model.ProcessInstance{}).Count(&totalCount)

	return &response.PagingResponse{
		TotalCount:   totalCount,
		CurrentCount: int64(len(instances)),
		Data:         &instances,
	}, err
}

// 获取实例的某一个变量
func (i *instanceService) GetVariable(r *request.GetVariableRequest) (*response.InstanceVariableResponse, error) {
	var str string
	err := global.BankDb.Model(&model.ProcessInstance{}).Where("id=?", r.InstanceId).Select("variables").First(&str).Error
	if err != nil {
		return nil, errors.New("当前instance实例不存在")
	}

	var variables []model.InstanceVariable
	err = json.Unmarshal([]byte(str), &variables)
	if err != nil {
		return nil, errors.New("获取失败")
	}

	for _, v := range variables {
		if v.Name == r.VariableName {
			return &response.InstanceVariableResponse{
				Name:  v.Name,
				Type:  v.Type,
				Value: v.Value,
			}, nil
		}
	}

	return nil, errors.New("变量不存在")
}

func (i *instanceService) ListVariables(r *request.GetVariableListRequest) (*response.PagingResponse, error) {
	var str string
	err := global.BankDb.Model(&model.ProcessInstance{}).Where("id=?", r.InstanceId).Select("variables").First(&str).Error
	if err != nil {
		return nil, errors.New("当前instance实例不存在")
	}

	var variables []model.InstanceVariable
	err = json.Unmarshal([]byte(str), &variables)
	if err != nil {
		return nil, errors.New("获取失败")
	}

	//var finalVariables []model.InstanceVariable
	//util.NewPaging(variables).Offset(r.Offset).Limit(r.Limit).Get(&finalVariables)

	return &response.PagingResponse{
		TotalCount:   int64(len(variables)),
		CurrentCount: int64(len(variables)),
		Data:         &variables,
	}, err
}

// 创建流程实例记录
func createProcessInstance(r *request.ProcessInstanceRequest) (*model.ProcessInstance, error) {
	// 检查流程是否存在
	var process model.Process
	err := global.BankDb.Where("code=?", r.ProcessCode).First(&process).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("当前流程不存在，请检查后重试")
	}

	// 创建流程
	instance := r.ProcessInstance(process.Id)
	err = global.BankDb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&instance).Error; err != nil {
			return err
		}

		// 返回nil提交事务
		return nil
	})

	return &instance, err
}

// 创建流程实例的时候预先处理
func PrefillVariables(stateList []interface{}, creator int) (err error) {
	var (
		userInfo system.SysUser
		deptInfo system.Dept
	)

	// 变量转为实际的数据
	for _, stateItem := range stateList {
		if stateItem.(map[string]interface{})["process_method"] == "variable" {
			for processorIndex, processor := range stateItem.(map[string]interface{})["processor"].([]interface{}) {
				if int(processor.(float64)) == 1 {
					// 创建者
					stateItem.(map[string]interface{})["processor"].([]interface{})[processorIndex] = creator
				} else if int(processor.(float64)) == 2 {
					// 1. 查询用户信息
					err = orm.Eloquent.Model(&userInfo).Where("user_id = ?", creator).Find(&userInfo).Error
					if err != nil {
						return
					}
					// 2. 查询部门信息
					err = orm.Eloquent.Model(&deptInfo).Where("dept_id = ?", userInfo.DeptId).Find(&deptInfo).Error
					if err != nil {
						return
					}

					// 3. 替换处理人信息
					stateItem.(map[string]interface{})["processor"].([]interface{})[processorIndex] = deptInfo.Leader
				}
			}
			stateItem.(map[string]interface{})["process_method"] = "person"
		}
	}

	return
}