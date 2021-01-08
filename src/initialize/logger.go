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
	"time"

	"workflow/src/config"
	"workflow/src/global"
	"workflow/src/util"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	golog "github.com/op/go-logging"
)

const (
	LOG_DIR      = "log"
	LOG_SOFTLINK = "latest_log"
	MODULE       = "bank/distributedquery"
)

var (
	defaultFormatter = `%{time:2006/01/02 - 15:04:05.000} %{longfile} %{color:bold}▶ [%{level:.6s}] %{message}%{color:reset}`
)

func init() {
	c := global.BANK_CONFIG.Log
	if c.Prefix == "" {
		_ = fmt.Errorf("logger prefix not found")
	}
	logger := golog.MustGetLogger(MODULE)
	var backends []golog.Backend
	registerStdout(c, &backends)
	if c.LogFile {
		if fileWriter := registerFile(c, &backends); fileWriter != nil {
			gin.DefaultWriter = io.MultiWriter(fileWriter, os.Stdout)
		}
	}

	golog.SetBackend(backends...)
	global.BANK_LOGGER = logger
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

func registerFile(c config.Log, backends *[]golog.Backend) io.Writer {
	if c.File != "" {
		if !util.PathExists(LOG_DIR) {
			// directory not exist
			fmt.Println("create log directory")
			_ = os.Mkdir(LOG_DIR, os.ModePerm)
		}
		fileWriter, err := rotatelogs.New(
			LOG_DIR+string(os.PathSeparator)+"%Y-%m-%d-%H-%M.log",
			// generate soft link, point to latest log file
			rotatelogs.WithLinkName(LOG_SOFTLINK),
			// maximum time to save log files
			rotatelogs.WithMaxAge(7*24*time.Hour),
			// time period of log file switching
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			fmt.Println(err)
		}
		level, err := golog.LogLevel(c.File)
		if err != nil {
			fmt.Println(err)
		}
		*backends = append(*backends, createBackend(fileWriter, c, level))

		return fileWriter
	}
	return nil
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
	if !c.LogFile {
		// Remove %{logfile} tag
		pattern = strings.Replace(pattern, "%{longfile}", "", -1)
	}

	return golog.MustStringFormatter(pattern)
}
