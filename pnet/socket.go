package pnet

import (
	"fmt"
	"net"

	"github.com/caticat/go_game_server/plog"
)

type PSocket struct {
	m_conn              net.Conn
	m_chaSend           chan *PMessage
	m_chaRecv           chan *PRecvData
	m_sessionID         int64 // sessionID/serverID共用字段,同时只能存在一个
	m_connType          int
	m_isInnerConnection bool
}

func NewPSocket(c net.Conn, chaRecv chan *PRecvData) *PSocket {
	t := &PSocket{
		m_conn:              c,
		m_chaSend:           make(chan *PMessage, PSocket_ChanLen),
		m_chaRecv:           chaRecv,
		m_isInnerConnection: false,
	}
	return t
}

func (t *PSocket) Start() {
	go t.runSend()
	go t.runRecv()
}

func (t *PSocket) Send(m *PMessage) {
	if t.getConn() == nil {
		plog.ErrorLn("t.getConn() == nil,socket not connected")
		return
	}
	t.getChaSend() <- m
}

func (t *PSocket) Close() {
	c := t.getConn()
	if c == nil {
		return
	}

	if sm := GetSocketMgr(); sm != nil {
		// debug.PrintStack()
		sm.OnDisconnect(t)
	}

	c.Close()
	t.setConn(nil)
}

func (t *PSocket) String() string {
	addr := "?.?.?.?:?"
	if c := t.getConn(); c != nil {
		addr = c.RemoteAddr().String()
	}

	strID := "SessionID"
	if t.GetIsInnerConnection() {
		strID = "ServerID"
	}

	return fmt.Sprintf("Host:[%q],InInner:%v,%v:%v,ConnectionType:%v",
		addr,
		t.GetIsInnerConnection(),
		strID,
		t.GetSessionID(),
		t.GetConnectionType())
}
