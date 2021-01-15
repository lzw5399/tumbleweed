/**
 * @Author: lzw5399
 * @Date: 2021/1/15 23:35
 * @Desc:
 */
package service

import (
	"errors"
	"workflow/src/global"
	. "workflow/src/model"
	"workflow/src/model/request"
)

// 创建新的process流程
func CreateProcess(r *request.Definitions) error {
	// 检查流程是否已存在
	var c int64
	global.BankDb.Model(&ProcessDefinition{}).Where("id=?", r.Process.ID).Count(&c)
	if c != 0 {
		return errors.New("当前流程标识已经在，请检查后重试")
	}

	// 创建表数据
	process := r.ToProcess()
	process.Create()

	return nil
}
