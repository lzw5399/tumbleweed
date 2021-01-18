/**
 * @Author: lzw5399
 * @Date: 2021/1/16 23:30
 * @Desc:
 */
package util

import (
	"encoding/json"
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

func StructToBytes(m interface{}) []byte {
	bytes, _ := json.Marshal(m)

	return bytes
}
