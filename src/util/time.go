/**
 * @Author: lzw5399
 * @Date: 2021/3/19 19:05
 * @Desc:
 */
package util

import (
	"fmt"
	"time"
)

// 时间格式化
func FmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute

	return fmt.Sprintf("%02d小时 %02d分钟", h, m)
}
