package main

import (
	"strings"
	"time"

	ProtoExample "github.com/caticat/go_game_server/example/proto"
	"github.com/caticat/go_game_server/plog"
	"github.com/caticat/go_game_server/pnet"
)

const (
	PPORT      = 6666
	ChaRecvLen = 100
)

var (
	g_s              *pnet.PSocket
	g_chaRecv        chan *pnet.PRecvData = make(chan *pnet.PRecvData, ChaRecvLen)
	g_messageManager                      = MessageManager{}.New()
)

func main() {
	// 初始化
	plog.LogInit()
	getMessageManager().Init()

	// 连接
	g_s = pnet.Dial("127.0.0.1", 6666, getChaRecv())
	g_s.Start()

	// 收取协议
	go run()

	// 发送协议
	for i := 0; i < 10; i++ {
		msg := &ProtoExample.HelloReq{
			Msg: strings.Repeat("a", i),
		}
		d := pnet.PMessage{}.New(int32(ProtoExample.MsgID_HelloReqID), msg)
		g_s.Send(d)
		time.Sleep(time.Second)
	}
}

func run() {
	t := time.Tick(100 * time.Millisecond)
	chaRecv := getChaRecv()
	for {
		select {
		case r := <-chaRecv:
			getMessageManager().Trigger(r)
		case <-t:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func getChaRecv() chan *pnet.PRecvData   { return g_chaRecv }
func getMessageManager() *MessageManager { return g_messageManager }
