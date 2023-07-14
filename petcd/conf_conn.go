package petcd

type ConfigEtcdConn struct {
	Endpoints []string        `json:"endpoints"`
	Auth      *ConfigEtcdAuth `json:"auth"`
}

func NewConfigEtccdConn() *ConfigEtcdConn {
	return &ConfigEtcdConn{
		Auth: NewConfigEtcdAuth(),
	}
}
