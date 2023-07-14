package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/caticat/go_game_server/plog"
)

func initGUILog(w fyne.Window) fyne.CanvasObject {
	// 上侧日志等级选择与功能按钮
	guiSelLogLevel := widget.NewSelect([]string{
		string(plog.SLogLevel_Debug),
		string(plog.SLogLevel_Info),
		string(plog.SLogLevel_Warn)}, func(s string) {
		logLevel := plog.ToLogLevel(s)
		getConf().SetLogLevel(logLevel)
		plog.InfoLn("change log level to:", s)
	})
	guiSelLogLevel.SetSelected(plog.ToLogLevelName(plog.ELogLevel(getConf().LogLevel)))
	guiButClearLog := widget.NewButtonWithIcon(STR_EMPTY, theme.DeleteIcon(), func() {
		getLogData().Set(STR_EMPTY)
	})
	guiConLogLevel := container.NewAdaptiveGrid(2, guiSelLogLevel, container.NewHBox(guiButClearLog))
	guiForLogLevel := widget.NewForm(widget.NewFormItem("LogLevel", guiConLogLevel))

	// 日志内容
	guiScrLog := container.NewScroll(widget.NewLabelWithData(getLogData()))

	// 界面组合
	return container.NewBorder(guiForLogLevel, nil, nil, nil, guiScrLog)
}
