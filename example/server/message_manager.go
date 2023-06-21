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
	t.Regist(int32(pproto.MsgID_InitSessionReqID), t.initSessionReqID)
	t.Regist(int32(pproto.MsgID_InitSessionAckID), t.initSessionAckID)
	t.Regist(int32(pproto.MsgID_HelloReqID), t.helloReqHandler)
}

func (t *MessageManager) initSessionReqID(r *pnet.PRecvData) bool {
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
	msg := &pproto.InitSessionReq{}
	m.Unmarshal(msg)

	if s.GetSessionID() != 0 {
		plog.ErrorLn("s.GetSessionID() != 0,sessionID:", s.GetSessionID())
		return false
	}

	// 逻辑处理
	sessionID := msg.GetSessionID()
	s.SetSessionID(sessionID)
	getSocketManager().Add(s)
	plog.DebugLn("连接初始化:", sessionID)

	// 返回协议
	a := &pproto.InitSessionAck{
		Error:     pproto.ErrorCode_OK,
		SessionID: sessionID,
	}
	s.Send(pnet.NewPMessage(int32(pproto.MsgID_InitSessionAckID), a))

	return false
}

func (t *MessageManager) initSessionAckID(r *pnet.PRecvData) bool {
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
	msg := &pproto.InitSessionAck{}
	m.Unmarshal(msg)

	// 逻辑处理
	plog.DebugLn("收到返回:", m.GetMsgID(), msg.GetError(), msg.GetSessionID())

	// 测试发送协议
	a := &pproto.HelloReq{
		Msg: "收到返回了!!!",
	}
	s.Send(pnet.NewPMessage(int32(pproto.MsgID_HelloReqID), a))

	return false
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
	msg := &pproto.HelloReq{}
	m.Unmarshal(msg)

	// 逻辑处理
	plog.DebugLn("收到消息:", m.GetMsgID(), msg.GetMsg())

	// 返回协议
	a := &pproto.HelloAck{
		Error: pproto.ErrorCode_OK,
		Msg:   "收到消息了:" + msg.GetMsg(),
	}
	s.Send(pnet.NewPMessage(int32(pproto.MsgID_HelloAckID), a))

	return false
}
