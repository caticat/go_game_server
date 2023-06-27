package pnet

import (
	"time"

	"github.com/caticat/go_game_server/plog"
)

var (
	g_socketMgr PSocketManager = nil
	g_sessionID int64          = 0
)

func GetSocketMgr() PSocketManager          { return g_socketMgr }
func setSocketMgr(socketMgr PSocketManager) { g_socketMgr = socketMgr }
func GenSessionID() int64                   { g_sessionID += 1; return g_sessionID }

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

func runConnect(serverConfigs []*ConfRemoteServer) {
	socketMgr := GetSocketMgr()
	if socketMgr == nil {
		plog.PanicLn("socketMgr == nil")
	}

	t := time.Tick(time.Second * 10)
	for range t {
		for _, cfg := range serverConfigs {
			if socketMgr.HasConnect(cfg.String()) {
				continue
			}

			connect(cfg)
		}
	}
}

func connect(cfg *ConfRemoteServer) {
	socketMgr := GetSocketMgr()
	if socketMgr == nil {
		plog.PanicLn("socketMgr == nil")
	}

	plog.DebugLn("try connect to:", cfg)
	s := Dial(cfg.GetIP(), cfg.GetPort(), socketMgr.GetChaRecv())
	if s == nil {
		plog.ErrorLn("Dail failed,cfg:", cfg)
		return
	}
	s.SetServerID(int64(cfg.GetServerID()))
	s.SetIsInnerConnection(true)
	s.SetConnectionType(cfg.GetConnectionType())

	socketMgr.OnConnect(s)
}
