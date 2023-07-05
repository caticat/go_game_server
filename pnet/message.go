package pnet

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/caticat/go_game_server/plog"
	"google.golang.org/protobuf/proto"
)

type PMessage struct {
	m_msgID   int32
	m_msgData []byte

	m_sessionID int64 // 会话唯一ID
}

func NewPMessage(msgID int32, msg proto.Message) *PMessage {
	t := &PMessage{
		m_msgID: msgID,
	}
	msgData, err := proto.Marshal(msg)
	if err != nil {
		plog.ErrorLn("proto.Marshal failed,error:", err)
	} else {
		t.m_msgData = msgData
	}
	return t
}

func (t *PMessage) GetMsgID() int32              { return t.m_msgID }
func (t *PMessage) GetSessionID() int64          { return t.m_sessionID }
func (t *PMessage) setSessionID(sessionID int64) { t.m_sessionID = sessionID }

func (t PMessage) NewByHead(bufferHead []byte, sessionID int64) (*PMessage, int32) {
	buffer := bytes.NewReader(bufferHead)

	var l int32 = 0
	binary.Read(buffer, binary.BigEndian, &l)
	binary.Read(buffer, binary.BigEndian, &(t.m_msgID))

	t.setSessionID(sessionID)

	return &t, l
}

func (t *PMessage) SetMsgData(data string) {
	t.m_msgData = []byte(data)
}

func (t *PMessage) Marshal() string {
	// 消息体长度,不包含消息头
	var l int32 = int32(len(t.m_msgData))

	// 消息头
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, &l)
	binary.Write(buffer, binary.BigEndian, &(t.m_msgID))

	// 消息体
	buffer.Write(t.m_msgData)

	return buffer.String()
}

func (t *PMessage) Unmarshal(msg proto.Message) {
	proto.Unmarshal(t.m_msgData, msg)
}

func (t *PMessage) String() string {
	return fmt.Sprintf("msgID:%d,msgData:%q", t.m_msgID, t.m_msgData)
}
