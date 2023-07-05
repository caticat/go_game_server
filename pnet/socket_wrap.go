package pnet

import (
	"bufio"
	"io"
	"net"
	"runtime/debug"

	"github.com/caticat/go_game_server/plog"
)

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
		p, l := NewPMessageByHead(bufferHead, t.GetSessionID())
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
