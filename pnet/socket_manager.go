package pnet

type PSocketManager interface {
	GetChaRecv() chan *PRecvData
	OnConnect(*PSocket)
	OnDisconnect(*PSocket)
	HasConnect(host string) bool
	Add(*PSocket)
	Del(*PSocket)
	Get(int64) *PSocket
	GetAll() map[int64]*PSocket
	GetServer(int64) *PSocket
	GetServerAll() map[int64]*PSocket
	Close()
}
