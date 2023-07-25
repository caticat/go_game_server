package pfyne_theme_cn

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	pfont "github.com/caticat/go_game_server/pfyne_theme_cn/font"
)

type ThemeCN struct {
	bold, italic, symbol fyne.Resource
}

var _ fyne.Theme = (*ThemeCN)(nil)

func NewThemeCN() *ThemeCN {
	return &ThemeCN{
		bold: &fyne.StaticResource{
			StaticName:    "OPPOSans-H-2.ttf",
			StaticContent: pfont.TTF_OPPOSans_H_2,
		},
		italic: &fyne.StaticResource{
			StaticName:    "OPPOSans-L-2.ttf",
			StaticContent: pfont.TTF_OPPOSans_L_2,
		},
		symbol: &fyne.StaticResource{
			StaticName:    "OPPOSans-M-2.ttf",
			StaticContent: pfont.TTF_OPPOSans_M_2,
		},
	}
}

func (t *ThemeCN) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}
func (t *ThemeCN) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *ThemeCN) Font(style fyne.TextStyle) fyne.Resource {
	if style.Bold {
		return t.bold
	} else if style.Italic {
		return t.italic
	} else if style.Symbol || style.Monospace {
		return t.symbol
	} else {
		return t.symbol
		// return theme.DefaultTheme().Font(style)
	}

	// fmt.Println("style:", style)
	//return theme.DefaultTheme().Font(style)
	// return &fyne.StaticResource{
	// 	StaticName:    "simhei.ttf",
	// 	StaticContent: pfont.TTF_SIMHEI,
	// }
	// return &fyne.StaticResource{
	// 	StaticName:    "OPPOSans-R-2.ttf",
	// 	StaticContent: pfont.TTF_OPPOSans_R_2,
	// }
}

func (t *ThemeCN) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
