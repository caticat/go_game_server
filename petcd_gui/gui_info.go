package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func initGUIInfo(w fyne.Window, c binding.String) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabel("Configuration"),
		canvas.NewLine(color.Black),
		widget.NewLabelWithData(c),
	)
}
