package main

import (
	"strings"
	"time"

	pproto "github.com/caticat/go_game_server/example/proto"
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
	g_messageManager                      = NewMessageManager()
)

func main() {
	// 初始化
	plog.Init(plog.ELogLevel_Debug, "")
	getMessageManager().Init()

	// 连接
	g_s = pnet.Dial("127.0.0.1", 6666, getChaRecv())
	if g_s == nil {
		plog.PanicLn("g_s == nil")
	}
	g_s.Start()
	defer g_s.Close()

	// 收取协议
	go run()

	// 发送协议
	for i := 0; i < 10; i++ {
		msg := &pproto.HelloReq{
			Msg: strings.Repeat("a", i),
		}
		d := pnet.NewPMessage(int32(pproto.MsgID_HelloReqID), msg)
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
