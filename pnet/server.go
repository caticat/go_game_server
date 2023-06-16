package pnet

import (
	"fmt"
	"log"
	"net"
)

func ListenAndServe(port int, socketMgr PSocketManager) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		socketMgr.OnConnect(conn)
	}
}
