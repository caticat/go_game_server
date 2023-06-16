package main

import "github.com/caticat/go_game_server/pnet"

type MessageManager struct {
	*pnet.PMessageManager
}

func (t MessageManager) New() *MessageManager {
	t.PMessageManager = pnet.PMessageManager{}.New()
	return &t
}

func (t *MessageManager) Init() {
	// TODO: 消息待注册
}
