package main

import (
	"sync"

	"github.com/caticat/go_game_server/example/proto"
	"github.com/caticat/go_game_server/example/server/conf"
	"github.com/caticat/go_game_server/plog"
	"github.com/caticat/go_game_server/pnet"
)

// TODO: 功能待制作

type SocketManager struct {
	m_mapConnection       map[string]bool                          // 所有连接,是否有连接 <"ip:port", true>
	m_mapSocketPre        map[*pnet.PSocket]bool                   // 无sessionID/serverID连接
	m_mapSocket           map[int64]*pnet.PSocket                  // 有sessionID连接 <sessionID, socket>
	m_mapSocketServer     map[int64]*pnet.PSocket                  // 有serverID连接 <serverID, socket>
	m_mapConnectionSocket map[proto.ConnectionType][]*pnet.PSocket // 有serverID连接 <connectionType, []socket>
	m_chaRecv             chan *pnet.PRecvData
	m_chaMainLoopFun      chan func() // 主线程函数调用

	m_mutConnection sync.Mutex // 连接锁
}

func (t *SocketManager) GetChaRecv() chan *pnet.PRecvData { return t.m_chaRecv }
func (t *SocketManager) GetChaMainLoopFun() chan func()   { return t.m_chaMainLoopFun }
func (t *SocketManager) AddMainLoopFun(fun func())        { t.GetChaMainLoopFun() <- fun }
func (t *SocketManager) getMutConnection() *sync.Mutex    { return &t.m_mutConnection }

func NewSocketManager() *SocketManager {
	t := &SocketManager{
		m_mapConnection:       make(map[string]bool),
		m_mapSocketPre:        make(map[*pnet.PSocket]bool),
		m_mapSocket:           make(map[int64]*pnet.PSocket),
		m_mapSocketServer:     make(map[int64]*pnet.PSocket),
		m_mapConnectionSocket: make(map[proto.ConnectionType][]*pnet.PSocket),
		m_chaRecv:             make(chan *pnet.PRecvData, conf.ChaRecvLen),
		m_chaMainLoopFun:      make(chan func(), conf.ChaMainLoopFun),
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

	plog.InfoLn("socket connect:", s)
	s.Start()

	t.m_mapSocketPre[s] = true
	t.m_mapConnection[s.GetHost()] = true

	// sessionID处理
	t.AddMainLoopFun(func() {
		if s.GetIsInnerConnection() { // 服务器连接
			ct := s.GetConnectionType()
			if ct != int(proto.ConnectionType_ConnectionType_Default) { // 主动连接
				msgID := int32(proto.MsgID_InitConnectionNtfID)
				getMessageManager().Handle(s, msgID, &proto.InitConnectionNtf{ // 本地处理
					ServerID:       s.GetServerID(),
					ConnectionType: proto.ConnectionType(ct),
				})
				c := getConf()
				s.Send(pnet.NewPMessage(msgID, &proto.InitConnectionNtf{ // 远端处理
					ServerID:       c.GetID(),
					ConnectionType: proto.ConnectionType(c.GetConnectionType()),
				}))
			}
		} else { // 客户端连接
			getMessageManager().Handle(s, int32(proto.MsgID_InitSessionNtfID), &proto.InitSessionNtf{ // 本地处理
				SessionID: pnet.GenSessionID(),
			})
		}
	})
}

func (t *SocketManager) OnDisconnect(s *pnet.PSocket) {
	if s == nil {
		plog.ErrorLn("s == nil")
		return
	}

	t.getMutConnection().Lock()
	defer t.getMutConnection().Unlock()

	plog.InfoLn("socket disconnect:", s)
	host := s.GetHost()
	delete(t.m_mapConnection, host)
	delete(t.m_mapSocketPre, s)
	t.AddMainLoopFun(func() { t.Del(s) })
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

	if s.GetIsInnerConnection() {
		t.m_mapSocketServer[s.GetSessionID()] = s
	} else {
		t.m_mapSocket[s.GetSessionID()] = s
	}
	t.m_mapConnectionSocket[proto.ConnectionType(s.GetConnectionType())] = append(t.m_mapConnectionSocket[proto.ConnectionType(s.GetConnectionType())], s)

	plog.InfoLn("socket add:", s)
}

func (t *SocketManager) Del(s *pnet.PSocket) {
	if s == nil {
		plog.ErrorLn("s == nil")
		return
	}

	t.getMutConnection().Lock()
	defer t.getMutConnection().Unlock()

	if s.GetIsInnerConnection() {
		delete(t.m_mapSocketServer, s.GetSessionID())
	} else {
		delete(t.m_mapSocket, s.GetSessionID())
	}

	l := t.m_mapConnectionSocket[proto.ConnectionType(s.GetConnectionType())]
	i := -1
	for m, n := range l {
		if n == s {
			i = m
			break
		}
	}
	if i >= 0 {
		if i+1 >= len(l) { // 最后一个
			l = l[:i]
		} else {
			l = append(l[:i], l[i+1:]...)
		}
		t.m_mapConnectionSocket[proto.ConnectionType(s.GetConnectionType())] = l
	} else {
		plog.ErrorLn("i < 0,connectionType:", s.GetConnectionType(), ",sessionID:", s.GetSessionID())
	}

	plog.InfoLn("socket del:", s)
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

func (t *SocketManager) GetServer(sessionID int64) *pnet.PSocket {
	t.getMutConnection().Lock()
	defer t.getMutConnection().Unlock()
	s, ok := t.m_mapSocketServer[sessionID]
	if !ok {
		plog.ErrorLn("sessionID:%v not exist\n", sessionID)
		return nil
	}
	return s
}

func (t *SocketManager) GetServerAll() map[int64]*pnet.PSocket {
	return t.m_mapSocketServer
}
