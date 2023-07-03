package main

import (
	"flag"
	"os"

	"github.com/caticat/go_game_server/plog"
	"github.com/caticat/go_game_server/pnet"
	"gopkg.in/yaml.v3"
)

const (
	ChaRecvLen             = 100
	ChaMainLoopFun         = 10
	TimePrecision    int64 = 1000
	TimeMinuteSecond int64 = 60
)

var (
	FileConfig string = "server.yaml"
)

type ConfServer struct {
	ID             int64                    `yaml:"id" json:"id"`
	ConnectionType int                      `yaml:"connection_type" json:"connection_type"`
	Port           int                      `yaml:"port" json:"port"`
	PortIn         int                      `yaml:"port_in" json:"port_in"`
	Log            *plog.ConfLog            `yaml:"log" json:"log"`
	RemoteServers  []*pnet.ConfRemoteServer `yaml:"remote_server" json:"remote_server"`
}

func NewConfServer() *ConfServer {
	t := &ConfServer{
		Log: plog.NewConfLog(),
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

	if t.GetID() <= 0 {
		plog.FatalLn("config id <= 0")
	}
}

func (t *ConfServer) GetID() int64                               { return t.ID }
func (t *ConfServer) GetConnectionType() int                     { return t.ConnectionType }
func (t *ConfServer) GetPort() int                               { return t.Port }
func (t *ConfServer) GetPortIn() int                             { return t.PortIn }
func (t *ConfServer) GetLog() *plog.ConfLog                      { return t.Log }
func (t *ConfServer) GetRemoteServers() []*pnet.ConfRemoteServer { return t.RemoteServers }

func (t *ConfServer) parseArgs() {
	flag.StringVar(&FileConfig, "c", "server.yaml", "-c=server.yaml")
	flag.StringVar(&FileConfig, "conf", "server.yaml", "-conf=server.yaml")
	flag.StringVar(&FileConfig, "config", "server.yaml", "-config=server.yaml")

	flag.Parse()
}
