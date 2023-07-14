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

	// 设置Home刷新
	var guiSelSearch *widget.SelectEntry = nil // 顶部搜索栏
	var guiButCollapse *widget.Button = nil    // 折叠目录按钮
	var guiTreKeys *widget.Tree = nil          // 左下目录
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

	// 左下目录
	initGUIHomeKeys(&guiTreKeys, &guiSelSearch, &guiButCollapse)

	// 右下key对应的值
	guiLabValue := widget.NewLabelWithData(g_etcdValue)

	// 主界面
	guiHSpMain := container.NewHSplit(guiTreKeys, guiLabValue)
	guiHSpMain.SetOffset(GUI_HOME_MAIN_OFFSET)

	// 顶部刷新
	guiSelSearch = widget.NewSelectEntry(root.AllKeys())
	guiSelSearch.SetPlaceHolder(GUI_HOME_SEARCH_PLACEHOLDER)
	guiSelSearch.OnSubmitted = func(s string) {
		if _, ok := root.GetValue(s); !ok {
			dialog.NewError(ErrorPathHasNoData, w).Show()
			return
		}
		guiTreKeys.Select(s)
	}
	guiSelSearch.OnChanged = func(s string) {
		if _, ok := root.GetValue(s); !ok {
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
		key := getEtcdKey()
		if key == STR_EMPTY {
			dialog.NewError(ErrorNoPathSelect, w).Show()
			return
		}
		v, _ := getEtcdValue().Get()
		binV := binding.BindString(&v)
		guiEntV := widget.NewMultiLineEntry()
		guiEntV.Bind(binV)
		guiEntV.SetMinRowsVisible(GUI_HOME_EDIT_ENTRY_LINE_NUM)
		guiDia := dialog.NewForm("Edit", "OK", "Cancel",
			[]*widget.FormItem{widget.NewFormItem(key, guiEntV)}, func(b bool) {
				if !b {
					return
				}
				v, _ := binV.Get()
				petcd.PutKeepLease(key, v)
				guiButRefresh.OnTapped() // 刷新界面
				guiTreKeys.Select(key)   // 重新选择指定条目
			}, w)

		guiDia.Resize(w.Canvas().Size())
		guiDia.Show()
	})
	guiButAdd := widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
		k := getEtcdKey()
		binK := binding.NewString()
		binK.Set(path.Join(k, pdata.PDATA_PREFIX)) // 初始化输入参数
		binV := binding.NewString()
		guiEntKey := widget.NewEntry()
		guiEntKey.Bind(binK)
		guiEntKey.Validator = func(s string) error {
			if _, ok := root.GetValue(s); ok {
				return ErrorPathAlreadyHasData
			}
			return nil
		}
		guiEntV := widget.NewMultiLineEntry()
		guiEntV.Bind(binV)
		guiEntV.SetMinRowsVisible(GUI_HOME_EDIT_ENTRY_LINE_NUM)
		guiDia := dialog.NewForm("New", "OK", "Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("key", guiEntKey),
				widget.NewFormItem("value", guiEntV),
			}, func(b bool) {
				if !b {
					return
				}
				k, _ := binK.Get()
				v, _ := binV.Get()
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

		guiDia.Resize(w.Canvas().Size())
		guiDia.Show()
	})
	guiButDel := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
		k := getEtcdKey()
		if k == STR_EMPTY {
			dialog.NewError(ErrorNoPathSelect, w).Show()
			return
		}
		if _, ok := root.GetValue(k); !ok {
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
			guiIco := widget.NewIcon(theme.RadioButtonIcon())
			guiLab := widget.NewLabel("")
			return container.NewHBox(guiIco, guiLab)
		},
		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
			if _, ok := root.GetValue(id); ok {
				o.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(theme.RadioButtonCheckedIcon())
			} else {
				o.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(theme.RadioButtonIcon())
			}
			text := path.Join(pdata.PDATA_PREFIX, path.Base(id))
			o.(*fyne.Container).Objects[1].(*widget.Label).SetText(text)
		})

	(*pKeys).OnSelected = func(uid widget.TreeNodeID) {
		setEtcdKey(uid)
		value, ok := root.GetValue(uid)
		if ok {
			(*pSearch).SetText(uid)
		} else {
			(*pSearch).SetText(STR_EMPTY)
			value = STR_NIL

			// 没有数据的话,相当于切换子树的展开状态,这里重复点击无效,体验不好,所以注释掉了
			// if (*pKeys).IsBranchOpen(uid) {
			// 	(*pKeys).CloseBranch(uid)
			// } else {
			// 	(*pKeys).OpenBranch(uid)
			// }
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
