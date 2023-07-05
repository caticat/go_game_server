package pnet

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"runtime/debug"

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
		debug.PrintStack()
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

// =============================================================================

func (t *PSocket) getConn() net.Conn                     { return t.m_conn }
func (t *PSocket) setConn(c net.Conn)                    { t.m_conn = c }
func (t *PSocket) getChaSend() chan *PMessage            { return t.m_chaSend }
func (t *PSocket) getChaRecv() chan *PRecvData           { return t.m_chaRecv }
func (t *PSocket) GetSessionID() int64                   { return t.m_sessionID }
func (t *PSocket) SetSessionID(sessionID int64)          { t.m_sessionID = sessionID }
func (t *PSocket) GetServerID() int64                    { return t.m_sessionID }
func (t *PSocket) SetServerID(serverID int64)            { t.m_sessionID = serverID }
func (t *PSocket) GetConnectionType() int                { return t.m_connType }
func (t *PSocket) SetConnectionType(connType int)        { t.m_connType = connType }
func (t *PSocket) GetIsInnerConnection() bool            { return t.m_isInnerConnection }
func (t *PSocket) SetIsInnerConnection(isInnerConn bool) { t.m_isInnerConnection = isInnerConn }
func (t *PSocket) GetHost() string {
	if c := t.m_conn; c != nil {
		return c.RemoteAddr().String()
	} else {
		debug.PrintStack()
		plog.ErrorLn("t.m_conn == nil,sessionID:", t.GetSessionID(), ",connType:", t.GetConnectionType())
		return ""
	}
}

// =============================================================================

func (t *PSocket) runSend() {
	for data := range t.getChaSend() {
		t.send(data)
	}
}

func (t *PSocket) runRecv() {
	r := bufio.NewReader(t.getConn())
	bufferHead := make([]byte, 8)
	chaRecv := t.getChaRecv()

	for {
		// 读取数据头
		lh, err := io.ReadFull(r, bufferHead)
		if err != nil {
			if lh != 0 {
				plog.InfoLn("conn read failed, head, err:", err)
			}
			t.Close()
			break
		}

		// 读取数据体
		p, l := PMessage{}.NewByHead(bufferHead, t.GetSessionID())
		bufferBody := make([]byte, l)
		_, err = io.ReadFull(r, bufferBody)
		if err != nil {
			plog.InfoLn("conn read failed, body, err:", err)
			t.Close()
			break
		}
		p.SetMsgData(string(bufferBody))

		// 返回
		chaRecv <- NewPRecvData(t, p)
	}
}

func (t *PSocket) send(data *PMessage) {
	// bufio.NewReader(t.getConn())
	sliData := []byte(data.Marshal())

	conn := t.getConn()
	lenWaitSend := len(sliData)
	lenDoneSend := 0
	lenSend := 0
	var err error
	for lenWaitSend > lenDoneSend {
		lenSend, err = conn.Write(sliData[lenDoneSend:])
		if err != nil {
			plog.InfoLn("conn write failed, err:", err)
			t.Close()
			break
		}
		lenDoneSend += lenSend
	}
}
