package main

import (
	"fmt"
	"image/color"
	"path"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/caticat/go_game_server/petcd"
	"github.com/caticat/go_game_server/petcd/pdata"
	ptheme "github.com/caticat/go_game_server/petcd_gui/theme"
	"github.com/caticat/go_game_server/plog"
)

func runGUI() {
	a := app.NewWithID("com.github.caticat.go_game_server.petcd_gui")
	a.Settings().SetTheme(ptheme.NewThemeCN())
	w := a.NewWindow("PEtcdGUI")
	w.SetMaster()

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), initGUIHome(w)),
		container.NewTabItemWithIcon("Setting", theme.SettingsIcon(), initGUISetting(w)),
		container.NewTabItemWithIcon("Info", theme.InfoIcon(), initGUIInfo(w)),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(1000, 700))
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
	key.OnSelected = func(uid widget.TreeNodeID) {
		setEtcdKey(uid)
		value, ok := root.GetValue(uid)
		if ok {
			search.SetText(uid)
		} else {
			search.SetText("")
		}
		getEtcdValue().Set(value)
	}

	// 右下key对应的值
	value := widget.NewLabelWithData(g_etcdValue)
	main := container.NewHSplit(key, value)
	main.SetOffset(0.4)

	// 顶部刷新
	search = widget.NewSelectEntry(root.AllKeys())
	search.SetPlaceHolder("Search Key...")
	search.OnSubmitted = func(s string) {
		key.Select(s)
	}
	search.OnChanged = func(s string) {
		_, ok := root.GetValue(s)
		if !ok {
			return
		}
		search.OnSubmitted(s)
	}
	refresh := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		root.Clear()
		setEtcdKey("")
		getEtcdValue().Set("")
		initData()
		key.UnselectAll()
		key.Refresh()
		search.SetText("")
	})
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
	title := container.NewGridWithColumns(2, search, container.NewHBox(refresh, edit, add, del))

	// 总布局
	all := container.NewBorder(title, nil, nil, nil, container.NewMax(main))
	return all
}

func initGUISetting(w fyne.Window) fyne.CanvasObject {
	conf := getConf()

	endpoints := widget.NewMultiLineEntry()
	endpoints.SetText(strings.Join(conf.Etcd.Endpoints, "\n"))

	form := widget.NewForm(widget.NewFormItem("EndPoints", endpoints))
	form.OnSubmit = func() {
		// plog.InfoLn("TODO 保存")
		dialog.NewInformation("保存", "TODO:功能待制作", w).Show()
	}
	form.OnCancel = func() {
		// plog.InfoLn("TODO 取消")
		dialog.NewInformation("取消", "TODO:功能待制作", w).Show()
	}

	return form
}

func initGUIInfo(w fyne.Window) fyne.CanvasObject {
	conf := getConf()
	c, err := conf.ToJson()
	if err != nil {
		plog.ErrorLn(err)
	}

	return container.NewVBox(widget.NewLabel("Config"), canvas.NewLine(color.Black), widget.NewLabel(c))
	// form := widget.NewForm(widget.NewFormItem("Config", widget.NewLabel(c)))

	// return form
}
