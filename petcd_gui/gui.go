package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"path"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/caticat/go_game_server/petcd"
	"github.com/caticat/go_game_server/petcd/pdata"
	ptheme "github.com/caticat/go_game_server/petcd_gui/theme"
	"github.com/caticat/go_game_server/phelp"
	"github.com/caticat/go_game_server/plog"
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
			fmt.Errorf("etcd Selected:%s", getConf().EndPointSelect),
			ErrorConnectToEtcdFailedInfo), w).Show()
	} else {
		setConnected(true)
		getFunUpdateTitle()()
	}

	w.ShowAndRun()
}

func initGUIHome(w fyne.Window) fyne.CanvasObject {
	// 数据
	root := getEtcdData()

	// 左下目录
	key := widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			if id == "" {
				return []string{root.GetKey()}
			} else {
				node := root.Get(id)
				if node == nil {
					return []string{}
				}
				return node.ChildKeys()
			}
		},
		func(id widget.TreeNodeID) bool {
			if node := root.Get(id); node == nil {
				return false
			} else {
				return node.IsBranch()
			}
		},
		func(branch bool) fyne.CanvasObject {
			if branch {
				return widget.NewLabel("Branch template")
			}
			return widget.NewLabel("Leaf template")
		},
		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
			text := path.Join(pdata.PDATA_PREFIX, path.Base(id))
			if branch {
				text += " (branch)"
			}
			o.(*widget.Label).SetText(text)
		})

	var search *widget.SelectEntry = nil
	var collapse *widget.Button = nil
	key.OnSelected = func(uid widget.TreeNodeID) {
		setEtcdKey(uid)
		value, ok := root.GetValue(uid)
		if ok {
			search.SetText(uid)
		} else {
			search.SetText("")
		}
		getEtcdValue().Set(value)

		// 展开父分支
		for p := uid; p != pdata.PDATA_PREFIX; {
			key.OpenBranch(path.Dir(p))
			p = path.Dir(p)
		}
		collapse.SetIcon(theme.ZoomOutIcon())
	}

	// 右下key对应的值
	value := widget.NewLabelWithData(g_etcdValue)
	main := container.NewHSplit(key, value)
	main.SetOffset(0.4)

	// 顶部刷新
	search = widget.NewSelectEntry(root.AllKeys())
	search.SetPlaceHolder("Search Key...")
	search.OnSubmitted = func(s string) {
		_, ok := root.GetValue(s)
		if !ok {
			dialog.NewError(ErrorPathHasNoData, w).Show()
			return
		}
		key.Select(s)
	}
	search.OnChanged = func(s string) {
		_, ok := root.GetValue(s)
		if !ok {
			return
		}
		search.OnSubmitted(s)
	}
	connectStatus := widget.NewIcon(theme.ConfirmIcon())
	collapse = widget.NewButtonWithIcon("", theme.ZoomInIcon(), func() {
		if key.IsBranchOpen(pdata.PDATA_PREFIX) {
			key.CloseAllBranches()
			collapse.SetIcon(theme.ZoomInIcon())
		} else {
			key.OpenAllBranches()
			collapse.SetIcon(theme.ZoomOutIcon())
		}
	})
	g_funGUIRefresh = func() {
		root.Clear()
		setEtcdKey("")
		getEtcdValue().Set("")
		initData()
		key.UnselectAll()
		key.Refresh()
		search.SetText("")
		search.SetOptions(root.AllKeys())
	}
	refresh := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), getFunGUIRefresh())
	edit := widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
		k := getEtcdKey()
		if k == "" {
			dialog.NewError(ErrorNoPathSelect, w).Show()
			return
		}
		v, _ := getEtcdValue().Get()
		bv := binding.BindString(&v)
		en := widget.NewMultiLineEntry()
		en.Bind(bv)
		di := dialog.NewForm("Edit", "OK", "Cancel",
			[]*widget.FormItem{widget.NewFormItem(k, en)}, func(b bool) {
				if !b {
					return
				}
				v, _ := bv.Get()
				petcd.PutKeepLease(k, v)
				refresh.OnTapped() // 刷新界面
				key.Select(k)      // 重新选择指定条目
			}, w)

		di.Resize(w.Canvas().Size())
		di.Show()
	})
	add := widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
		bk := binding.NewString()
		bv := binding.NewString()
		enk := widget.NewEntry()
		enk.Bind(bk)
		env := widget.NewMultiLineEntry()
		env.Bind(bv)
		di := dialog.NewForm("New", "OK", "Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("key", enk),
				widget.NewFormItem("value", env),
			}, func(b bool) {
				if !b {
					return
				}
				k, _ := bk.Get()
				v, _ := bv.Get()
				if k == "" {
					dialog.NewError(ErrorEmptyPath, w).Show()
					return
				}
				if !strings.HasPrefix(k, pdata.PDATA_PREFIX) {
					dialog.NewError(ErrorBadPathPrefix, w).Show()
					return
				}
				petcd.PutKeepLease(k, v)
				refresh.OnTapped() // 刷新界面
				key.Select(k)      // 重新选择指定条目
			}, w)

		di.Resize(w.Canvas().Size())
		di.Show()
	})
	del := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
		k := getEtcdKey()
		if k == "" {
			dialog.NewError(ErrorNoPathSelect, w).Show()
			return
		}
		_, ok := root.GetValue(k)
		if !ok {
			dialog.NewError(ErrorPathHasNoData, w).Show()
			return
		}
		dialog.NewConfirm("Delete", fmt.Sprintf("Delete %q ?(Single Node, Not include Hierarchies)", k), func(b bool) {
			if !b {
				return
			}
			petcd.Del(k)
			refresh.OnTapped() // 刷新界面
		}, w).Show()
	})
	title := container.NewGridWithColumns(2, search, container.NewHBox(connectStatus, collapse, refresh, edit, add, del))

	// 总布局
	all := container.NewBorder(title, nil, nil, nil, container.NewMax(main))

	// 更新标题状态函数
	g_funUpdateTitle = func() {
		updateWindowTitle(w)
		if getConnected() {
			connectStatus.SetResource(theme.ConfirmIcon())
		} else {
			connectStatus.SetResource(theme.ErrorIcon())
		}
	}

	return all
}

