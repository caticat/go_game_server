package main

import (
	"sync"

	"github.com/caticat/go_game_server/example/server/conf"
	"github.com/caticat/go_game_server/plog"
	"github.com/caticat/go_game_server/pnet"
)

// TODO: 功能待制作

type SocketManager struct {
	m_mapConnection map[string]bool         // 所有连接,是否有连接 <"ip:port", true>
	m_mapSocketPre  map[*pnet.PSocket]bool  // 无sessionID连接
	m_mapSocket     map[int64]*pnet.PSocket // 有sessionID连接
	m_chaRecv       chan *pnet.PRecvData

	m_mutConnection sync.Mutex // 连接锁
}

func (t *SocketManager) GetChaRecv() chan *pnet.PRecvData { return t.m_chaRecv }
func (t *SocketManager) getMutConnection() *sync.Mutex    { return &t.m_mutConnection }

func NewSocketManager() *SocketManager {
	t := &SocketManager{
		m_mapConnection: make(map[string]bool),
		m_mapSocketPre:  make(map[*pnet.PSocket]bool),
		m_chaRecv:       make(chan *pnet.PRecvData, conf.ChaRecvLen),
	}
	return t
}

func (t *SocketManager) OnConnect(s *pnet.PSocket) {
	if s == nil {
		plog.ErrorLn("s == nil")
		return
	}

	t.getMutConnection().Lock()
	defer t.getMutConnection().Unlock()

	plog.InfoLn("connect to:", s.GetHost(), "done!")
	s.Start()

	t.m_mapSocketPre[s] = true
	t.m_mapConnection[s.GetHost()] = true
}

func (t *SocketManager) OnDisconnect(s *pnet.PSocket) {
	if s == nil {
		plog.ErrorLn("s == nil")
		return
	}

	t.getMutConnection().Lock()
	defer t.getMutConnection().Unlock()

	plog.InfoLn("server disconnect:", s.GetHost())
	host := s.GetHost()
	delete(t.m_mapConnection, host)
	delete(t.m_mapSocketPre, s)
	t.Del(s.GetSessionID())
}

func (t *SocketManager) HasConnect(host string) bool {
	t.getMutConnection().Lock()
	defer t.getMutConnection().Unlock()

	_, ok := t.m_mapConnection[host]
	return ok
}

func (t *SocketManager) Add(s *pnet.PSocket) {
	if s == nil {
		plog.ErrorLn("s == nil")
		return
	}

	t.getMutConnection().Lock()
	defer t.getMutConnection().Unlock()

	delete(t.m_mapSocketPre, s)
	t.m_mapSocket[s.GetSessionID()] = s
}

func (t *SocketManager) Del(sessionID int64) {
	if sessionID <= 0 {
		return
	}

	// 这里不能加锁,调用处已经加过了
	// t.getMutConnection().Lock()
	// defer t.getMutConnection().Unlock()

	delete(t.m_mapSocket, sessionID)
}

func (t *SocketManager) Get(sessionID int64) *pnet.PSocket {
	t.getMutConnection().Lock()
	defer t.getMutConnection().Unlock()

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
