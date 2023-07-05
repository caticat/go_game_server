package plog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/caticat/go_game_server/ptime"
)

func doLog(l ELogLevel, v ...any) {
	if l < getLogLevel() {
		return
	}
	c := g_mapLog[l]
	trySetPrefix(l, c.Prefix)
	checkLogFile()

	c.Log(v...)
}

func doLogF(l ELogLevel, format string, v ...any) {
	if l < getLogLevel() {
		return
	}
	c := g_mapLog[l]
	trySetPrefix(l, c.Prefix)
	checkLogFile()

	c.LogF(format, v...)
}

func doLogLn(l ELogLevel, v ...any) {
	if l < getLogLevel() {
		return
	}
	c := g_mapLog[l]
	trySetPrefix(l, c.Prefix)
	checkLogFile()

	c.LogLn(v...)
}

func trySetPrefix(l ELogLevel, prefix string) {
	if l == g_logLevelPre {
		return
	}
	g_logLevelPre = l

	log.SetPrefix(prefix)
}

// 检测设置日志文件
func checkLogFile() {
	if len(getLogFilePrefix()) == 0 {
		return
	}
	n := time.Now().Local()
	if n.Unix() < getLogFileScrollTime() {
		return
	}
	setLogFileScrollTime(ptime.GetNextHourTime(n))

	f := getLogFile()
	if f != nil {
		f.Close()
	}

	// 创建/打开日志文件
	fileName := fmt.Sprintf("%s.%s", getLogFilePrefix(), n.Format("2006-01-02-15"))
	_, err := os.Stat(fileName)
	if err != nil {
		// 路径创建
		d := filepath.Dir(fileName)
		_, err = os.Stat(d)
		if err != nil {
			os.MkdirAll(d, c_logDirMod)
		}
		f, err = os.Create(fileName)
	} else {
		f, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, c_logFileMod)
	}
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Println("os.OpenFile failed,error:", err)
		return
	}
	setLogFile(f)
	log.SetOutput(f)
}

func logPrint(v ...any)                 { log.Output(c_calldepth, fmt.Sprint(v...)) }
func logPrintF(format string, v ...any) { log.Output(c_calldepth, fmt.Sprintf(format, v...)) }
func logPrintLn(v ...any)               { log.Output(c_calldepth, fmt.Sprintln(v...)) }

func logFatal(v ...any) { log.Output(c_calldepth, fmt.Sprint(v...)); os.Exit(1) }
func logFatalF(format string, v ...any) {
	log.Output(c_calldepth, fmt.Sprintf(format, v...))
	os.Exit(1)
}
func logFatalLn(v ...any) { log.Output(c_calldepth, fmt.Sprintln(v...)); os.Exit(1) }

func logPanic(v ...any) { s := fmt.Sprint(v...); log.Output(c_calldepth, s); panic(s) }
func logPanicF(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	log.Output(c_calldepth, s)
	panic(s)
}
func logPanicLn(v ...any) { s := fmt.Sprintln(v...); log.Output(c_calldepth, s); panic(s) }
