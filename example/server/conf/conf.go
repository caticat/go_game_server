package conf

import (
	"os"

	"github.com/caticat/go_game_server/plog"
	"gopkg.in/yaml.v2"
)

const (
	ChaRecvLen = 100
	FileConfig = "server.yaml"
)

type ConfServer struct {
	Port int      `yaml:"port"`
	Log  *ConfLog `yaml:"log"`
}

func (t ConfServer) New() *ConfServer {
	t.Log = NewConfLog()
	return &t
}

func (t *ConfServer) Init() {
	f, err := os.ReadFile(FileConfig)
	if err != nil {
		plog.FatalLn("ioutil.ReadFile failed,error:", err)
	}

	err = yaml.Unmarshal(f, t)
	if err != nil {
		plog.FatalLn("yaml.Unmarshal failed,error:", err)
	}
}

func (t *ConfServer) GetPort() int     { return t.Port }
func (t *ConfServer) GetLog() *ConfLog { return t.Log }
