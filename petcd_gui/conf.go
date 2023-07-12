package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caticat/go_game_server/petcd"
	"gopkg.in/yaml.v3"
)

var FileConfig = ""

type ConfigEtcdGUI struct {
	Etcd *petcd.ConfigEtcd `yaml:"etcd"`
}

func NewConfigEtcdGUI() *ConfigEtcdGUI {
	return &ConfigEtcdGUI{}
}

func (t *ConfigEtcdGUI) Init() error {
	flag.StringVar(&FileConfig, "c", "petcd_gui.yml", "-c petcd_gui.yml")
	flag.Parse()

	f, err := os.ReadFile(FileConfig)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(f, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *ConfigEtcdGUI) ToJson() (string, error) {
	s, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		return "", err
	}

	return string(s), err
}
