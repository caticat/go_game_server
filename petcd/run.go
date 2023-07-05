package petcd

import (
	"github.com/caticat/go_game_server/plog"
	"go.etcd.io/etcd/clientv3"
)

func run(cha <-chan *clientv3.LeaseKeepAliveResponse) {
	plog.InfoLn("petcd run begin")

	for range cha {
		// plog.DebugLn("run tick")
	}

	plog.InfoLn("petcd run end")
}

func runWatch(cha clientv3.WatchChan, prefix string, fun funWatchCallback_t) {
	plog.InfoLn("petcd runWatch begin")

	for resp := range cha {
		for _, e := range resp.Events {

			fun(e.Type, prefix, e.Kv)
		}
	}

	plog.InfoLn("petcd runWatch end")
}
