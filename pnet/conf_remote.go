package pnet

import "fmt"

type ConfRemoteServer struct {
	ServerID       int    `yaml:"server_id" json:"server_id"`
	ConnectionType int    `yaml:"connection_type" json:"connection_type"`
	IP             string `yaml:"ip" json:"ip"`
	Port           int    `yaml:"port" json:"port_in"`
}

func NewConfRemoteServer() *ConfRemoteServer       { return &ConfRemoteServer{} }
func (t *ConfRemoteServer) GetServerID() int       { return t.ServerID }
func (t *ConfRemoteServer) GetConnectionType() int { return t.ConnectionType }
func (t *ConfRemoteServer) GetIP() string          { return t.IP }
func (t *ConfRemoteServer) GetPort() int           { return t.Port }
func (t *ConfRemoteServer) String() string         { return fmt.Sprintf("%s:%d", t.GetIP(), t.GetPort()) }
