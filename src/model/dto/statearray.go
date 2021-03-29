/**
 * @Author: lzw5399
 * @Date: 2021/3/25 22:14
 * @Desc:
 */
package dto

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type StateArray []State

type State struct {
	Id                 string `json:"id"`
	Label              string `json:"label"`
	Processor          []int  `json:"processor"`          // 完整的处理人列表
	CompletedProcessor []int  `json:"completedProcessor"` // 已处理的人
	ProcessMethod      string `json:"processMethod"`      // 处理方式(角色 用户等)
	AssignValue        []int  `json:"assignValue"`        // 指定的处理者(用户的id或者角色的id)
	AvailableEdges     []Edge `json:"availableEdges"`     // 可走的线路
	IsCounterSign      bool   `json:"isCounterSign"`      // 是否是会签
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 StateArray
func (j *StateArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal dto.StateArray value:", value))
	}

	var result StateArray
	err := json.Unmarshal(bytes, &result)
	*j = result

	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j StateArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}

	v, err := json.Marshal(j)
	return string(v), err
}
