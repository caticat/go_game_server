package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/caticat/go_game_server/petcd"
	"github.com/caticat/go_game_server/plog"
	"github.com/caticat/go_game_server/pnet"
)

var (
	g_socketManager  = NewSocketManager()
	g_messageManager = NewMessageManager()
	g_conf           = NewConfServer()
	g_chaSig         = make(chan os.Signal, 1)
)

func main() {
	// 初始化
	initServer()

	// 主线程
	run()
}

func initServer() {
	// 信号
	signal.Notify(g_chaSig, os.Interrupt)

	// 配置
	c := getConf()
	c.Init()

	// 日志
	plog.Init(c.GetLog().GetLevel(), c.GetLog().GetFile())

	// 协议
	getMessageManager().Init()

	// 定时器
	initTimer()

	// 网络
	pnet.Init(getSocketManager())
	pnet.Connect(getConf().GetRemoteServers())
	go pnet.ListenAndServe(getConf().GetPort(), getConf().GetPortIn())
}

func run() {
	t := time.Tick(50 * time.Millisecond)
	chaRecv := getSocketManager().GetChaRecv()
	chaFun := getSocketManager().GetChaMainLoopFun()
	for {
		select {
		case r := <-chaRecv:
			getMessageManager().Trigger(r)
		case f := <-chaFun:
			f()
		case <-t:
			runTimer(time.Now().Local().UnixMilli())
		case s := <-g_chaSig:
			onExit(s)
			return
		}
	}
}

func onExit(s os.Signal) {
	plog.InfoLn("receive signal:", s)
	getSocketManager().Close()
	petcd.Close()
}

func getSocketManager() *SocketManager   { return g_socketManager }
func getMessageManager() *MessageManager { return g_messageManager }
func getConf() *ConfServer               { return g_conf }
