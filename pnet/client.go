package pnet

import (
	"net"
	"strconv"
	"time"

	"github.com/caticat/go_game_server/plog"
)

func Dial(ip string, port int, chaRecv chan *PRecvData) *PSocket {
	conn, err := net.DialTimeout("tcp", ip+":"+strconv.Itoa(port), time.Second)
	if err != nil {
		plog.ErrorLn(err)
		return nil
	}

	return NewPSocket(conn, chaRecv)
}
