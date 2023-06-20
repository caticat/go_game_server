package main

import (
	"time"

	"github.com/caticat/go_game_server/example/server/conf"
	"github.com/caticat/go_game_server/plog"
	"github.com/caticat/go_game_server/pnet"
)

var (
	g_socketManager  = SocketManager{}.New()
	g_messageManager = MessageManager{}.New()
	g_conf           = conf.ConfServer{}.New()
)

func main() {
	c := getConf()
	c.Init()
	plog.Init(c.GetLog().GetLevel(), c.GetLog().GetFile())
	getMessageManager().Init()

	go run()
	pnet.ListenAndServe(getConf().GetPort(), getSocketManager())
}

func run() {
	t := time.Tick(100 * time.Millisecond)
	chaRecv := getSocketManager().getChaRecv()
	for {
		select {
		case r := <-chaRecv:
			getMessageManager().Trigger(r)
		case <-t:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func getSocketManager() *SocketManager   { return g_socketManager }
func getMessageManager() *MessageManager { return g_messageManager }
func getConf() *conf.ConfServer          { return g_conf }
