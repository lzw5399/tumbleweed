/**
 * @Author: lzw5399
 * @Date: 2021/1/17 21:34
 * @Desc:
 */
package response

type InstanceVariableResponse struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
	// Scope string      `json:"scope"`
}
