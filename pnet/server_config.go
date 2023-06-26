package pnet

import "fmt"

type ConfRemoteServer struct {
	ServerID       int    `yaml:"server_id"`
	ConnectionType int    `yaml:"connection_type"`
	IP             string `yaml:"ip"`
	Port           int    `yaml:"port"`
}

func NewConfRemoteServer() *ConfRemoteServer       { return &ConfRemoteServer{} }
func (t *ConfRemoteServer) GetServerID() int       { return t.ServerID }
func (t *ConfRemoteServer) GetConnectionType() int { return t.ConnectionType }
func (t *ConfRemoteServer) GetIP() string          { return t.IP }
func (t *ConfRemoteServer) GetPort() int           { return t.Port }
func (t *ConfRemoteServer) String() string         { return fmt.Sprintf("%s:%d", t.GetIP(), t.GetPort()) }
