package main

import (
	"github.com/caticat/go_game_server/example/proto"
	"github.com/caticat/go_game_server/pnet"
	"github.com/caticat/go_game_server/ptime"
)

var (
	g_ticker1s *ptime.PTicker
	g_ticker1m *ptime.PTicker
)

func initTimer() {
	g_ticker1s = ptime.NewPTicker(1, runSecond)
	g_ticker1m = ptime.NewPTicker(TimeMinuteSecond, runMinute)
}

func runTimer(unixTimeNowMill int64) {
	unixTimeNow := unixTimeNowMill / TimePrecision

	g_ticker1s.TryRun(unixTimeNow)
	g_ticker1m.TryRun(unixTimeNow)
}

func runSecond(unixTimeNow int64) {
	// 每秒调用
}

func runMinute(unixTimeNow int64) {
	// 每分调用
	ss := getSocketManager().GetServerAll()
	if len(ss) > 0 {
		m := pnet.NewPMessage(int32(proto.MsgID_TickNtfID), &proto.TickNtf{})
		for _, s := range ss {
			s.Send(m)
		}
	}
}
