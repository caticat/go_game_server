package pnet

import "net"

type PSocketManager interface {
	OnConnect(net.Conn)
	Add(*PSocket)
	Del(int64)
	Get(int64) *PSocket
	GetAll() map[int64]*PSocket
}
