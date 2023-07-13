package main

import (
	"encoding/json"
	"fmt"

	"fyne.io/fyne/v2"
)

func updateWindowTitle(w fyne.Window) {
	connTxt := "Connecting"
	if !getConnected() {
		connTxt = "Disconnected"
	}
	w.SetTitle(fmt.Sprintf("%s[%s:%s]", WINDOW_TITLE, connTxt, getConf().EndPointSelect))
}

func toJsonIndent[T any](t T) (string, error) {
	b, e := json.MarshalIndent(t, "", "\t")
	return string(b), e
}
