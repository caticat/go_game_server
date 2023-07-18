package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/caticat/go_game_server/petcd"
	"github.com/caticat/go_game_server/plog"
)

func main() {
	// 初始化日志
	plog.Init(plog.ELogLevel_Debug, "")

	// 初始化App
	setApp(app.NewWithID(APP_ID))

	// 初始化配置
	conf := getConf()
	if err := conf.Init(); err != nil {
		plog.FatalLn(err)
	}

	// 日志额外设置
	plog.SetOutput(NewPLogWriter()) // 配置读取后才将日志输出只向GUI界面
	plog.SetShortFile()             // GUI中采用短文件名记录日志

	// 初始化ETCD连接
	err := petcd.Init(conf.GetCfgETCD())
	if err == nil {
		initData() // 同步ETCD的数据到内存数据结构
	}

	// 界面运行
	runGUI(err)

	// 关闭
	close(true)
}

func close(isExit bool) {
	petcd.Close()
	if isExit {
		a := getApp()
		if !a.Driver().Device().IsMobile() {
			a.Quit()
		}
	}
}
