package main

import (
	"fmt"
	"path"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/caticat/go_game_server/petcd"
	"github.com/caticat/go_game_server/petcd/pdata"
)

func initGUIHome(w fyne.Window) fyne.CanvasObject {
	// 数据
	root := getEtcdData()

	var guiSelSearch *widget.SelectEntry = nil
	var guiButCollapse *widget.Button = nil
	// 左下目录
	var guiTreKeys *widget.Tree = nil
	initGUIHomeKeys(&guiTreKeys, &guiSelSearch, &guiButCollapse)

	// 右下key对应的值
	guiLabValue := widget.NewLabelWithData(g_etcdValue)
	guiHSpMain := container.NewHSplit(guiTreKeys, guiLabValue)
	guiHSpMain.SetOffset(GUI_HOME_MAIN_OFFSET)

	// 顶部刷新
	guiSelSearch = widget.NewSelectEntry(root.AllKeys())
	guiSelSearch.SetPlaceHolder(GUI_HOME_SEARCH_PLACEHOLDER)
	guiSelSearch.OnSubmitted = func(s string) {
		_, ok := root.GetValue(s)
		if !ok {
			dialog.NewError(ErrorPathHasNoData, w).Show()
			return
		}
		guiTreKeys.Select(s)
	}
	guiSelSearch.OnChanged = func(s string) {
		_, ok := root.GetValue(s)
		if !ok {
			return
		}
		guiSelSearch.OnSubmitted(s)
	}
	connectStatus := widget.NewIcon(theme.ConfirmIcon())
	guiButCollapse = widget.NewButtonWithIcon(STR_EMPTY, theme.ZoomInIcon(), func() {
		if guiTreKeys.IsBranchOpen(pdata.PDATA_PREFIX) {
			guiTreKeys.CloseAllBranches()
			guiButCollapse.SetIcon(theme.ZoomInIcon())
		} else {
			guiTreKeys.OpenAllBranches()
			guiButCollapse.SetIcon(theme.ZoomOutIcon())
		}
	})
	guiButRefresh := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), getFunGUIHomeRefresh())
	guiButEdit := widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
		k := getEtcdKey()
		if k == STR_EMPTY {
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
				guiButRefresh.OnTapped() // 刷新界面
				guiTreKeys.Select(k)     // 重新选择指定条目
			}, w)

		di.Resize(w.Canvas().Size())
		di.Show()
	})
	guiButAdd := widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
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
				if k == STR_EMPTY {
					dialog.NewError(ErrorEmptyPath, w).Show()
					return
				}
				if !strings.HasPrefix(k, pdata.PDATA_PREFIX) {
					dialog.NewError(ErrorBadPathPrefix, w).Show()
					return
				}
				petcd.PutKeepLease(k, v)
				guiButRefresh.OnTapped() // 刷新界面
				guiTreKeys.Select(k)     // 重新选择指定条目
			}, w)

		di.Resize(w.Canvas().Size())
		di.Show()
	})
	guiButDel := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
		k := getEtcdKey()
		if k == STR_EMPTY {
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
			guiButRefresh.OnTapped() // 刷新界面
		}, w).Show()
	})
	guiConTitle := container.NewGridWithColumns(2, guiSelSearch, container.NewHBox(connectStatus, guiButCollapse, guiButRefresh, guiButEdit, guiButAdd, guiButDel))

	// 总布局
	guiConAll := container.NewBorder(guiConTitle, nil, nil, nil, container.NewMax(guiHSpMain))

	// 设置Home刷新
	setFunGUIHomeRefresh(func() {
		root.Clear()
		setEtcdKey(STR_EMPTY)
		getEtcdValue().Set(STR_EMPTY)
		initData()
		guiTreKeys.UnselectAll()
		guiTreKeys.Refresh()
		guiSelSearch.SetText(STR_EMPTY)
		guiSelSearch.SetOptions(root.AllKeys())
	})

	// 设置标题状态函数
	setFunUpdateTitle(func() {
		updateWindowTitle(w)
		if getConnected() {
			connectStatus.SetResource(theme.ConfirmIcon())
		} else {
			connectStatus.SetResource(theme.ErrorIcon())
		}
	})

	return guiConAll
}

func initGUIHomeKeys(pKeys **widget.Tree, pSearch **widget.SelectEntry, pCollapse **widget.Button) {
	root := getEtcdData()

	*pKeys = widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			if id == STR_EMPTY {
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

	(*pKeys).OnSelected = func(uid widget.TreeNodeID) {
		setEtcdKey(uid)
		value, ok := root.GetValue(uid)
		if ok {
			(*pSearch).SetText(uid)
		} else {
			(*pSearch).SetText(STR_EMPTY)
		}
		getEtcdValue().Set(value)

		// 展开父分支
		for p := uid; p != pdata.PDATA_PREFIX; {
			(*pKeys).OpenBranch(path.Dir(p))
			p = path.Dir(p)
		}
		(*pCollapse).SetIcon(theme.ZoomOutIcon())
	}
}
