package pnet

import "net"

type PSocketManager interface {
	OnConnect(net.Conn)
	Add(*PSocket)
	Del(int64)
	Get(int64) *PSocket
	GetAll() map[int64]*PSocket
}

type PRecvData struct {
	m_socket  *PSocket
	m_message *PMessage
}

func NewPRecvData(s *PSocket, m *PMessage) *PRecvData {
	t := &PRecvData{
		m_socket:  s,
		m_message: m,
	}
	return t
}
func (t *PRecvData) GetSocket() *PSocket   { return t.m_socket }
func (t *PRecvData) GetMessage() *PMessage { return t.m_message }
