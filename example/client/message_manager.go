package main

import (
	"log"

	ProtoExample "github.com/caticat/go_game_server/example/proto"
	"github.com/caticat/go_game_server/pnet"
)

type MessageManager struct {
	*pnet.PMessageManager
}

func (t MessageManager) New() *MessageManager {
	t.PMessageManager = pnet.PMessageManager{}.New()
	return &t
}

func (t *MessageManager) Init() {
	// 消息待注册
	t.Regist(int32(ProtoExample.MsgID_HelloAckID), t.helloAckHandler)
}

func (t *MessageManager) helloAckHandler(r *pnet.PRecvData) bool {
	m := r.GetMessage()
	if m == nil {
		log.Println("msg == nil")
		return false
	}
	msg := &ProtoExample.HelloAck{}
	m.Unmarshal(msg)
	log.Println("收到消息:", m.GetMsgID(), msg.GetError(), msg.GetMsg())

	return false
}