func initGUISetting(w fyne.Window) fyne.CanvasObject {
	conf := getConf()

	// base
	lbs := widget.NewForm(
		widget.NewFormItem("dial-timeout", initGUISettingForm(conf.Base.DialTimeout, func(v int64) {
			conf.Base.DialTimeout = v
			conf.SetBase(conf.Base)
		})),
		widget.NewFormItem("operation-timeout", initGUISettingForm(conf.Base.OperationTimeout, func(v int64) {
			conf.Base.OperationTimeout = v
			conf.SetBase(conf.Base)
		})),
		widget.NewFormItem("lease-timeout-before-keep-aliv", initGUISettingForm(conf.Base.LeaseTimeoutBeforeKeepAlive, func(v int64) {
			conf.Base.LeaseTimeoutBeforeKeepAlive = v
			conf.SetBase(conf.Base)
		})),
		widget.NewFormItem("reset-local-data", widget.NewButtonWithIcon("reset all local config data", theme.ErrorIcon(), func() {
			dialog.NewConfirm("reset-local-data", "Are you serious?\nDelete All Local Storage Configuration Data", func(b bool) {
				if !b {
					return
				}
				getConf().clearPreferences()
				d := dialog.NewInformation("reset-local-data", "Delete Local Config Done!\n Need Restart Progress Now!", w)
				d.SetOnClosed(func() {
					getApp().Quit()
				})
				d.Show()
			}, w).Show()
		})),
	)

	// endpoint
	lep := widget.NewLabel("default text")
	var search *widget.Select
	var butConnect *widget.Button
	butConnect = widget.NewButtonWithIcon("Connect", theme.LoginIcon(), func() { // 连接按钮前移,解决sarch的空指针问题
		k := search.Selected
		_, ok := conf.MapEndPoint.Get(k)
		if !ok {
			dialog.NewError(ErrorPathHasNoData, w).Show()
			return
		}
		if err := conf.SetSelect(k); err != nil {
			dialog.NewError(err, w).Show()
			return
		}
		butConnect.Disable()
		dialog.NewInformation("Information", fmt.Sprintf("%s now is Connecting to: %q", WINDOW_TITLE, k), w).Show()
	})
	search = widget.NewSelect(conf.MapEndPoint.M_sliKey, func(s string) {
		c, ok := conf.MapEndPoint.Get(s)
		if !ok {
			return
		}
		lep.SetText(strings.Join(c, "\n"))
		if conf.EndPointSelect == s {
			butConnect.Disable()
		} else {
			butConnect.Enable()
		}
	})
	search.SetSelected(conf.EndPointSelect)
	butEdit := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		k := search.Selected
		v, ok := conf.MapEndPoint.Get(k)
		if !ok {
			dialog.NewError(ErrorPathHasNoData, w).Show()
			return
		}
		vOri := phelp.ToJsonIndent(v)
		bv := binding.NewString()
		bv.Set(vOri)
		en := widget.NewMultiLineEntry()
		en.Bind(bv)
		en.SetPlaceHolder("json format")
		en.Validator = func(s string) error {
			if json.Valid([]byte(s)) {
				return nil
			}
			return ErrorInvalidInput
		}
		di := dialog.NewForm("Edit", "OK", "Cancel",
			[]*widget.FormItem{widget.NewFormItem(k, en)}, func(b bool) {
				if !b {
					return
				}
				v, _ := bv.Get()
				if v == vOri {
					return
				}

				var c []string
				if err := json.Unmarshal([]byte(v), &c); err != nil {
					dialog.NewError(err, w).Show()
					return
				}

				if err := conf.ModEndPoint(k, c); err != nil {
					dialog.NewError(err, w).Show()
					return
				}

				search.OnChanged(k)
			}, w)

		di.Resize(w.Canvas().Size())
		di.Show()
	})
	butCreate := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		bk := binding.NewString()
		ek := widget.NewEntry()
		ek.Bind(bk)

		bv := binding.NewString()
		bv.Set("[\n\n]")
		en := widget.NewMultiLineEntry()
		en.Bind(bv)
		en.SetPlaceHolder("json format")
		en.Validator = func(s string) error {
			if json.Valid([]byte(s)) {
				return nil
			}
			return ErrorInvalidInput
		}
		di := dialog.NewForm("Add", "OK", "Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("EndPointName", ek),
				widget.NewFormItem("EndPoints", en),
			}, func(b bool) {
				if !b {
					return
				}
				k, _ := bk.Get()
				v, _ := bv.Get()
				if k == "" || v == "" {
					dialog.NewError(ErrorInputDataEmpty, w).Show()
					return
				}
				if _, ok := conf.MapEndPoint.Get(k); ok {
					dialog.NewError(ErrorEndPointNameAlreadyExist, w).Show()
					return
				}

				var c []string
				if err := json.Unmarshal([]byte(v), &c); err != nil {
					dialog.NewError(err, w).Show()
					return
				}

				conf.AddEndPoint(k, c)

				search.Options = conf.MapEndPoint.M_sliKey
				search.Refresh()
				search.SetSelected(k)
			}, w)

		di.Resize(w.Canvas().Size())
		di.Show()
	})
	butDelete := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		k := search.Selected
		_, ok := conf.MapEndPoint.Get(k)
		if !ok {
			dialog.NewError(ErrorPathHasNoData, w).Show()
			return
		}
		di := dialog.NewConfirm("Delete", "Delete EtcdPointName:"+k, func(b bool) {
			if !b {
				return
			}

			if err := conf.DelEndPoint(k); err != nil {
				dialog.NewError(err, w).Show()
				return
			}

			search.Options = conf.MapEndPoint.M_sliKey
			search.Refresh()
			search.ClearSelected()
		}, w)

		di.Show()
	})

	header := container.NewGridWithColumns(2, search, container.NewHBox(butConnect, butEdit, butCreate, butDelete))
	cep := container.NewBorder(header, nil, nil, nil, lep)

	con := container.NewAppTabs(
		container.NewTabItem("EndPoint", cep),
		container.NewTabItem("Base", lbs),
	)

	return con
}

