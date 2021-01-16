/**
 * @Author: lzw5399
 * @Date: 2021/1/16 13:52
 * @Desc:
 */
package util

import "fmt"

func PropertyNotFound(name string) string {
	return fmt.Sprintf("当前bpmn xml不合法, 缺少必要的属性: %s", name)
}
