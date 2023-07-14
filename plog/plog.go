package plog

// 日志模块
// 支持
//	控制台日志
//	文件日志(按小时拆分)

import (
	"io"
	"log"
)

func Init(l ELogLevel, logFile string) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	SetLogLevel(l)
	setLogFilePrefix(logFile)
}

func SetOutput(o io.Writer) {
	log.SetOutput(o)
}

func SetShortFile() {
	flags := log.Flags()
	flags = flags&(0^log.Llongfile) | log.Lshortfile
	log.SetFlags(flags)
}

func ToLogLevel(s string) ELogLevel {
	switch SLogLevel(s) {
	case SLogLevel_Debug:
		return ELogLevel_Debug
	case SLogLevel_Info:
		return ELogLevel_Info
	case SLogLevel_Warn:
		return ELogLevel_Warn
	default:
		return ELogLevel_Info
	}
}

func ToLogLevelName(logLevel ELogLevel) string {
	switch logLevel {
	case ELogLevel_Debug:
		return string(SLogLevel_Debug)
	case ELogLevel_Info:
		return string(SLogLevel_Info)
	case ELogLevel_Warn:
		return string(SLogLevel_Warn)
	default:
		return string(SLogLevel_Info)
	}
}

func Debug(v ...any)                 { doLog(ELogLevel_Debug, v...) }
func DebugF(format string, v ...any) { doLogF(ELogLevel_Debug, format, v...) }
func DebugLn(v ...any)               { doLogLn(ELogLevel_Debug, v...) }
func Info(v ...any)                  { doLog(ELogLevel_Info, v...) }
func InfoF(format string, v ...any)  { doLogF(ELogLevel_Info, format, v...) }
func InfoLn(v ...any)                { doLogLn(ELogLevel_Info, v...) }
func Warn(v ...any)                  { doLog(ELogLevel_Warn, v...) }
func WarnF(format string, v ...any)  { doLogF(ELogLevel_Warn, format, v...) }
func WarnLn(v ...any)                { doLogLn(ELogLevel_Warn, v...) }
func Error(v ...any)                 { doLog(ELogLevel_Error, v...) }
func ErrorF(format string, v ...any) { doLogF(ELogLevel_Error, format, v...) }
func ErrorLn(v ...any)               { doLogLn(ELogLevel_Error, v...) }
func Fatal(v ...any)                 { doLog(ELogLevel_Fatal, v...) }
func FatalF(format string, v ...any) { doLogF(ELogLevel_Fatal, format, v...) }
func FatalLn(v ...any)               { doLogLn(ELogLevel_Fatal, v...) }
func Panic(v ...any)                 { doLog(ELogLevel_Panic, v...) }
func PanicF(format string, v ...any) { doLogF(ELogLevel_Panic, format, v...) }
func PanicLn(v ...any)               { doLogLn(ELogLevel_Panic, v...) }
