package pnet

import (
	"fmt"
	"net"

	"github.com/caticat/go_game_server/plog"
)

func ListenAndServe(port int, socketMgr PSocketManager) {
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
		socketMgr.OnConnect(conn)
	}
}
