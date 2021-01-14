/**
 * @Author: lzw5399
 * @Date: 2021/1/14 23:13
 * @Desc: 用户组同任务的关系
 */
package model

// 用户组同任务的关系
type IdentityLink struct {
	DbBase
	Group      string `json:"group,omitempty"`
	Type       string `json:"type,omitempty"`
	UserID     string `json:"userid,omitempty"`
	UserName   string `json:"username,omitempty"`
	TaskID     int    `json:"taskID,omitempty"`
	Step       int    `json:"step"`
	ProcInstID int    `json:"procInstID,omitempty"`
	Company    string `json:"company,omitempty"`
	Comment    string `json:"comment,omitempty"`
}
