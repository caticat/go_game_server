package main

import (
	"net"

	"github.com/caticat/go_game_server/example/server/conf"
	"github.com/caticat/go_game_server/plog"
	"github.com/caticat/go_game_server/pnet"
)

// TODO: 功能待制作

type SocketManager struct {
	m_mapSocketPre map[*pnet.PSocket]bool  // 无sessionID连接
	m_mapSocket    map[int64]*pnet.PSocket // 有sessionID连接
	m_chaRecv      chan *pnet.PRecvData
}

func (t *SocketManager) getChaRecv() chan *pnet.PRecvData { return t.m_chaRecv }

func NewSocketManager() *SocketManager {
	t := &SocketManager{
		m_mapSocketPre : make(map[*pnet.PSocket]bool),
		m_chaRecv : make(chan *pnet.PRecvData, conf.ChaRecvLen)
	}
	return &t
}

func (t *SocketManager) OnConnect(conn net.Conn) {
	s := pnet.NewPSocket(conn, t.getChaRecv())
	s.Start()
	t.m_mapSocketPre[s] = true
}

func (t *SocketManager) Add(s *pnet.PSocket) {
	if s == nil {
		plog.ErrorLn("s == nil")
		return
	}

	delete(t.m_mapSocketPre, s)
	t.m_mapSocket[s.GetSessionID()] = s
}

func (t *SocketManager) Del(sessionID int64) { delete(t.m_mapSocket, sessionID) }

func (t *SocketManager) Get(sessionID int64) *pnet.PSocket {
	s, ok := t.m_mapSocket[sessionID]
	if !ok {
		plog.ErrorLn("sessionID:%v not exist\n", sessionID)
		return nil
	}
	return s
}

func (t *SocketManager) GetAll() map[int64]*pnet.PSocket {
	return t.m_mapSocket
}