func initGUISettingForm(value int64, confirm func(int64)) *fyne.Container {
	bindData := binding.NewString()
	bindData.Set(fmt.Sprintf("%v", value))
	ent := widget.NewEntryWithData(bindData)
	ent.Disable()
	var butEdit, butOK, butCancel *widget.Button
	butEdit = widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		s, _ := bindData.Get()
		value, _ = strconv.ParseInt(s, 10, 64)

		ent.Enable()
		butEdit.Hide()
		butOK.Show()
		butCancel.Show()
	})
	butOK = widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
		if s, err := bindData.Get(); err == nil {
			if v, err := strconv.ParseInt(s, 10, 64); err == nil {
				if v != value {
					confirm(v)
				}
			} else {
				plog.Error(err)
			}
		} else {
			plog.Error(err)
		}
		ent.Disable()
		butEdit.Show()
		butOK.Hide()
		butCancel.Hide()
	})
	butCancel = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		bindData.Set(fmt.Sprintf("%v", value))
		ent.Disable()
		butEdit.Show()
		butOK.Hide()
		butCancel.Hide()
	})
	butOK.Hide()
	butCancel.Hide()

	return container.NewGridWithColumns(2, ent, container.NewHBox(butEdit, butOK, butCancel))
}

func initGUIInfo(w fyne.Window, c binding.String) fyne.CanvasObject {
	return container.NewVBox(widget.NewLabel("Config"), canvas.NewLine(color.Black), widget.NewLabelWithData(c))
}

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
