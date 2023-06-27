package ptime

import "time"

// 获取下一个小时的整点时间戳
func GetNextHourTime(t time.Time) int64 {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.Local).Add(time.Hour).Unix()
}
