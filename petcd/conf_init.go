package petcd

type ConfigEtcdInit struct {
	*ConfigEtcd `yaml:"etcd" json:"etcd"` // 连接配置
	FlushDB     bool                      `yaml:"flush_db" json:"flush_db"`             // 是否清档
	EtcdKVList  []*EtcdKV                 `yaml:"list_key_value" json:"list_key_value"` // 初始化键值对
}

type EtcdKV struct {
	Key   string `yaml:"key" json:"key"`
	Value string `yaml:"value" json:"value"`
}

func NewConfigEtcdInit() *ConfigEtcdInit {
	return &ConfigEtcdInit{ConfigEtcd: NewConfigEtcd()}
}

func NewEtcdKV(key, value string) *EtcdKV {
	return &EtcdKV{key, value}
}

// 跑初始化配置
func ProcessEtcdInit(cfg *ConfigEtcdInit) error {
	if cfg == nil {
		return ErrorNilConfig
	}
	if err := Init(cfg.ConfigEtcd); err != nil {
		return err
	}
	defer Close()

	// 清档
	if cfg.FlushDB {
		FlushDB()
	}

	// 初始化数据
	for _, p := range cfg.EtcdKVList {
		Put(p.Key, p.Value)
	}

	return nil
}
