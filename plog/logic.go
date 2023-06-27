package plog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/caticat/go_game_server/ptime"
)

const (
	c_calldepth  = 4
	c_logDirMod  = 0766
	c_logFileMod = 0666
)

var (
	g_mapLog = map[ELogLevel]*struct {
		Prefix string
		Log    func(v ...any)
		LogF   func(format string, v ...any)
		LogLn  func(v ...any)
	}{
		ELogLevel_None:  {"[None ]", logPrint, logPrintF, logPrintLn},
		ELogLevel_Debug: {"[Debug]", logPrint, logPrintF, logPrintLn},
		ELogLevel_Info:  {"[Info ]", logPrint, logPrintF, logPrintLn},
		ELogLevel_Warn:  {"[Warn ]", logPrint, logPrintF, logPrintLn},
		ELogLevel_Error: {"[Error]", logPrint, logPrintF, logPrintLn},
		ELogLevel_Fatal: {"[Fatal]", logFatal, logFatalF, logFatalLn},
		ELogLevel_Panic: {"[Panic]", logPanic, logPanicF, logPanicLn},
	}

	g_logLevel      = ELogLevel_None
	g_logFilePrefix = ""

	// 前缀逻辑
	g_logLevelPre = ELogLevel_None

	// 日志文件名逻辑
	g_logFileScrollTime int64    = -1
	g_logFile           *os.File = nil
)

func getLogLevel() ELogLevel       { return g_logLevel }
func setLogLevel(l ELogLevel)      { g_logLevel = l }
func getLogFilePrefix() string     { return g_logFilePrefix }
func setLogFilePrefix(f string)    { g_logFilePrefix = f }
func getLogFileScrollTime() int64  { return g_logFileScrollTime }
func setLogFileScrollTime(t int64) { g_logFileScrollTime = t }
func getLogFile() *os.File         { return g_logFile }
func setLogFile(f *os.File)        { g_logFile = f }

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
