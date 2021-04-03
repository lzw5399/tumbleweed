/**
 * @Author: lzw5399
 * @Date: 2021/3/31 13:28
 * @Desc: 流转历史
 */
package controller

import (
	"github.com/labstack/echo/v4"

	"workflow/src/global/response"
	"workflow/src/model/request"
	"workflow/src/service"
)

// @Tags process-instances
// @Summary 获取流转历史列表
// @Produce json
// @param id path int true "实例id"
// @param request query request.HistoryListRequest true "request"
// @param WF-TENANT-CODE header string true "WF-TENANT-CODE"
// @param WF-CURRENT-USER header string true "WF-CURRENT-USER"
// @Success 200 {object} response.HttpResponse
// @Router /api/wf/process-instances/{id}/history [GET]
func ListHistory(c echo.Context) error {
	// 从queryString获取分页参数
	var r request.HistoryListRequest
	if err := c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	trainNodes, err := service.ListHistory(&r, c)
	if err != nil {
		return response.Failed(c, err)
	}

	return response.OkWithData(c, trainNodes)
}
