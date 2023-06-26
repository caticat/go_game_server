package pnet

import (
	"fmt"
	"net"

	"github.com/caticat/go_game_server/plog"
)

func Init(socketMgr PSocketManager) { setSocketMgr(socketMgr) }

func ListenAndServe(port int, portIn int) {
	if port == 0 && portIn == 0 {
		plog.PanicLn("port == 0 && portIn == 0")
	}

	if portIn > 0 {
		if port == 0 {
			listenAndServe(portIn, true)
		} else {
			go listenAndServe(portIn, true)
		}
	}

	if port > 0 {
		listenAndServe(port, false)
	}
}

func listenAndServe(port int, isInnerConnection bool) {
	socketMgr := GetSocketMgr()
	if socketMgr == nil {
		plog.PanicLn("socketMgr == nil")
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		plog.PanicLn(err)
	}

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

func Connect(serverConfigs []*ConfRemoteServer) {
	socketMgr := GetSocketMgr()
	if socketMgr == nil {
		plog.PanicLn("socketMgr == nil")
	}
	go runConnect(serverConfigs)
}
