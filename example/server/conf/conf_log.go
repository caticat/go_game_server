package conf

import "github.com/caticat/go_game_server/plog"

type ConfLog struct {
	Level plog.ELogLevel `yaml:"level"`
	File  string         `yaml:"file"`
}

func NewConfLog() *ConfLog {
	return &ConfLog{}
}

func (t *ConfLog) GetLevel() plog.ELogLevel { return t.Level }
func (t *ConfLog) GetFile() string          { return t.File }
