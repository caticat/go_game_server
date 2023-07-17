package petcd

type ConfigEtcdInitBase struct {
	FlushDB    bool      `yaml:"flush_db" json:"flush_db"`             // 是否清档
	EtcdKVList []*EtcdKV `yaml:"list_key_value" json:"list_key_value"` // 初始化键值对
}

func NewConfigEtcdInitBase() *ConfigEtcdInitBase {
	return &ConfigEtcdInitBase{}
}
