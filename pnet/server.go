package pnet

import (
	"fmt"
	"net"

	"github.com/caticat/go_game_server/plog"
)

func Init(socketMgr PSocketManager) { setSocketMgr(socketMgr) }

func ListenAndServe(port int) {
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

		socketMgr.OnConnect(NewPSocket(conn, socketMgr.GetChaRecv()))
	}
}

func Connect(serverConfigs []*ConfRemoteServer) {
	socketMgr := GetSocketMgr()
	if socketMgr == nil {
		plog.PanicLn("socketMgr == nil")
	}
	go runConnect(serverConfigs)
}
