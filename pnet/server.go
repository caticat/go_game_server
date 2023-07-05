package pnet

import (
	"github.com/caticat/go_game_server/plog"
	"github.com/caticat/go_game_server/pnet/conf"
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

func Connect(serverConfigs []*conf.ConfServerRemote) {
	socketMgr := GetSocketMgr()
	if socketMgr == nil {
		plog.PanicLn("socketMgr == nil")
	}
	go runConnect(serverConfigs)
}
