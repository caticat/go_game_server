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
	t.Regist(int32(pproto.MsgID_TickNtfID), t.tickNtfID)
	t.Regist(int32(pproto.MsgID_InitSessionNtfID), t.initSessionNtfID)
	t.Regist(int32(pproto.MsgID_InitConnectionNtfID), t.initConnectionNtfID)
	t.Regist(int32(pproto.MsgID_HelloReqID), t.helloReqHandler)
}

func (t *MessageManager) tickNtfID(r *pnet.PRecvData) bool {
	// 参数校验
	s := r.GetSocket()
	if s == nil {
		plog.ErrorLn("s == nil")
		return false
	}

	// 逻辑处理
	plog.DebugF("收到服务器[%s]心跳包\n", s)

	return false
}

func (t *MessageManager) initSessionNtfID(r *pnet.PRecvData) bool {
	// 参数校验
	s := r.GetSocket()
	if s == nil {
		plog.ErrorLn("s == nil")
		return false
	}
	if s.GetIsInnerConnection() {
		plog.ErrorLn("s.GetIsInnerConnection(),can't set sessionID")
		return false
	}
	m := r.GetMessage()
	if m == nil {
		plog.ErrorLn("msg == nil")
		return false
	}
	msg := &pproto.InitSessionNtf{}
	m.Unmarshal(msg)

	if s.GetSessionID() != 0 {
		plog.ErrorLn("s.GetSessionID() != 0,sessionID:", s.GetSessionID())
		return false
	}

	// 逻辑处理
	sessionID := msg.GetSessionID()
	plog.DebugLn("收到客户端连接初始化,sessionID:", sessionID)
	s.SetSessionID(sessionID)
	s.SetConnectionType(int(pproto.ConnectionType_ConnectionType_Client))
	getSocketManager().Add(s)

	return false
}

func (t *MessageManager) initConnectionNtfID(r *pnet.PRecvData) bool {
	// 参数校验
	s := r.GetSocket()
	if s == nil {
		plog.ErrorLn("s == nil")
		return false
	}
	if !s.GetIsInnerConnection() {
		plog.ErrorLn("!s.GetIsInnerConnection(),can't set serverID")
		return false
	}
	m := r.GetMessage()
	if m == nil {
		plog.ErrorLn("msg == nil")
		return false
	}
	msg := &pproto.InitConnectionNtf{}
	m.Unmarshal(msg)

	// 逻辑处理
	serverID := msg.GetServerID()
	connectionType := msg.GetConnectionType()
	plog.DebugF("收到服务端连接初始化,serverID:%v,connectionType:%v\n", serverID, connectionType)
	s.SetServerID(serverID)
	s.SetConnectionType(int(connectionType))
	getSocketManager().Add(s)

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
