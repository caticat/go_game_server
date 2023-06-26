package main

import (
	"github.com/caticat/go_game_server/example/proto"
	"github.com/caticat/go_game_server/example/server/conf"
	"github.com/caticat/go_game_server/pnet"
)

var (
	unixTimeLast       int64 = 0
	unixTimeLastMinute int64 = 0
)

func runTimer(unixTimeNowMill int64) {
	unixTimeNow := unixTimeNowMill / conf.TimePrecision
	if unixTimeLast == 0 {
		unixTimeLast = unixTimeNowMill / conf.TimePrecision
		unixTimeLastMinute = unixTimeLast + conf.TimeMinuteSecond
	}

	if unixTimeNow != unixTimeLast {
		unixTimeLast = unixTimeNow
		runSecond(unixTimeNow)
	}
	if unixTimeNow >= unixTimeLastMinute {
		unixTimeLastMinute = unixTimeNow + conf.TimeMinuteSecond
		runMinute(unixTimeNow)
	}
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
