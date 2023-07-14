package petcd

type ConfigEtcdAuth struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

func NewConfigEtcdAuth() *ConfigEtcdAuth {
	return &ConfigEtcdAuth{}
}
