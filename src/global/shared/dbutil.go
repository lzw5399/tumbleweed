/**
 * @Author: lzw5399
 * @Date: 2021/1/17 18:43
 * @Desc:
 */
package shared

import (
	"fmt"
	
	"workflow/src/model/request"

	"gorm.io/gorm"
)

func ApplyPaging(db *gorm.DB, r *request.PagingRequest) *gorm.DB {
	// 如果等于0说明没传Limit参数，那么等于-1(不限制)
	limit := r.Limit
	if limit == 0 {
		limit = -1
	}

	order := "asc"
	if r.Order == "desc" {
		order = r.Order
	}

	sort := r.Sort
	if sort == "" {
		sort = "id"
	}

	return db.Offset(r.Offset).Limit(r.Limit).Order(fmt.Sprintf("%s %s", sort, order))
}
