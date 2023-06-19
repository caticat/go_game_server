package main

import (
	"fmt"
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
	t.Regist(int32(ProtoExample.MsgID_HelloReqID), t.helloReqHandler)
}

func (t *MessageManager) helloReqHandler(r *pnet.PRecvData) bool {
	// 参数校验
	s := r.GetSocket()
	if s == nil {
		log.Println("s == nil")
		return false
	}
	m := r.GetMessage()
	if m == nil {
		log.Println("msg == nil")
		return false
	}
	msg := &ProtoExample.HelloReq{}
	m.Unmarshal(msg)

	// 逻辑处理
	fmt.Println("收到消息:", m.GetMsgID(), msg.GetMsg())

	// 返回协议
	a := &ProtoExample.HelloAck{
		Error: ProtoExample.ErrorCode_OK,
		Msg:   "收到消息了:" + msg.GetMsg(),
	}
	s.Send(pnet.PMessage{}.New(int32(ProtoExample.MsgID_HelloAckID), a))

	return false
}