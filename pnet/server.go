package pnet

import (
	"fmt"
	"log"
	"net"
)

func ListenAndServe(port int, fun func(net.Conn)) {
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
		fun(conn)
	}
}
