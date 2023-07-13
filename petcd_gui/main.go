package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/caticat/go_game_server/petcd"
	"github.com/caticat/go_game_server/plog"
)

func main() {
	plog.Init(plog.ELogLevel_Debug, "")

	g_app = app.NewWithID("com.github.caticat.go_game_server.petcd_gui")
	defer g_app.Quit()

	conf := getConf()
	if err := conf.Init(); err != nil {
		plog.FatalLn(err)
	}
	plog.SetOutput(NewPLogWriter()) // 配置读取后才将日志输出只向GUI界面

	err := petcd.Init(conf.GetCfgETCD())
	initData()

	runGUI(err)

	close()
}

func close() {
	petcd.Close()
}
