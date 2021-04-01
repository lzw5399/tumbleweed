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

	order := "desc"
	if r.Order == "asc" {
		order = r.Order
	}

	sort := r.Sort
	if sort == "" {
		sort = "update_time"
	}

	return db.Offset(r.Offset).Limit(r.Limit).Order(fmt.Sprintf("%s %s", sort, order))
}

func ApplyRawPaging(sql string, r *request.PagingRequest) string {
	order := "desc"
	if r.Order == "asc" {
		order = r.Order
	}
	sort := r.Sort
	if sort == "" {
		sort = "update_time"
	}

	sql += fmt.Sprintf(" order by %s %s ", sort, order)


	if r.Limit > 0 {
		sql += fmt.Sprintf(" limit %d ", r.Limit)
	}

	if r.Offset > 0 {
		sql += fmt.Sprintf(" offset %d ", r.Offset)
	}



	return sql
}
