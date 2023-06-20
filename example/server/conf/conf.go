package conf

import (
	"flag"
	"os"

	"github.com/caticat/go_game_server/plog"
	"gopkg.in/yaml.v2"
)

const (
	ChaRecvLen = 100
	// FileConfig = "server.yaml"
)

var (
	FileConfig string
)

type ConfServer struct {
	Port int      `yaml:"port"`
	Log  *ConfLog `yaml:"log"`
}

func NewConfServer() *ConfServer {
	t := &ConfServer{
		Log: NewConfLog(),
	}
	return t
}

func (t *ConfServer) Init() {
	t.parseArgs()

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

func (t *ConfServer) parseArgs() {
	flag.StringVar(&FileConfig, "c", "server.yaml", "-c=server.yaml")
	flag.StringVar(&FileConfig, "conf", "server.yaml", "-conf=server.yaml")
	flag.StringVar(&FileConfig, "config", "server.yaml", "-config=server.yaml")

	flag.Parse()
}
