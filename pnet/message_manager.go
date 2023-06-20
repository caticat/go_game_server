package pnet

import (
	"github.com/caticat/go_game_server/plog"
)

type messageHandler_t func(*PRecvData) bool

type PMessageManager struct {
	m_mapMsgIDHandler map[int32]messageHandler_t
}

func (t PMessageManager) New() *PMessageManager {
	t.m_mapMsgIDHandler = make(map[int32]messageHandler_t)
	return &t
}

func (t *PMessageManager) Regist(msgID int32, fun messageHandler_t) {
	t.m_mapMsgIDHandler[msgID] = fun
}

func (t *PMessageManager) Trigger(r *PRecvData) bool {
	if r == nil {
		plog.ErrorLn("r == nil")
		return false
	}
	m := r.GetMessage()
	if m == nil {
		plog.ErrorLn("m == nil")
		return false
	}

	fun, ok := t.m_mapMsgIDHandler[m.GetMsgID()]
	if !ok {
		plog.ErrorLn("msgID:%v not found\n", m.GetMsgID())
		return false
	}

	return fun(r)
}
