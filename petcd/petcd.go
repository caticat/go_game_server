package petcd

import (
	"context"
	"strings"

	"go.etcd.io/etcd/clientv3"
)

func Init(cfg *ConfigEtcd) error {
	var err error = nil
	getInit().Do(
		func() {
			// 初始化连接,配置
			if cfg == nil {
				err = ErrorNilConfig
				return
			}
			setConfig(cfg)
			cli, err := clientv3.New(cfg.ToConfig())
			if err != nil {
				return
			}
			setClient(cli)

			// 初始化lease
			ctx, cancel := context.WithTimeout(context.Background(), getConfigOperationTimeout())
			resp, err := cli.Lease.Grant(ctx, cfg.LeaseTimeoutBeforeKeepAlive)
			cancel()
			if err != nil {
				return
			}
			leaseID := resp.ID
			setLeaseID(leaseID)

			// 启动lease的KeepAlive
			ctx, cancel = context.WithCancel(context.Background())
			cha, err := cli.KeepAlive(ctx, leaseID)
			if err != nil {
				cancel()
				return
			}
			setLeaseCancel(cancel)
			go run(cha)
		})

	return err
}

func Close() {
	cli := getClient()
	if cli == nil {
		return
	}

	leaseCancel := getLeaseCancel()
	if leaseCancel != nil {
		leaseCancel()
		setLeaseCancel(nil)
	}

	sliCancelWatch := getSliCancelWatch()
	for _, watchCancel := range sliCancelWatch {
		watchCancel()
	}

	cli.Close()
	setClient(nil)
}

func Put(key, val string, opts ...clientv3.OpOption) error {
	cli := getClient()
	if cli == nil {
		return ErrorNilClient
	}

	ctx, cancel := context.WithTimeout(context.Background(), getConfigOperationTimeout())
	defer cancel()
	_, err := cli.Put(ctx, key, val, opts...)

	return err
}

// 设置键值并绑定生命周期
func PutAlive(key, val string) error {
	return Put(key, val, clientv3.WithLease(getLeaseID()))
}

// 设置键值并保持原有的生命周期
func PutKeepLease(key, val string) error {
	if _, err := GetString(key); err != nil { // 没有的话就直接创建
		return Put(key, val)
	} else {
		return Put(key, val, clientv3.WithIgnoreLease())
	}
}

func Get(key string, mapResult map[string]string, opts ...clientv3.OpOption) error {
	cli := getClient()
	if cli == nil {
		return ErrorNilClient
	}

	ctx, cancel := context.WithTimeout(context.Background(), getConfigOperationTimeout())
	defer cancel()
	resp, err := cli.Get(ctx, key, opts...)
	if err != nil {
		return err
	}

	for _, p := range resp.Kvs {
		mapResult[string(p.Key)] = string(p.Value)
	}

	if len(mapResult) == 0 {
		return ErrorKeyNotFound
	} else {
		return nil
	}
}

func GetPrefix(key string, mapResult map[string]string) error {
	return Get(key, mapResult, clientv3.WithPrefix())
}

// 简单的KV获取封装
func GetString(key string) (string, error) {
	mapResult := make(map[string]string)
	if err := Get(key, mapResult); err != nil {
		return "", err
	}
	cntResult := len(mapResult)
	if cntResult == 0 {
		return "", nil
	} else {
		for _, v := range mapResult {
			return v, nil
		}
	}
	return "", ErrorInvalidOperation
}

func Del(key string, opts ...clientv3.OpOption) error {
	cli := getClient()
	if cli == nil {
		return ErrorNilClient
	}

	ctx, cancel := context.WithTimeout(context.Background(), getConfigOperationTimeout())
	defer cancel()
	_, err := cli.Delete(ctx, key, opts...)
	if err != nil {
		return err
	}

	return nil
}

func Watch(key string, fun funWatchCallback_t, opts ...clientv3.OpOption) error {
	cli := getClient()
	if cli == nil {
		return ErrorNilClient
	}

	ctx, cancel := context.WithCancel(context.Background())
	appendSliCancelWatch(cancel)
	cha := cli.Watch(ctx, key, opts...)
	go runWatch(cha, key, fun)

	return nil
}

func WatchPrefix(prefix string, fun funWatchCallback_t) error {
	return Watch(prefix, fun, clientv3.WithPrefix())
}

// 清空数据库
func FlushDB() error {
	if err := Del("", clientv3.WithPrefix()); err != nil {
		return err
	}

	if err := Compact(); err != nil {
		return err
	}

	if err := Defrag(); err != nil {
		return err
	}

	return nil
}

// 清理历史版本数据记录
func Compact() error {
	cli := getClient()
	if cli == nil {
		return ErrorNilClient
	}

	revisionID, err := getRevision()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), getConfigOperationTimeout())
	defer cancel()
	_, err = cli.Compact(ctx, revisionID)

	return err
}

// 恢复硬盘空间
func Defrag() error {
	cli := getClient()
	if cli == nil {
		return ErrorNilClient
	}
	cfg := getConfig()
	if cfg == nil {
		return ErrorNilConfig
	}

	for _, ep := range cfg.Endpoints {
		ctx, cancel := context.WithTimeout(context.Background(), getConfigOperationTimeout()) // MARK: 这里的超时时间可能过短
		defer cancel()
		cli.Defragment(ctx, ep)
	}

	return nil
}

func TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}
