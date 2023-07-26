package phelp

import (
	"encoding/json"
)

func ToJson[T any](t T) string {
	b, e := json.Marshal(t)
	if e != nil {
		return ""
	}
	return string(b)
}

func ToJsonIndent[T any](t T) string {
	b, e := json.MarshalIndent(t, "", "\t")
	if e != nil {
		return ""
	}
	return string(b)
}
