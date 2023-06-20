package pnet

import (
	"bufio"
	"io"
	"net"

	"github.com/caticat/go_game_server/plog"
)

const PSocket_ChanLen = 10 // 收发消息阻塞长度

type PSocket struct {
	m_conn      net.Conn
	m_chaSend   chan *PMessage
	m_chaRecv   chan *PRecvData
	m_sessionID int64
}

func (t *PSocket) getConn() net.Conn {
	return t.m_conn
}

func (t *PSocket) getChaSend() chan *PMessage {
	return t.m_chaSend
}

func (t *PSocket) getChaRecv() chan *PRecvData {
	return t.m_chaRecv
}

func (t *PSocket) GetSessionID() int64 { return t.m_sessionID }

func (t *PSocket) SetSessionID(sessionID int64) { t.m_sessionID = sessionID }

func (t PSocket) New(c net.Conn, chaRecv chan *PRecvData) *PSocket {
	t.m_conn = c
	t.m_chaSend = make(chan *PMessage, PSocket_ChanLen)
	t.m_chaRecv = chaRecv
	return &t
}

func (t *PSocket) Start() {
	go t.runSend()
	go t.runRecv()
}

func (t *PSocket) Send(m *PMessage) {
	t.m_chaSend <- m
}

func (t *PSocket) runSend() {
	for {
		select {
		case data := <-t.getChaSend():
			t.send(data)
			// TODO: 是否需要退出循环
		}
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
		chaRecv <- PRecvData{}.New(t, p)
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
	t.getConn().Close()
}
