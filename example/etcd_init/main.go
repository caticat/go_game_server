package main

// 初始化ETCD的数据

import (
	"flag"
	"os"

	"github.com/caticat/go_game_server/petcd"
	"github.com/caticat/go_game_server/plog"

	"gopkg.in/yaml.v3"
)

var g_fileConfig string

func main() {
	plog.Init(plog.ELogLevel_Debug, "")
	flag.StringVar(&g_fileConfig, "c", "etcd_init.yaml", "-c=etcd_init.yaml")
	flag.Parse()

	f, err := os.ReadFile(g_fileConfig)
	if err != nil {
		plog.FatalLn("error:", err)
	}

	cfg := petcd.NewConfigEtcdInit()
	err = yaml.Unmarshal(f, cfg)
	if err != nil {
		plog.FatalLn("error:", err)
	}

	err = petcd.ProcessEtcdInit(cfg)
	if err != nil {
		plog.FatalLn("error:", err)
	}
	plog.InfoLn("init etcd done")
}
