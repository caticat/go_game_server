package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caticat/go_game_server/petcd"
	"github.com/caticat/go_game_server/plog"
	"github.com/caticat/go_game_server/pnet"
	"github.com/caticat/go_game_server/pnet/conf"
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

type ConfServerEtcd struct {
	Server *conf.ConfServer  `yaml:"server" json:"server"`
	Etcd   *petcd.ConfigEtcd `yaml:"etcd"`
	Status *ConfServerStatus
}

type ConfServerStatus struct {
	IP string `json:"ip"`
}

func NewConfServerEtcd() *ConfServerEtcd {
	t := &ConfServerEtcd{
		Server: conf.NewConfServer(),
		Etcd:   petcd.NewConfigEtcd(),
		Status: NewConfServerStatus(),
	}
	return t
}

func NewConfServerStatus() *ConfServerStatus {
	return &ConfServerStatus{}
}

func (t *ConfServerStatus) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		plog.ErrorLn(err)
		return ""
	}

	return string(b)
}

func (t *ConfServerEtcd) Init() {
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

	if t.GetConnectionType() <= 0 {
		plog.FatalLn("config connection type <= 0")
	}

	t.updateIP()

	t.initConfFromEtcd()
}

func (t *ConfServerEtcd) GetServer() *conf.ConfServer { return t.Server }
func (t *ConfServerEtcd) GetID() int64                { return t.GetServer().GetID() }
func (t *ConfServerEtcd) GetConnectionType() int      { return t.GetServer().GetConnectionType() }
func (t *ConfServerEtcd) GetPort() int                { return t.GetServer().GetPort() }
func (t *ConfServerEtcd) GetPortIn() int              { return t.GetServer().GetPortIn() }
func (t *ConfServerEtcd) GetRemoteServers() []*conf.ConfServerRemote {
	return t.GetServer().GetRemoteServers()
}
func (t *ConfServerEtcd) GetLog() *plog.ConfLog        { return t.GetServer().GetLog() }
func (t *ConfServerEtcd) GetEtcd() *petcd.ConfigEtcd   { return t.Etcd }
func (t *ConfServerEtcd) GetStatus() *ConfServerStatus { return t.Status }

func (t *ConfServerEtcd) parseArgs() {
	flag.StringVar(&FileConfig, "c", "server.yaml", "-c=server.yaml")
	flag.StringVar(&FileConfig, "conf", "server.yaml", "-conf=server.yaml")
	flag.StringVar(&FileConfig, "config", "server.yaml", "-config=server.yaml")

	flag.Parse()
}

func (t *ConfServerEtcd) updateIP() {
	ips := pnet.GetIPs()

	// plog.InfoLn("ip长度:", len(ips))
	// for i, ip := range ips {
	// 	plog.InfoLn(i, "->", ip)
	// }

	if len(ips) > 0 {
		t.Status.IP = ips[0] // 本地测试,只取第一个就可以了
	}
}

func (t *ConfServerEtcd) initConfFromEtcd() {
	c := t.GetEtcd()
	if c == nil {
		plog.FatalLn("config etcd is nil")
	}
	petcd.Init(c)

	if err := petcd.GetServerConfig(t.GetServer()); err != nil {
		plog.FatalLn(err)
	}

	t.RegistService()
}

func (t *ConfServerEtcd) RegistService() {
	petcd.RegistService(t.GetServer(), t.GetStatus().String())
}
