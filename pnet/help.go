package pnet

import (
	"net"

	"github.com/caticat/go_game_server/plog"
)

// 获取本机IPv4地址
func GetIPs() []string {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		plog.ErrorLn(err)
		return ips
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips
}
