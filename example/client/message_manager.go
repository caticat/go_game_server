package main

import (
	pproto "github.com/caticat/go_game_server/example/proto"
	"github.com/caticat/go_game_server/plog"
	"github.com/caticat/go_game_server/pnet"
)

type MessageManager struct {
	*pnet.PMessageManager
}

func NewMessageManager() *MessageManager {
	t := &MessageManager{
		PMessageManager: pnet.NewPMessageManager(),
	}
	return t
}

func (t *MessageManager) Init() {
	// 消息待注册
	t.Regist(int32(pproto.MsgID_HelloAckID), t.helloAckHandler)
}

func (t *MessageManager) helloAckHandler(r *pnet.PRecvData) bool {
	m := r.GetMessage()
	if m == nil {
		plog.ErrorLn("msg == nil")
		return false
	}
	msg := &pproto.HelloAck{}
	m.Unmarshal(msg)
	plog.InfoLn("收到消息:", m.GetMsgID(), msg.GetError(), msg.GetMsg())

	return false
}
