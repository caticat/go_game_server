package pnet

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/caticat/go_game_server/plog"
)

const PSocket_ChanLen = 10 // 收发消息阻塞长度

type PSocket struct {
	m_conn              net.Conn
	m_chaSend           chan *PMessage
	m_chaRecv           chan *PRecvData
	m_sessionID         int64 // sessionID/serverID共用字段,同时只能存在一个
	m_connType          int
	m_isInnerConnection bool
}

func (t *PSocket) getConn() net.Conn                     { return t.m_conn }
func (t *PSocket) getChaSend() chan *PMessage            { return t.m_chaSend }
func (t *PSocket) getChaRecv() chan *PRecvData           { return t.m_chaRecv }
func (t *PSocket) GetSessionID() int64                   { return t.m_sessionID }
func (t *PSocket) SetSessionID(sessionID int64)          { t.m_sessionID = sessionID }
func (t *PSocket) GetServerID() int64                    { return t.m_sessionID }
func (t *PSocket) SetServerID(serverID int64)            { t.m_sessionID = serverID }
func (t *PSocket) GetHost() string                       { return t.m_conn.RemoteAddr().String() }
func (t *PSocket) GetConnectionType() int                { return t.m_connType }
func (t *PSocket) SetConnectionType(connType int)        { t.m_connType = connType }
func (t *PSocket) GetIsInnerConnection() bool            { return t.m_isInnerConnection }
func (t *PSocket) SetIsInnerConnection(isInnerConn bool) { t.m_isInnerConnection = isInnerConn }

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

func (t *PSocket) Send(m *PMessage) { t.getChaSend() <- m }

func (t *PSocket) runSend() {
	for data := range t.getChaSend() {
		t.send(data)
	}

	// for {
	// 	select {
	// 	case data := <-t.getChaSend():
	// 		t.send(data)
	// 		// TODO: 是否需要退出循环
	// 	}
	// }
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

func (t *PSocket) Close() {
	GetSocketMgr().OnDisconnect(t)
	t.getConn().Close()
}

func (t *PSocket) String() string {
	strID := "SessionID"
	if t.GetIsInnerConnection() {
		strID = "ServerID"
	}
	return fmt.Sprintf("Host:[%q],%v:%v,InInner:%v,ConnectionType:%v",
		t.getConn().RemoteAddr().String(),
		strID,
		t.GetSessionID(),
		t.GetIsInnerConnection(),
		t.GetConnectionType())
}
