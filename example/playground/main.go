package main

import (
	"context"
	"time"

	"github.com/caticat/go_game_server/petcd"
	"github.com/caticat/go_game_server/plog"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
)

func main() {
	plog.Init(plog.ELogLevel_Debug, "")

	// testETCD()
	// testETCDLease()
	// testPETCD()
	// testPETCDWatch()
	// testPETCDFlushDB()
	// testPETCDCompact()
}

func testETCD() {
	plog.InfoLn("init client")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:60001", "http://localhost:60002", "http://localhost:60003"},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		plog.FatalLn("error:", err)
	}
	defer cli.Close()

	plog.InfoLn("begin put")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	_, err = cli.Put(ctx, "/tmp", "tmp_value1")
	cancel()
	if err != nil {
		plog.FatalLn("error:", err)
	}

	plog.InfoLn("begin get")
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*1)
	// resp, err := cli.Get(ctx, "/tmp")
	resp, err := cli.Get(ctx, "/", clientv3.WithFromKey())
	cancel()
	if err != nil {
		plog.FatalLn("error:", err)
	}
	for i, p := range resp.Kvs {
		plog.InfoF("[%d]%q->%q\n", i, string(p.Key), string(p.Value))
	}
}

func testETCDLease() {
	plog.InfoLn("init client")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		plog.FatalLn("error:", err)
	}
	defer cli.Close()

	const (
		etcdTimeout            = time.Second
		etcdTimeoutLease int64 = 10
	)

	// 初始化租约
	plog.InfoLn("begin put")
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeout)
	resp, err := cli.Lease.Grant(ctx, etcdTimeoutLease)
	cancel()
	if err != nil {
		plog.FatalLn("error:", err)
	}
	leaseID := resp.ID

	// 初始化值
	ctx, cancel = context.WithTimeout(context.Background(), etcdTimeout)
	cli.Put(ctx, "/test", "test_value", clientv3.WithLease(leaseID))
	cancel()

	// 租约续期
	ctx, cancel = context.WithCancel(context.Background())
	r, err := cli.KeepAlive(ctx, leaseID)
	if err != nil {
		plog.FatalLn("error:", err)
	}
	go func(r <-chan *clientv3.LeaseKeepAliveResponse) {
		for range r {
			plog.InfoF("goroutine tick")
		}
		plog.InfoF("goroutine done")
	}(r)

	time.Sleep(time.Second * 10)
	cancel()

	plog.InfoLn("cancel done")
	time.Sleep(time.Second * 60)
	plog.InfoLn("progress done")
}

func testPETCD() {
	t := petcd.NewConfigEtcd()
	t.Endpoints = append(t.Endpoints, "http://127.0.0.1:2379")
	t.DialTimeout = 1
	t.OperationTimeout = 1
	t.LeaseTimeoutBeforeKeepAlive = 30
	petcd.Init(t)
	defer petcd.Close()

	// plog.InfoLn(petcd.GetString("/hello"))
	// plog.InfoLn(petcd.GetString("/tmp"))
	// plog.InfoLn(petcd.GetString("/abc"))

	// // petcd.Put("/abc", "def")
	// petcd.PutAlive("/abc", "ghi1")
	// time.Sleep(time.Second * 10)
	// plog.InfoLn("test done")

	// petcd.PutAlive("/server/127.0.0.1:1", "alive")
	// petcd.PutAlive("/server/127.0.0.1:2", "alive")
	// petcd.PutAlive("/server/127.0.0.1:3", "alive")

	// mapServer := make(map[string]string)
	// petcd.GetPrefix("/server/", mapServer)
	// for k, v := range mapServer {
	// 	plog.InfoF("%v[%v]->%v\n", k, petcd.TrimPrefix(k, "/server/"), v)
	// }

	// time.Sleep(time.Second * 10)
}

func testPETCDWatch() {
	t := petcd.NewConfigEtcd()
	t.Endpoints = append(t.Endpoints, "http://localhost:60001", "http://localhost:60002", "http://localhost:60003")
	t.DialTimeout = 1
	t.OperationTimeout = 1
	t.LeaseTimeoutBeforeKeepAlive = 10
	petcd.Init(t)
	defer petcd.Close()

	sliServer := []string{
		"127.0.0.1:1",
		"127.0.0.1:2",
		"127.0.0.1:3",
	}
	for _, server := range sliServer {
		petcd.PutAlive("/server/"+server, "alive")
	}

	chaFunc := make(chan func(), 10)

	petcd.WatchPrefix("/server/", func(eventType mvccpb.Event_EventType, prefix string, kv *mvccpb.KeyValue) {
		chaFunc <- func() {
			plog.InfoF("server change,type:%v,prefix:%q,kv:%v\n", eventType, prefix, kv)
			s := petcd.TrimPrefix(string(kv.Key), prefix)
			switch eventType {
			case mvccpb.PUT:
				f := false
				for _, server := range sliServer {
					if s == server {
						f = true
						break
					}
				}
				if !f {
					sliServer = append(sliServer, s)
				}
			case mvccpb.DELETE:
				for i, server := range sliServer {
					if s == server {
						sliServer = append(sliServer[:i], sliServer[i+1:]...)
						break
					}
				}
			}
		}
	})

	ti := time.Tick(time.Second * 30)
	for {
		select {
		case f := <-chaFunc:
			f()
		case <-ti:
			plog.InfoLn("server count:", len(sliServer))
			for i, server := range sliServer {
				plog.InfoF("[%v]%v", i, server)
			}
		}
	}
}

func testPETCDFlushDB() {
	t := petcd.NewConfigEtcd()
	t.Endpoints = append(t.Endpoints, "http://127.0.0.1:2379")
	t.DialTimeout = 1
	t.OperationTimeout = 1
	t.LeaseTimeoutBeforeKeepAlive = 10
	petcd.Init(t)
	defer petcd.Close()

	plog.InfoLn(petcd.FlushDB())
}

func testPETCDCompact() {
	t := petcd.NewConfigEtcdInit()
	t.ConfigEtcd.Endpoints = append(t.ConfigEtcd.Endpoints, "http://localhost:60001", "http://localhost:60002", "http://localhost:60003")
	t.ConfigEtcd.DialTimeout = 1
	t.ConfigEtcd.OperationTimeout = 1
	t.ConfigEtcd.LeaseTimeoutBeforeKeepAlive = 10

	t.FlushDB = true
	t.EtcdKVList = append(t.EtcdKVList, petcd.NewEtcdKV("/server/127.0.0.1:60001", "alive"))
	t.EtcdKVList = append(t.EtcdKVList, petcd.NewEtcdKV("/server/127.0.0.1:60002", "alive"))
	t.EtcdKVList = append(t.EtcdKVList, petcd.NewEtcdKV("/server/127.0.0.1:60003", "alive"))

	petcd.ProcessEtcdInit(t)
}
