package pnet

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PMessage struct {
	m_msgID   int32
	m_msgData string // TODO: 临时类型,待修改

	m_sessionID int64 // 会话唯一ID
}

func (t *PMessage) GetMsgID() int32              { return t.m_msgID }
func (t *PMessage) GetSessionID() int64          { return t.m_sessionID }
func (t *PMessage) setSessionID(sessionID int64) { t.m_sessionID = sessionID }

func (t PMessage) New(msgID int32, msgData string) *PMessage {
	t.m_msgID = msgID
	t.m_msgData = msgData
	return &t
}

func (t PMessage) NewByHead(bufferHead []byte, sessionID int64) (*PMessage, int32) {
	buffer := bytes.NewReader(bufferHead)

	var l int32 = 0
	binary.Read(buffer, binary.BigEndian, &l)
	binary.Read(buffer, binary.BigEndian, &(t.m_msgID))

	t.setSessionID(sessionID)

	return &t, l
}

func (t *PMessage) ParseFromString(data string) {
	t.m_msgData = data // TODO: 正常解析待制作
}

func (t *PMessage) SerializeAsString() string {
	var l int32 = int32(len(t.m_msgData)) // 消息体长度,不包含消息头

	// 消息头
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, &l)
	binary.Write(buffer, binary.BigEndian, &(t.m_msgID))

	// 消息体
	buffer.Write([]byte(t.m_msgData))

	return buffer.String()
}

func (t *PMessage) String() string {
	return fmt.Sprintf("msgID:%d,msgData:%q", t.m_msgID, t.m_msgData)
}
