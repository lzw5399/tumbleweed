/**
 * @Author: lzw5399
 * @Date: 2021/3/26 23:58
 * @Desc:
 */
package dto

type Node struct {
	X             float64  `json:"x"`
	Y             float64  `json:"y"`
	Id            string   `json:"id"`
	Size          []int    `json:"size"`
	Sort          string   `json:"sort"`
	Clazz         string   `json:"clazz"`
	Label         string   `json:"label"`
	Shape         string   `json:"shape"`
	IsHideNode    bool     `json:"isHideNode,omitempty"`
	AssignType    string   `json:"assignType,omitempty"`
	ActiveOrder   bool     `json:"activeOrder,omitempty"`
	AssignValue   []string `json:"assignValue,omitempty"`
	IsCounterSign bool     `json:"isCounterSign,omitempty"`
}
