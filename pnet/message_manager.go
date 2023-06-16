package pnet

import "log"

type messageHandler_t struct {
	needSoasdf
	fun func(*PSocket, *PMessage) bool
}

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

func (t *PMessageManager) Trigger(s *PSocket, m *PMessage) bool {
	if s == nil {
		log.Printf("s == nil")
		return false
	}
	if m == nil {
		log.Printf("m == nil")
		return false
	}

	fun, ok := t.m_mapMsgIDHandler[m.GetMsgID()]
	if !ok {
		log.Printf("msgID:%v not found\n", m.GetMsgID())
		return false
	}

	return fun(s, m)
}
