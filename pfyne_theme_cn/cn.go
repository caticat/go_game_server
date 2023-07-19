package pfyne_theme_cn

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	pfont "github.com/caticat/go_game_server/pfyne_theme_cn/font"
)

type ThemeCN struct{}

var _ fyne.Theme = (*ThemeCN)(nil)

func NewThemeCN() *ThemeCN {
	return &ThemeCN{}
}

func (t *ThemeCN) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}
func (t *ThemeCN) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *ThemeCN) Font(style fyne.TextStyle) fyne.Resource {
	//return theme.DefaultTheme().Font(style)
	return &fyne.StaticResource{
		StaticName:    "simhei.ttf",
		StaticContent: pfont.TTF_SIMHEI,
	}
}

func (t *ThemeCN) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
