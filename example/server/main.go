package main

import (
	"fmt"
	"net"
	"time"

	"github.com/caticat/go_game_server/plog"
	"github.com/caticat/go_game_server/pnet"
)

const (
	PPORT      = 6666
	ChaRecvLen = 100
)

var (
	g_s       *pnet.PSocket
	g_chaRecv chan *pnet.PMessage = make(chan *pnet.PMessage, ChaRecvLen)
)

func main() {
	plog.LogInit()

	go run()
	pnet.ListenAndServe(PPORT, handleConnection)
}

func handleConnection(conn net.Conn) {
	g_s = pnet.PSocket{}.New(conn, g_chaRecv)
	g_s.Start()
}

func run() {
	t := time.Tick(100 * time.Millisecond)
	for {
		select {
		case m := <-g_chaRecv:
			fmt.Println("recv msg:", m)
		case <-t:
			time.Sleep(50 * time.Millisecond)
		}
	}
}
