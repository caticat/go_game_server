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
			ctx, cancel := context.WithTimeout(context.Background(), cfg.OperationTimeout)
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
	cfg := getConfig()
	if cfg == nil {
		return ErrorNilConfig
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.OperationTimeout)
	defer cancel()
	_, err := cli.Put(ctx, key, val, opts...)

	return err
}

// 设置键值并绑定生命周期
func PutAlive(key, val string) error {
	return Put(key, val, clientv3.WithLease(getLeaseID()))
}

func Get(key string, mapResult map[string]string, opts ...clientv3.OpOption) error {
	cli := getClient()
	if cli == nil {
		return ErrorNilClient
	}
	cfg := getConfig()
	if cfg == nil {
		return ErrorNilConfig
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.OperationTimeout)
	defer cancel()
	resp, err := cli.Get(ctx, key, opts...)
	if err != nil {
		return err
	}

	for _, p := range resp.Kvs {
		mapResult[string(p.Key)] = string(p.Value)
	}

	return nil
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
	cfg := getConfig()
	if cfg == nil {
		return ErrorNilConfig
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.OperationTimeout)
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

func TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}
