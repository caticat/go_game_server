package ptime

import "time"

type PTicker struct {
	m_interval     int64
	m_fun          func(int64)
	m_unixTimeTick int64
}

func NewPTicker(interval int64, fun func(unixTimeNow int64)) *PTicker {
	t := &PTicker{
		m_interval: interval,
		m_fun:      fun,
	}
	t.updateUnixTimeTick(time.Now().Unix())
	return t
}

func (t *PTicker) TryRun(unixTimeNow int64) bool {
	if unixTimeNow < t.m_unixTimeTick {
		return false
	}

	t.m_fun(unixTimeNow)

	t.updateUnixTimeTick(unixTimeNow)

	return true
}

func (t *PTicker) updateUnixTimeTick(unixTimeNow int64) {
	t.m_unixTimeTick = unixTimeNow + t.m_interval
}
