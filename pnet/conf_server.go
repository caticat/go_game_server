package pnet

import "github.com/caticat/go_game_server/plog"

type ConfServer struct {
	ID             int64               `yaml:"id" json:"id"`
	ConnectionType int                 `yaml:"connection_type" json:"connection_type"`
	Port           int                 `yaml:"port" json:"port"`
	PortIn         int                 `yaml:"port_in" json:"port_in"`
	RemoteServers  []*ConfRemoteServer `yaml:"remote_server" json:"remote_server"`
	Log            *plog.ConfLog       `yaml:"log" json:"log"`
}

func NewConfServer() *ConfServer {
	return &ConfServer{
		Log: plog.NewConfLog(),
	}
}

func (t *ConfServer) GetID() int64                          { return t.ID }
func (t *ConfServer) GetConnectionType() int                { return t.ConnectionType }
func (t *ConfServer) GetPort() int                          { return t.Port }
func (t *ConfServer) GetPortIn() int                        { return t.PortIn }
func (t *ConfServer) GetRemoteServers() []*ConfRemoteServer { return t.RemoteServers }
func (t *ConfServer) GetLog() *plog.ConfLog                 { return t.Log }
