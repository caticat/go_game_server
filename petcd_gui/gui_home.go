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
	"github.com/caticat/go_game_server/plog"
)

func initGUIHome(w fyne.Window) fyne.CanvasObject {
	// 数据
	root := getEtcdData()

	// 控件 先声明的原因是要把他们提出来作为其他控件的函数调用参数,否则放在后面顺序不好控制
	var (
		guiSelSearch   *widget.SelectEntry = nil // 头部 搜索栏
		connectStatus  *widget.Icon        = nil // 头部 连接状态
		guiButCollapse *widget.Button      = nil // 头部 折叠目录
		guiButRefresh  *widget.Button      = nil // 头部 刷新
		guiButEdit     *widget.Button      = nil // 头部 编辑
		guiButAdd      *widget.Button      = nil // 头部 添加
		guiButDel      *widget.Button      = nil // 头部 删除
		guiTreKeys     *widget.Tree        = nil // 身体 目录
		guiLabValue    *widget.Label       = nil // 身体 值显示
		guiLabLogLast  *widget.Label       = nil // 身体 日志最后一行
		guiConHead     *fyne.Container     = nil // 布局 头
		guiConBody     *container.Split    = nil // 布局 身
		guiConAll      *fyne.Container     = nil // 布局 总
	)

	// 设置Home刷新 设置函数要放在界面初始化前面
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

	// 身体
	initGUIHomeBody(
		&guiTreKeys,
		&guiLabValue,
		&guiLabLogLast,
		&guiConBody,
		&guiSelSearch,
		&guiButCollapse,
	)

	// 头部
	initGUIHomeHead(
		w,
		&guiSelSearch,
		&connectStatus,
		&guiButCollapse,
		&guiButRefresh,
		&guiButEdit,
		&guiButAdd,
		&guiButDel,
		&guiConHead,
		guiTreKeys,
	)

	// 设置标题状态函数
	setFunUpdateTitle(func() {
		updateWindowTitle(w)
		if getConnected() {
			connectStatus.SetResource(theme.ConfirmIcon())
		} else {
			connectStatus.SetResource(theme.ErrorIcon())
		}
	})

	// 总布局
	guiConAll = container.NewBorder(guiConHead, guiLabLogLast, nil, nil, container.NewMax(guiConBody))
	return guiConAll
}

