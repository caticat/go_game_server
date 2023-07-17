package petcd

type EtcdKV struct {
	Key   string `yaml:"key" json:"key"`
	Value string `yaml:"value" json:"value"`
}

func NewEtcdKV(key, value string) *EtcdKV {
	return &EtcdKV{key, value}
}

type ConfigEtcdInit struct {
	*ConfigEtcd `yaml:"etcd" json:"etcd"` // 连接配置
	FlushDB     bool                      `yaml:"flush_db" json:"flush_db"`             // 是否清档
	EtcdKVList  []*EtcdKV                 `yaml:"list_key_value" json:"list_key_value"` // 初始化键值对
}

func NewConfigEtcdInit() *ConfigEtcdInit {
	return &ConfigEtcdInit{ConfigEtcd: NewConfigEtcd()}
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

	return ProcessEtcdInitWithoutConn(cfg)
}

// 跑初始化配置
func ProcessEtcdInitWithoutConn(cfg *ConfigEtcdInit) error {
	if cfg == nil {
		return ErrorNilConfig
	}

	// 清档
	if cfg.FlushDB {
		FlushDB()
	}

	// 初始化数据
	for _, p := range cfg.EtcdKVList {
		if err := Put(p.Key, p.Value); err != nil {
			return err
		}
	}

	return nil
}

func (t *ConfigEtcdInit) SetBase(b *ConfigEtcdInitBase) {
	t.FlushDB = b.FlushDB
	t.EtcdKVList = b.EtcdKVList
}

func (t *ConfigEtcdInit) SetConfigEtcd(c *ConfigEtcd) {
	t.ConfigEtcd = c
}
