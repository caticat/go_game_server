package main

import "github.com/caticat/go_game_server/plog"

const (
	// Title
	WINDOW_TITLE = "PEtcdGUI"

	// key
	PETCD_CFG_BASE            = "cfg-base"
	PETCD_CFG_ENDPOINT_SELECT = "cfg-endpoint-select"
	PETCD_CFG_ENDPOINT_LIST   = "cfg-endpoint-list"
	PETCD_CFG_LOG_LEVEL       = "cfg-log-level"

	// value
	PETCD_CFG_ENDPOINT_SELECT_DEFAULT = "default"
	PETCD_CFG_BASE_DEFAULT            = `
	{
		"dial-timeout": 1,
		"operation-timeout": 1,
		"lease-timeout-before-keep-alive": 10
	}
	`
	PETCD_CFG_ENDPOINT_LIST_DEFAULT = `
	{
		"default":[
			"localhost:60001",
			"localhost:60002",
			"localhost:60003"
		]
	}
	`
	PETCD_CFG_LOG_LEVEL_DEFAULT = plog.ELogLevel_Debug

	// plog
	PLOG_MAX_SIZE = 1 << 20 // 日志数据量
	PLOG_CHAN_LEN = 100     // 数据管道长度
)
