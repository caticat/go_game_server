package petcd

// 配置文件拆分,序列化用

type ConfigEtcdBase struct {
	DialTimeout                 int64 `yaml:"dial-timeout" json:"dial-timeout"`                                       // 连接超时时间 秒
	OperationTimeout            int64 `yaml:"operation-timeout" json:"operation-timeout"`                             // 操作超时时间 秒
	LeaseTimeoutBeforeKeepAlive int64 `yaml:"lease-timeout-before-keep-alive" json:"lease-timeout-before-keep-alive"` // 租约续期前的过期时间(连接断开后多长时间ETCD数据会消失)
	EnableReInit                bool  `yaml:"enable-reinit" json:"enable-reinit"`
}

func NewConfigEtcdBase() *ConfigEtcdBase {
	return &ConfigEtcdBase{}
}

func NewConfigEtcdBaseByConfig(c *ConfigEtcd) *ConfigEtcdBase {
	return &ConfigEtcdBase{
		DialTimeout:                 c.DialTimeout,
		OperationTimeout:            c.OperationTimeout,
		LeaseTimeoutBeforeKeepAlive: c.LeaseTimeoutBeforeKeepAlive,
		EnableReInit:                c.EnableReInit,
	}
}
