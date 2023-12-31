package plog

type ELogLevel int

const (
	ELogLevel_None ELogLevel = iota
	ELogLevel_Debug
	ELogLevel_Info
	ELogLevel_Warn
	ELogLevel_Error
	ELogLevel_Fatal
	ELogLevel_Panic
)

const (
	c_calldepth  = 4
	c_logDirMod  = 0766
	c_logFileMod = 0666
)

type SLogLevel string

const (
	SLogLevel_Debug SLogLevel = "debug"
	SLogLevel_Info  SLogLevel = "info"
	SLogLevel_Warn  SLogLevel = "warn"
)
