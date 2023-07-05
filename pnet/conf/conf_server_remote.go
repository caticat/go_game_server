package conf

import "fmt"

type ConfServerRemote struct {
	ServerID       int    `yaml:"server_id" json:"server_id"`
	ConnectionType int    `yaml:"connection_type" json:"connection_type"`
	IP             string `yaml:"ip" json:"ip"`
	Port           int    `yaml:"port" json:"port_in"` // 配置文件直接配置对应服务器的内部监听端口,在etcd的配置中读取需要读取对方配置的内部监听端口才行,所以两个值是不一样的
}

func NewConfServerRemote() *ConfServerRemote       { return &ConfServerRemote{} }
func (t *ConfServerRemote) GetServerID() int       { return t.ServerID }
func (t *ConfServerRemote) GetConnectionType() int { return t.ConnectionType }
func (t *ConfServerRemote) GetIP() string          { return t.IP }
func (t *ConfServerRemote) GetPort() int           { return t.Port }
func (t *ConfServerRemote) String() string         { return fmt.Sprintf("%s:%d", t.GetIP(), t.GetPort()) }
