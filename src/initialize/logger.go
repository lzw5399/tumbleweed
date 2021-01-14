/**
 * @Author: lzw5399
 * @Date: 2020/9/30 15:25
 * @Desc: auto load logger after app start
 */
package initialize

import (
	"fmt"
	"io"
	"os"
	"strings"

	"workflow/src/config"
	"workflow/src/global"

	golog "github.com/op/go-logging"
)

const (
	MODULE = "bank-workflow-engine"
)

var (
	defaultFormatter = `%{time:2006/01/02 - 15:04:05.000} %{longfile} %{color:bold}▶ [%{level:.6s}] %{message}%{color:reset}`
)

func init() {
	c := global.BankConfig.Log
	if c.Prefix == "" {
		_ = fmt.Errorf("logger prefix not found")
	}
	logger := golog.MustGetLogger(MODULE)
	var backends []golog.Backend
	registerStdout(c, &backends)

	golog.SetBackend(backends...)
	global.BankLogger = logger
}

func registerStdout(c config.Log, backends *[]golog.Backend) {
	if c.Stdout != "" {
		level, err := golog.LogLevel(c.Stdout)
		if err != nil {
			fmt.Println(err)
		}
		*backends = append(*backends, createBackend(os.Stdout, c, level))
	}
}

func createBackend(w io.Writer, c config.Log, level golog.Level) golog.Backend {
	backend := golog.NewLogBackend(w, c.Prefix, 0)
	stdoutWriter := false
	if w == os.Stdout {
		stdoutWriter = true
	}
	format := getLogFormatter(c, stdoutWriter)
	backendLeveled := golog.AddModuleLevel(golog.NewBackendFormatter(backend, format))
	backendLeveled.SetLevel(level, MODULE)
	return backendLeveled
}

func getLogFormatter(c config.Log, stdoutWriter bool) golog.Formatter {
	pattern := defaultFormatter
	if !stdoutWriter {
		// Color is only required for console output
		// Other writers don't need %{color} tag
		pattern = strings.Replace(pattern, "%{color:bold}", "", -1)
		pattern = strings.Replace(pattern, "%{color:reset}", "", -1)
	}

	return golog.MustStringFormatter(pattern)
}
