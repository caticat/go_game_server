package main

import "github.com/caticat/go_game_server/plog"

const (
	// App
	APP_ID       = "com.github.caticat.go_game_server.petcd_gui"
	WINDOW_TITLE = "PEtcdGUI"

	// key
	PETCD_CFG_BASE        = "cfg-base"
	PETCD_CFG_CONN_SELECT = "cfg-conn-select"
	PETCD_CFG_CONN_LIST   = "cfg-conn-list"
	PETCD_CFG_LOG_LEVEL   = "cfg-log-level"

	// value
	PETCD_CFG_BASE_DEFAULT = `{
	"dial-timeout": 1,
	"operation-timeout": 1,
	"lease-timeout-before-keep-alive": 10
}
`
	PETCD_CFG_CONN_SELECT_DEFAULT = "default"
	PETCD_CFG_CONN_LIST_DEFAULT   = `{
	"default":
	{
		"endpoints":
		[
			"localhost:60001",
			"localhost:60002",
			"localhost:60003"
		],
		"auth":
		{
			"username":"",
			"password":""
		}
	}
}
`
	PETCD_CFG_LOG_LEVEL_DEFAULT = plog.ELogLevel_Debug

	// plog
	PLOG_MAX_SIZE = 1 << 20 // 日志数据量
	PLOG_CHAN_LEN = 100     // 数据管道长度

	// gui
	STR_EMPTY                           = ""
	STR_NIL                             = "<nil>"
	GUI_HOME_MAIN_OFFSET                = 0.4             // Home界面主窗口分隔偏移
	GUI_HOME_SEARCH_PLACEHOLDER         = "Search Key..." // Home界面Search提示
	GUI_HOME_EDIT_ENTRY_LINE_NUM        = 20              // Home界面输入框的行数
	GUI_SETTING_CREATE_CONN_PLACEHOLDER = `{
	"endpoints":
	[
	],
	"auth":
	{
		"username":"",
		"password":""
	}
}
` // Setting界面创建连接提示
	GUI_SETTING_EDIT_CONN_ENTRY_LINE_NUM = 20   // 编辑连接界面输入框的行数
	GUI_WINDOW_INIT_SIZE_W               = 1000 // 窗口初始大小W
	GUI_WINDOW_INIT_SIZE_H               = 700  // 窗口初始大小H
)
