package main

import (
	"log"
	"net"

	"github.com/caticat/go_game_server/example/server/conf"
	"github.com/caticat/go_game_server/pnet"
)

// TODO: 功能待制作

type SocketManager struct {
	m_mapSocketPre map[*pnet.PSocket]bool  // 无sessionID连接
	m_mapSocket    map[int64]*pnet.PSocket // 有sessionID连接
	m_chaRecv      chan *pnet.PRecvData
}

func (t *SocketManager) getChaRecv() chan *pnet.PRecvData { return t.m_chaRecv }

func (t SocketManager) New() *SocketManager {
	t.m_mapSocketPre = make(map[*pnet.PSocket]bool)
	t.m_chaRecv = make(chan *pnet.PRecvData, conf.ChaRecvLen)
	return &t
}

func (t *SocketManager) OnConnect(conn net.Conn) {
	s := pnet.PSocket{}.New(conn, t.getChaRecv())
	s.Start()
	t.m_mapSocketPre[s] = true
}

func (t *SocketManager) Add(s *pnet.PSocket) {
	if s == nil {
		log.Println("s == nil")
		return
	}

	delete(t.m_mapSocketPre, s)
	t.m_mapSocket[s.GetSessionID()] = s
}

func (t *SocketManager) Del(sessionID int64) { delete(t.m_mapSocket, sessionID) }

func (t *SocketManager) Get(sessionID int64) *pnet.PSocket {
	s, ok := t.m_mapSocket[sessionID]
	if !ok {
		log.Printf("sessionID:%v not exist\n", sessionID)
		return nil
	}
	return s
}

func (t *SocketManager) GetAll() map[int64]*pnet.PSocket {
	return t.m_mapSocket
}
