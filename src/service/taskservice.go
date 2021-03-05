/**
 * @Author: lzw5399
 * @Date: 2021/2/9 18:47
 * @Desc:
 */
package service

import (
	"workflow/src/global"
	"workflow/src/global/shared"
	"workflow/src/model"
	"workflow/src/model/request"
	"workflow/src/model/response"
)

type TaskService interface {
	List(*request.PagingRequest) (*response.PagingResponse, error)
}
type taskService struct {
}

func NewTaskService() *taskService {
	return &taskService{}
}

func (t *taskService) List(r *request.PagingRequest) (*response.PagingResponse, error) {
	var tasks []model.UserTask
	db := shared.ApplyPaging(global.BankDb, r)
	err := db.Find(&tasks).Error

	var totalCount int64
	global.BankDb.Model(&model.UserTask{}).Count(&totalCount)

	return &response.PagingResponse{
		TotalCount:   totalCount,
		CurrentCount: int64(len(tasks)),
		Data:         &tasks,
	}, err
}
