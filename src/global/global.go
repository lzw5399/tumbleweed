/**
 * @Author: lzw5399
 * @Date: 2020/9/30 14:30
 * @Desc: global object, will initialized after project starting
 */
package global

import (
	"bank/distributedquery/src/config"

	golog "github.com/op/go-logging"
)

var (
	BANK_CONFIG config.Config
	BANK_LOGGER *golog.Logger
)
