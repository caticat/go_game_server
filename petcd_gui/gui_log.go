package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/caticat/go_game_server/plog"
)

func initGUILog(w fyne.Window) fyne.CanvasObject {
	s := widget.NewSelect([]string{"debug", "info", "warn"}, func(s string) {
		logLevel := plog.ToLogLevel(s)
		getConf().SetLogLevel(logLevel)
		plog.InfoLn("change log level to:", s)
	})
	s.SetSelected(plog.ToLogLevelName(plog.ELogLevel(getConf().LogLevel)))
	c := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		getLogData().Set("")
	})
	hs := container.NewAdaptiveGrid(2, s, container.NewHBox(c))
	h := widget.NewForm(widget.NewFormItem("LogLevel", hs))

	b := container.NewScroll(widget.NewLabelWithData(getLogData()))

	return container.NewBorder(h, nil, nil, nil, b)
}
