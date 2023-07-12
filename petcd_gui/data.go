package main

import (
	"github.com/caticat/go_game_server/petcd"
	"github.com/caticat/go_game_server/petcd/pdata"
)

func initData() {
	mapResult := make(map[string]string)
	petcd.GetPrefix(pdata.PDATA_PREFIX, mapResult)
	getEtcdData().SetAll(mapResult)
}
