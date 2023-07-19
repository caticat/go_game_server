package main

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"github.com/caticat/go_game_server/pfyne_theme_cn"
	"github.com/caticat/go_game_server/phelp"
	"github.com/caticat/go_game_server/plog"
)

func runGUI(err error) {
	// 窗口初始化
	a := getApp()
	a.Settings().SetTheme(pfyne_theme_cn.NewThemeCN())
	w := a.NewWindow(WINDOW_TITLE)
	w.SetMaster()

	// 窗口内容
	infoData := binding.NewString()
	guiLog := initGUILog(w)
	guiTabMain := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), initGUIHome(w)),
		container.NewTabItemWithIcon("Log", theme.DocumentIcon(), guiLog),
		container.NewTabItemWithIcon("Setting", theme.SettingsIcon(), initGUISetting(w)),
		container.NewTabItemWithIcon("Info", theme.InfoIcon(), initGUIInfo(w, infoData)),
	)
	guiTabMain.OnSelected = func(ti *container.TabItem) {
		iconName := ti.Icon.Name()
		if iconName == theme.InfoIcon().Name() { // 更新Info界面的数据
			infoData.Set(phelp.ToJsonIndent(getConf()))
		}
		if iconName == theme.DocumentIcon().Name() { // 日志界面隐藏逻辑,减少界面刷新
			guiLog.Show()
		} else {
			guiLog.Hide()
		}
	}
	guiTabMain.SetTabLocation(container.TabLocationLeading)

	// 窗口尺寸
	w.SetContent(guiTabMain)
	w.Resize(fyne.NewSize(GUI_WINDOW_INIT_SIZE_W, GUI_WINDOW_INIT_SIZE_H))

	// 初始连接错误提示框
	if err != nil {
		setConnected(false)
		getFunUpdateTitle()()
		e := errors.Join(ErrorConnectToEtcdFailed,
			err,
			fmt.Errorf("etcd Selected:%s", getConf().ConnSelect),
			ErrorConnectToEtcdFailedInfo)
		dialog.NewError(e, w).Show()
		plog.ErrorLn(e)
	} else {
		setConnected(true)
		getFunUpdateTitle()()
	}

	// 阻塞运行
	w.ShowAndRun()
}
