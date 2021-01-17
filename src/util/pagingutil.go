/**
 * @Author: lzw5399
 * @Date: 2021/1/17 22:59
 * @Desc:
 */
package util

import (
	"fmt"
	"reflect"
)

type PagingOption struct {
	offset     int
	limit      int
	originList interface{}
}

func NewPaging(list interface{}) *PagingOption {
	return &PagingOption{
		limit:      -1,
		originList: &list,
	}
}

func (option *PagingOption) Offset(offset int) *PagingOption {
	option.offset = offset
	return option
}

func (option *PagingOption) Limit(limit int) *PagingOption {
	option.limit = limit
	return option
}

func (option *PagingOption) Get(finalList interface{}) {
	//finalList = []interface{}{}
	//if option.offset > len(option.originList) {
	//	return
	//}

	v := reflect.ValueOf(option.originList)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t:= v.Type()
	fmt.Println(t)
	//for i, v := range option.originList {
	//	if i > option.offset-1 {
	//		finalList = append(finalList.([]interface{}), v)
	//	}
	//}
}
