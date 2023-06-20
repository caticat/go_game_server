package pnet

import (
	"net"
	"strconv"

	"github.com/caticat/go_game_server/plog"
)

func Dial(ip string, port int, chaRecv chan *PRecvData) *PSocket {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		plog.PanicLn(err)
	}

	return NewPSocket(conn, chaRecv)
}
