/**
 * @Author: lzw5399
 * @Date: 2021/3/19 0:09
 * @Desc:
 */
package util

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

func GenUUID() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}
