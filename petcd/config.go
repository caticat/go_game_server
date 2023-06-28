package petcd

import (
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
)

type funWatchCallback_t func(eventType mvccpb.Event_EventType, prefix string, kv *mvccpb.KeyValue)

type ConfigEtcd struct {
	Endpoints                   []string      `yaml:"endpoints"`
	DialTimeout                 time.Duration `yaml:"dial-timeout"`
	OperationTimeout            time.Duration `yaml:"operation-timeout"`               // 操作超时时间
	LeaseTimeoutBeforeKeepAlive int64         `yaml:"lease-timeout-before-keep-alive"` // 租约续期前的过期时间(连接断开后多长时间ETCD数据会消失)
}

func NewConfigEtcd() *ConfigEtcd {
	return &ConfigEtcd{}
}

func (t *ConfigEtcd) ToConfig() clientv3.Config {
	return clientv3.Config{
		Endpoints:   t.Endpoints,
		DialTimeout: t.DialTimeout,
	}
}
