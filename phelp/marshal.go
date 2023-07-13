package phelp

import (
	"encoding/json"

	"github.com/caticat/go_game_server/plog"
)

func ToJson[T any](t T) string {
	b, e := json.Marshal(t)
	if e != nil {
		plog.ErrorLn(e)
		return ""
	}
	return string(b)
}

func ToJsonIndent[T any](t T) string {
	b, e := json.MarshalIndent(t, "", "\t")
	if e != nil {
		plog.ErrorLn(e)
		return ""
	}
	return string(b)
}
