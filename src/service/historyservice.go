/**
 * @Author: lzw5399
 * @Date: 2021/3/31 13:44
 * @Desc: 流转历史
 */
package service

import (
	"github.com/labstack/echo/v4"

	"workflow/src/global"
	"workflow/src/global/constant"
	"workflow/src/global/shared"
	"workflow/src/model"
	"workflow/src/model/request"
	"workflow/src/model/response"
	"workflow/src/util"
)

func ListHistory(r *request.HistoryListRequest, c echo.Context) (*response.PagingResponse, error) {
	var (
		histories   []model.CirculationHistory
		tenantId, _ = util.GetWorkContext(c)
	)

	db := global.BankDb.
		Model(&model.ProcessInstance{}).
		Where("tenant_id = ?", tenantId).
		Where("process_instance.id = ?", r.ProcessInstanceId)

	// 根据type的不同有不同的逻辑
	switch r.Type {
	case constant.HistoryTypeFull:
	case constant.HistoryTypeSimple:
	default:
		return nil, util.BadRequest.New("type不合法")
	}

	if r.Keyword != "" {
		db = db.Where("title ~ ?", r.Keyword)
	}

	var count int64
	db.Count(&count)

	db = shared.ApplyPaging(db, &r.PagingRequest)
	err := db.Find(&histories).Error

	return &response.PagingResponse{
		TotalCount:   count,
		CurrentCount: int64(len(histories)),
		Data:         &histories,
	}, err
}
