package petcd

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/caticat/go_game_server/pnet"
)

// 拼接Key
func FormatKey(args ...any) string {
	sliKey := make([]string, 0, len(args))
	for _, k := range args {
		sliKey = append(sliKey, fmt.Sprint(k))
	}

	return path.Join(sliKey...)
}

// type key_t interface {
// 	int | int64 | string
// }

// func FormatKeySameType[T key_t](args ...T) string {
// 	sliKey := make([]string, 0, len(args))
// 	for _, k := range args {
// 		sliKey = append(sliKey, fmt.Sprint(k))
// 	}

// 	return path.Join(sliKey...)
// }

// 获取服务器配置
func GetServerConfig(c *pnet.ConfServer) error {
	if c == nil {
		return ErrorNilConfig
	}

	k := FormatKey(Config, c.GetConnectionType(), c.GetID())
	v, err := GetString(k)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(v), c)
	if err != nil {
		return err
	}

	ss := c.GetRemoteServers()
	for _, s := range ss {
		k = FormatKey(Config, s.GetConnectionType(), s.GetServerID())
		v, err = GetString(k)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(v), s)
		if err != nil {
			return err
		}
	}

	return nil
}

// 注册服务
func RegistService(c *pnet.ConfServer, value string) error {
	if c == nil {
		return ErrorNilConfig
	}

	k := FormatKey(Service, c.GetConnectionType(), c.GetID())

	return PutAlive(k, value)
}
