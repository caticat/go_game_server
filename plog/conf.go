package plog

type ConfLog struct {
	Level ELogLevel `yaml:"level"`
	File  string    `yaml:"file"`
}

func NewConfLog() *ConfLog {
	return &ConfLog{}
}

func (t *ConfLog) GetLevel() ELogLevel { return t.Level }
func (t *ConfLog) GetFile() string     { return t.File }
