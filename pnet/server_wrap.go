package pnet

import (
	"fmt"
	"net"
	"time"

	"github.com/caticat/go_game_server/plog"
)

func listenAndServe(port int, isInnerConnection bool) {
	socketMgr := GetSocketMgr()
	if socketMgr == nil {
		plog.PanicLn("socketMgr == nil")
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		plog.PanicLn(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			plog.ErrorLn(err)
			continue
		}

		s := NewPSocket(conn, socketMgr.GetChaRecv())
		if isInnerConnection {
			s.SetIsInnerConnection(true)
		}
		socketMgr.OnConnect(s)
	}
}

func runConnect(serverConfigs []*ConfServerRemote) {
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

func connect(cfg *ConfServerRemote) {
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
