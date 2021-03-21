/**
 * @Author: lzw5399
 * @Date: 2021/1/16 23:30
 * @Desc:
 */
package util

import (
	"encoding/json"

	"gorm.io/datatypes"

	"workflow/src/model"
)

func StringToMap(jsonStr string) (map[string]string, error) {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func MapToString(m map[string]interface{}) string {
	bytes, _ := json.Marshal(m)

	return string(bytes)
}

func MapToBytes(m map[string]interface{}) []byte {
	bytes, _ := json.Marshal(m)

	return bytes
}

func MarshalToBytes(m interface{}) []byte {
	bytes, _ := json.Marshal(m)

	return bytes
}

func MarshalToDbJson(m interface{}) datatypes.JSON {
	return datatypes.JSON(MarshalToBytes(m))
}

func MarshalToString(m interface{}) string {
	return string(MarshalToBytes(m))
}

func UnmarshalToInstanceVariables(m datatypes.JSON) []model.InstanceVariable {
	var variables []model.InstanceVariable
	_ = json.Unmarshal([]byte(m), &variables)

	return variables
}
