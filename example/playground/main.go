package main

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/caticat/go_game_server/petcd"
	"github.com/caticat/go_game_server/petcd/pdata"
	"github.com/caticat/go_game_server/phelp"
	"github.com/caticat/go_game_server/phelp/ppath"
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
	// testPETCDPrefix()
	// testPETCDAuth()
	// testCP()
	// testIsSubDir()
	testMv()
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
		plog.InfoF("[%d]%q->%q/n", i, string(p.Key), string(p.Value))
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
	t.Endpoints = append(t.Endpoints, "http://localhost:60001", "http://localhost:60002", "http://localhost:60003")
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

	petcd.PutAlive("/server/127.0.0.1:1", "alive")
	petcd.PutAlive("/server/127.0.0.1:2", "alive")
	petcd.PutAlive("/server/127.0.0.1:3", "alive")

	petcd.Put("/server/127.0.0.1:1", "abcdefg12345")

	mapServer := make(map[string]string)
	petcd.GetPrefix("/server/", mapServer)
	for k, v := range mapServer {
		plog.InfoF("%v[%v]->%v/n", k, petcd.TrimPrefix(k, "/server/"), v)
	}

	time.Sleep(time.Second * 10)
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
			plog.InfoF("server change,type:%v,prefix:%q,kv:%v/n", eventType, prefix, kv)
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

func testPETCDPrefix() {
	t := petcd.NewConfigEtcd()
	t.Endpoints = []string{"http://localhost:60001", "http://localhost:60002", "http://localhost:60003"}
	t.DialTimeout = 1
	t.OperationTimeout = 1
	t.LeaseTimeoutBeforeKeepAlive = 10
	petcd.Init(t)
	defer petcd.Close()

	plog.InfoLn("")
	mapResult := make(map[string]string)
	petcd.GetPrefix(pdata.PDATA_PREFIX, mapResult)

	root := pdata.NewPEtcdRoot()
	root.SetAll(mapResult)

	plog.DebugLn(mapResult)
}

func testPETCDAuth() {
	t := petcd.NewConfigEtcd()
	t.Endpoints = []string{"http://localhost:60001", "http://localhost:60002", "http://localhost:60003"}
	t.DialTimeout = 1
	t.OperationTimeout = 1
	t.LeaseTimeoutBeforeKeepAlive = 10
	t.Auth.Username = "pan"
	t.Auth.Password = "pan"
	petcd.Init(t)
	defer petcd.Close()

	mapResult := make(map[string]string)
	if err := petcd.GetPrefix(pdata.PDATA_PREFIX, mapResult); err != nil {
		plog.Error(err)
	}

	plog.DebugLn(mapResult)
}

func testCP() {
	from := "E:/pan/study_notes/go/git/tmp/a"
	to := "E:/pan/study_notes/go/git"
	if _, err := phelp.Cp(from, to, phelp.PBinFlag_Recursive|phelp.PBinFlag_Force); err != nil {
		plog.ErrorLn(err)
		return
	}
	plog.InfoLn("copy done")
}

func testIsSubDir() {
	d, err := os.Getwd()
	if err != nil {
		plog.InfoLn("err:", err)
		return
	}
	plog.InfoLn("base:", d)
	a := "./data"
	a, err = filepath.Abs(path.Join(d, a))
	if err != nil {
		plog.InfoLn("err:", err)
		return
	}

	f := func(s string) string {
		if filepath.IsAbs(s) {
			return s
		}
		s, err := filepath.Abs(path.Join(a, s))
		if err != nil {
			plog.InfoLn("err:", err)
			return s
		}
		return s
	}

	plog.InfoLn(ppath.IsSubDir(a, f("b/c")))
	plog.InfoLn(ppath.IsSubDir(a, f("../b/c")))
	plog.InfoLn(ppath.IsSubDir(a, f("../data/b/c")))
	plog.InfoLn(ppath.IsSubDir(a, f("/data/b/c")))
	plog.InfoLn(ppath.IsSubDir(a, f("e:/data/b/c")))

	plog.InfoLn(filepath.Abs("a/b/c"))
	plog.InfoLn(filepath.Abs("e:/a/b/c"))
}

func mv(from, to string, flags phelp.PBinFlags) error {
	if _, err := phelp.Cp(from, to, flags); err != nil {
		return err
	}
	return phelp.Rm(from)
}

func testMv() {
	from := "C:\\Users\\pan\\AppData\\Roaming\\fyne\\github.com.caticat.go_git_notebook_gui\\local.repo\\myapp.png"
	to := "C:\\Users\\pan\\AppData\\Roaming\\fyne\\github.com.caticat.go_git_notebook_gui\\local.repo\\a"
	from = phelp.Format(from)
	to = phelp.Format(to)
	if err := mv(from, to, phelp.PBinFlag_Recursive|phelp.PBinFlag_Force); err != nil {
		plog.ErrorLn(err)
		return
	}
	plog.InfoLn("mv done")
}
