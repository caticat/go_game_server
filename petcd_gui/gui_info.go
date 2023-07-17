package main

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func initGUIInfo(w fyne.Window, c binding.String) fyne.CanvasObject {
	urlRespitory, _ := url.Parse("https://github.com/caticat/go_game_server")
	urlLink := widget.NewHyperlink("github.com/caticat/go_game_server", urlRespitory)

	return container.NewScroll(
		container.NewVBox(
			widget.NewLabel("Information"),
			widget.NewForm(
				widget.NewFormItem("author", widget.NewLabel(APP_AUTHOR)),
				widget.NewFormItem("repository", urlLink),
				widget.NewFormItem("version", widget.NewLabel(APP_VER)),
			),
			layout.NewSpacer(),
			widget.NewLabel("Configuration"),
			widget.NewLabelWithData(c),
		),
	)
}
