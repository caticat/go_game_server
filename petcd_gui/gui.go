package main

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	ptheme "github.com/caticat/go_game_server/petcd_gui/theme"
	"github.com/caticat/go_game_server/phelp"
)

func runGUI(err error) {
	a := getApp()
	a.Settings().SetTheme(ptheme.NewThemeCN())
	w := a.NewWindow(WINDOW_TITLE)
	w.SetMaster()

	infoData := binding.NewString()
	l := initGUILog(w)
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), initGUIHome(w)),
		container.NewTabItemWithIcon("Setting", theme.SettingsIcon(), initGUISetting(w)),
		container.NewTabItemWithIcon("Info", theme.InfoIcon(), initGUIInfo(w, infoData)),
		container.NewTabItemWithIcon("Log", theme.DocumentIcon(), l),
	)
	tabs.OnSelected = func(ti *container.TabItem) {
		iconName := ti.Icon.Name()
		if iconName == theme.InfoIcon().Name() { // 更新Info界面的数据
			infoData.Set(phelp.ToJsonIndent(getConf()))
		}
		if iconName == theme.DocumentIcon().Name() { // 日志界面隐藏逻辑,减少界面刷新
			l.Show()
		} else {
			l.Hide()
		}
	}
	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(1000, 700))

	// 初始连接错误提示框
	if err != nil {
		setConnected(false)
		getFunUpdateTitle()()
		dialog.NewError(errors.Join(ErrorConnectToEtcdFailed,
			err,
			fmt.Errorf("etcd Selected:%s", getConf().ConnSelect),
			ErrorConnectToEtcdFailedInfo), w).Show()
	} else {
		setConnected(true)
		getFunUpdateTitle()()
	}

	w.ShowAndRun()
}
