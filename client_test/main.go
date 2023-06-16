package main

import (
	"fmt"
	"strings"
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

	g_s = pnet.Dial("127.0.0.1", 6666, g_chaRecv)
	g_s.Start()

	for i := 0; i < 10; i++ {
		d := pnet.PMessage{}.New(int32(1000+i), strings.Repeat("a", i))
		fmt.Println("send data:", d)
		g_s.Send(d)
		time.Sleep(time.Second)
	}
}
