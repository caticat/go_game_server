package plog

import (
	"os"
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
func SetLogLevel(l ELogLevel)      { g_logLevel = l }
func getLogFilePrefix() string     { return g_logFilePrefix }
func setLogFilePrefix(f string)    { g_logFilePrefix = f }
func getLogFileScrollTime() int64  { return g_logFileScrollTime }
func setLogFileScrollTime(t int64) { g_logFileScrollTime = t }
func getLogFile() *os.File         { return g_logFile }
func setLogFile(f *os.File)        { g_logFile = f }
