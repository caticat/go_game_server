package main

import (
	"fmt"
	"time"

	"github.com/caticat/go_game_server/plog"
	"github.com/caticat/go_game_server/pnet"
)

const (
	PPORT      = 6666
	ChaRecvLen = 100
)

var (
	g_socketManager  = SocketManager{}.New()
	g_messageManager = MessageManager{}.New()
)

func main() {
	plog.LogInit()
	getMessageManager().Init()

	go run()
	pnet.ListenAndServe(PPORT, getSocketManager())
}

func run() {
	t := time.Tick(100 * time.Millisecond)
	chaRecv := getSocketManager().getChaRecv()
	for {
		select {
		case m := <-chaRecv:
			getMessageManager().Trigger()
			fmt.Println("recv msg:", m)
		case <-t:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func getSocketManager() *SocketManager   { return g_socketManager }
func getMessageManager() *MessageManager { return g_messageManager }
