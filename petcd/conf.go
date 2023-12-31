package petcd

import (
	"time"

	"go.etcd.io/etcd/clientv3"
)

type ConfigEtcd struct {
	Endpoints                   []string        `yaml:"endpoints" json:"endpoints"`
	DialTimeout                 int64           `yaml:"dial-timeout" json:"dial-timeout"`                                       // 连接超时时间 秒
	OperationTimeout            int64           `yaml:"operation-timeout" json:"operation-timeout"`                             // 操作超时时间 秒
	LeaseTimeoutBeforeKeepAlive int64           `yaml:"lease-timeout-before-keep-alive" json:"lease-timeout-before-keep-alive"` // 租约续期前的过期时间(连接断开后多长时间ETCD数据会消失)
	EnableReInit                bool            `yaml:"enable-reinit" json:"enable-reinit"`
	Auth                        *ConfigEtcdAuth `yaml:"auth" json:"auth"`
}

func NewConfigEtcd() *ConfigEtcd {
	return &ConfigEtcd{
		Auth: NewConfigEtcdAuth(),
	}
}

func (t *ConfigEtcd) SetBase(b *ConfigEtcdBase) {
	t.DialTimeout = b.DialTimeout
	t.OperationTimeout = b.OperationTimeout
	t.LeaseTimeoutBeforeKeepAlive = b.LeaseTimeoutBeforeKeepAlive
	t.EnableReInit = b.EnableReInit
}

func (t *ConfigEtcd) SetConn(c *ConfigEtcdConn) {
	t.Endpoints = c.Endpoints
	t.Auth = c.Auth
}

func (t *ConfigEtcd) ToConfig() clientv3.Config {
	return clientv3.Config{
		Endpoints:   t.Endpoints,
		DialTimeout: time.Second * time.Duration(t.DialTimeout),
		Username:    t.Auth.Username,
		Password:    t.Auth.Password,
	}
}
