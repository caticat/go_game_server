package main

import (
	ProtoExample "github.com/caticat/go_game_server/example/proto"
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
	t.Regist(int32(ProtoExample.MsgID_HelloReqID), t.helloReqHandler)
}

func (t *MessageManager) helloReqHandler(r *pnet.PRecvData) bool {
	// 参数校验
	s := r.GetSocket()
	if s == nil {
		plog.ErrorLn("s == nil")
		return false
	}
	m := r.GetMessage()
	if m == nil {
		plog.ErrorLn("msg == nil")
		return false
	}
	msg := &ProtoExample.HelloReq{}
	m.Unmarshal(msg)

	// 逻辑处理
	plog.DebugLn("收到消息:", m.GetMsgID(), msg.GetMsg())

	// 返回协议
	a := &ProtoExample.HelloAck{
		Error: ProtoExample.ErrorCode_OK,
		Msg:   "收到消息了:" + msg.GetMsg(),
	}
	s.Send(pnet.NewPMessage(int32(ProtoExample.MsgID_HelloAckID), a))

	return false
}