func initGUIHomeHead(w fyne.Window,
	pGuiSelSearch **widget.SelectEntry,
	pConnectStatus **widget.Icon,
	pGuiButCollapse **widget.Button,
	pGuiButRefresh **widget.Button,
	pGuiButEdit **widget.Button,
	pGuiButAdd **widget.Button,
	pGuiButDel **widget.Button,
	pGuiConHead **fyne.Container,
	guiTreKeys *widget.Tree) {
	// 数据
	root := getEtcdData()

	// 查询
	*pGuiSelSearch = widget.NewSelectEntry(root.AllKeys())
	(*pGuiSelSearch).SetPlaceHolder(GUI_HOME_SEARCH_PLACEHOLDER)
	(*pGuiSelSearch).OnSubmitted = func(s string) {
		if _, ok := root.GetValue(s); !ok {
			dialog.NewError(ErrorPathHasNoData, w).Show()
			plog.ErrorLn(ErrorPathHasNoData)
			return
		}
		guiTreKeys.Select(s)
	}
	(*pGuiSelSearch).OnChanged = func(s string) {
		if _, ok := root.GetValue(s); !ok {
			return
		}
		(*pGuiSelSearch).OnSubmitted(s)
	}

	// 连接状态
	*pConnectStatus = widget.NewIcon(theme.ConfirmIcon())

	// 目录折叠
	*pGuiButCollapse = widget.NewButtonWithIcon(STR_EMPTY, theme.ZoomInIcon(), func() {
		if guiTreKeys.IsBranchOpen(pdata.PDATA_PREFIX) {
			guiTreKeys.CloseAllBranches()
			(*pGuiButCollapse).SetIcon(theme.ZoomInIcon())
		} else {
			guiTreKeys.OpenAllBranches()
			(*pGuiButCollapse).SetIcon(theme.ZoomOutIcon())
		}
	})

	// 数据刷新
	*pGuiButRefresh = widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), getFunGUIHomeRefresh())

	// 编辑
	*pGuiButEdit = widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
		key := getEtcdKey()
		if key == STR_EMPTY {
			dialog.NewError(ErrorNoPathSelect, w).Show()
			plog.ErrorLn(ErrorNoPathSelect)
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
				(*pGuiButRefresh).OnTapped() // 刷新界面
				guiTreKeys.Select(key)       // 重新选择指定条目
				plog.InfoF("put(mod) %q %q\n", key, v)
			}, w)

		guiDia.Resize(w.Canvas().Size())
		guiDia.Show()
	})

	// 添加
	*pGuiButAdd = widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
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
					plog.ErrorLn(ErrorEmptyPath)
					return
				}
				if !strings.HasPrefix(k, pdata.PDATA_PREFIX) {
					dialog.NewError(ErrorBadPathPrefix, w).Show()
					plog.ErrorLn(ErrorBadPathPrefix)
					return
				}
				petcd.PutKeepLease(k, v)
				(*pGuiButRefresh).OnTapped() // 刷新界面
				guiTreKeys.Select(k)         // 重新选择指定条目
				plog.InfoF("put(add) %q %q\n", k, v)
			}, w)

		guiDia.Resize(w.Canvas().Size())
		guiDia.Show()
	})

	// 删除
	*pGuiButDel = widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
		k := getEtcdKey()
		if k == STR_EMPTY {
			dialog.NewError(ErrorNoPathSelect, w).Show()
			plog.ErrorLn(ErrorNoPathSelect)
			return
		}
		if _, ok := root.GetValue(k); !ok {
			dialog.NewError(ErrorPathHasNoData, w).Show()
			plog.ErrorLn(ErrorPathHasNoData)
			return
		}
		dialog.NewConfirm("Delete", fmt.Sprintf("Delete %q ?(Single Node, Not include Hierarchies)", k), func(b bool) {
			if !b {
				return
			}
			petcd.Del(k)
			(*pGuiButRefresh).OnTapped() // 刷新界面
			plog.InfoF("del %q\n", k)
		}, w).Show()
	})

	// 容器 头
	*pGuiConHead = container.NewGridWithColumns(2, *pGuiSelSearch, container.NewHBox(*pConnectStatus, *pGuiButCollapse, *pGuiButRefresh, *pGuiButEdit, *pGuiButAdd, *pGuiButDel))
}

func initGUIHomeBody(pGuiTreKeys **widget.Tree,
	pGuiLabValue **widget.Label,
	pGuiLabLogLast **widget.Label,
	pGuiConBody **container.Split,
	pGuiSelSearch **widget.SelectEntry,
	pGuiButCollapse **widget.Button) {
	// 目录
	initGUIHomeBodyKeys(pGuiTreKeys, pGuiSelSearch, pGuiButCollapse)

	// 数据值
	*pGuiLabValue = widget.NewLabelWithData(g_etcdValue)

	// 日志 最后一行
	*pGuiLabLogLast = widget.NewLabelWithData(getLogLast())

	// 容器 身
	*pGuiConBody = container.NewHSplit(*pGuiTreKeys, *pGuiLabValue)
	(*pGuiConBody).SetOffset(GUI_HOME_MAIN_OFFSET)
}

func initGUIHomeBodyKeys(pGuiTreKeys **widget.Tree,
	pGuiSelSearch **widget.SelectEntry,
	pGuiButCollapse **widget.Button) {
	// 数据
	root := getEtcdData()

	// 目录
	*pGuiTreKeys = widget.NewTree(
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

	(*pGuiTreKeys).OnSelected = func(uid widget.TreeNodeID) {
		setEtcdKey(uid)
		value, ok := root.GetValue(uid)
		if ok {
			(*pGuiSelSearch).SetText(uid)
		} else {
			(*pGuiSelSearch).SetText(STR_EMPTY)
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
			(*pGuiTreKeys).OpenBranch(path.Dir(p))
			p = path.Dir(p)
		}
		(*pGuiButCollapse).SetIcon(theme.ZoomOutIcon())
	}
}
