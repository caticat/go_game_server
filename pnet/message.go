package pnet

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PMessage struct {
	MsgID   int32
	MsgData string // TODO: 临时类型,待修改
}

func (t PMessage) New(msgID int32, msgData string) *PMessage {
	t.MsgID = msgID
	t.MsgData = msgData
	return &t
}

func (t PMessage) NewByHead(bufferHead []byte) (*PMessage, int32) {
	buffer := bytes.NewReader(bufferHead)

	var l int32 = 0
	binary.Read(buffer, binary.BigEndian, &l)
	binary.Read(buffer, binary.BigEndian, &(t.MsgID))

	return &t, l
}

func (t *PMessage) ParseFromString(data string) {
	t.MsgData = data // TODO: 正常解析待制作
}

func (t *PMessage) SerializeAsString() string {
	var l int32 = int32(len(t.MsgData)) // 消息体长度,不包含消息头

	// 消息头
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, &l)
	binary.Write(buffer, binary.BigEndian, &(t.MsgID))

	// 消息体
	buffer.Write([]byte(t.MsgData))

	return buffer.String()
}

func (t *PMessage) String() string {
	return fmt.Sprintf("msgID:%d,msgData:%q", t.MsgID, t.MsgData)
}
