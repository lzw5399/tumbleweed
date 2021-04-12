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
		Where("process_instance.id = ?", r.Id).
		Joins("inner join wf.circulation_history on circulation_history.process_instance_id = process_instance.id")

	// 根据type的不同有不同的逻辑
	switch r.Type {
	case constant.HistoryTypeFull:
	case constant.HistoryTypeSimple:
		db.Where("source_id not like 'exclusiveGateway%' and source_id not like 'parallelGateway%'")
	default:
		return nil, util.BadRequest.New("type不合法")
	}

	if r.Keyword != "" {
		db = db.Where("title ~ ?", r.Keyword)
	}

	var count int64
	db.Count(&count)

	db = shared.ApplyPaging(db, &r.PagingRequest)
	err := db.Select("circulation_history.*").Scan(&histories).Error

	return &response.PagingResponse{
		TotalCount:   count,
		CurrentCount: int64(len(histories)),
		Data:         &histories,
	}, err
}
