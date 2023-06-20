package pnet

import (
	"log"
	"net"
	"strconv"
)

func Dial(ip string, port int, chaRecv chan *PRecvData) *PSocket {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		log.Panic(err)
	}

	return PSocket{}.New(conn, chaRecv)
}
