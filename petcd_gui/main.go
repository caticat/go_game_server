package main

import (
	"github.com/caticat/go_game_server/petcd"
	"github.com/caticat/go_game_server/plog"
)

func main() {
	plog.Init(plog.ELogLevel_Debug, "")

	conf := getConf()
	if err := conf.Init(); err != nil {
		plog.FatalLn(err)
	}

	petcd.Init(conf.Etcd)
	initData()

	runGUI()

	close()
}

func close() {
	petcd.Close()
}
