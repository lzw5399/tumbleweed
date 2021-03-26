/**
 * @Author: lzw5399
 * @Date: 2021/3/26 23:46
 * @Desc:
 */
package dto

type Edge struct {
	Id                  string `json:"id"`
	Sort                string `json:"sort"`
	Clazz               string `json:"clazz"`
	Label               string `json:"label"`
	Shape               string `json:"shape"`
	Source              string `json:"source"`
	Target              string `json:"target"`
	SourceAnchor        int64  `json:"sourceAnchor"`
	TargetAnchor        int64  `json:"targetAnchor"`
	FlowProperties      string `json:"flowProperties"`
	ConditionExpression string `json:"conditionExpression,omitempty"` // 表达式
}
