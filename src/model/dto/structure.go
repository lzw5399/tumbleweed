/**
 * @Author: lzw5399
 * @Date: 2021/3/26 23:46
 * @Desc:
 */
package dto

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Structure struct {
	Edges  []Edge        `json:"edges"`
	Nodes  []Node        `json:"nodes"`
	Groups []interface{} `json:"groups"`
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 StateArray
func (j *Structure) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal dto.Structure value:", value))
	}

	var result Structure
	err := json.Unmarshal(bytes, &result)
	*j = result

	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j Structure) Value() (driver.Value, error) {
	v, err := json.Marshal(j)
	return string(v), err
}
