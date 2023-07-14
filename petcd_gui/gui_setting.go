package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/caticat/go_game_server/petcd"
	"github.com/caticat/go_game_server/phelp"
	"github.com/caticat/go_game_server/plog"
)

func initGUISetting(w fyne.Window) fyne.CanvasObject {
	// 控件
	var (
		guiForBase       *widget.Form       = nil // 基础 基础页签
		guiLabConnection *widget.Label      = nil // 连接 配置信息
		guiSelConnection *widget.Select     = nil // 连接 选择框
		guiButConnection *widget.Button     = nil // 连接 连接
		guiButEdit       *widget.Button     = nil // 连接 编辑
		guiButCreate     *widget.Button     = nil // 连接 创建
		guiButDelete     *widget.Button     = nil // 连接 删除
		guiGriHeader     *fyne.Container    = nil // 布局 头
		guiBorConnection *fyne.Container    = nil // 布局 连接页签
		guiConTab        *container.AppTabs = nil // 布局 总
	)

	// base
	initGUISettingBase(w, &guiForBase)

	// conn
	initGUISettingConn(w,
		&guiLabConnection,
		&guiSelConnection,
		&guiButConnection,
		&guiButEdit,
		&guiButCreate,
		&guiButDelete,
		&guiGriHeader,
		&guiBorConnection,
	)

	// 页签
	guiConTab = container.NewAppTabs(
		container.NewTabItem("Connection", guiBorConnection),
		container.NewTabItem("Base", guiForBase),
	)

	return guiConTab
}

func initGUISettingBase(w fyne.Window, pGuiForBase **widget.Form) {
	conf := getConf()
	*pGuiForBase = widget.NewForm(
		widget.NewFormItem("dial-timeout", initGUISettingBaseFormItem(conf.Base.DialTimeout, func(v int64) {
			conf.Base.DialTimeout = v
			conf.SetBase(conf.Base)
		})),
		widget.NewFormItem("operation-timeout", initGUISettingBaseFormItem(conf.Base.OperationTimeout, func(v int64) {
			conf.Base.OperationTimeout = v
			conf.SetBase(conf.Base)
		})),
		widget.NewFormItem("lease-timeout-before-keep-aliv", initGUISettingBaseFormItem(conf.Base.LeaseTimeoutBeforeKeepAlive, func(v int64) {
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
					close(true)
				})
				d.Show()
			}, w).Show()
		})),
	)
}

func initGUISettingBaseFormItem(value int64, confirm func(int64)) *fyne.Container {
	bindData := binding.NewString()
	bindData.Set(fmt.Sprintf("%v", value))
	guiEntValue := widget.NewEntryWithData(bindData)
	guiEntValue.Disable()
	var guiButEdit, guiButOK, guiButCancel *widget.Button
	guiButEdit = widget.NewButtonWithIcon(STR_EMPTY, theme.DocumentCreateIcon(), func() {
		s, _ := bindData.Get()
		value, _ = strconv.ParseInt(s, 10, 64)

		guiEntValue.Enable()
		guiButEdit.Hide()
		guiButOK.Show()
		guiButCancel.Show()
	})
	guiButOK = widget.NewButtonWithIcon(STR_EMPTY, theme.ConfirmIcon(), func() {
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
		guiEntValue.Disable()
		guiButEdit.Show()
		guiButOK.Hide()
		guiButCancel.Hide()
	})
	guiButCancel = widget.NewButtonWithIcon(STR_EMPTY, theme.CancelIcon(), func() {
		bindData.Set(fmt.Sprintf("%v", value))
		guiEntValue.Disable()
		guiButEdit.Show()
		guiButOK.Hide()
		guiButCancel.Hide()
	})
	guiButOK.Hide()
	guiButCancel.Hide()

	return container.NewGridWithColumns(2, guiEntValue, container.NewHBox(guiButEdit, guiButOK, guiButCancel))
}

func initGUISettingConn(w fyne.Window,
	pGuiLabConnection **widget.Label,
	pGuiSelConnection **widget.Select,
	pGuiButConnection **widget.Button,
	pGuiButEdit **widget.Button,
	pGuiButCreate **widget.Button,
	pGuiButDelete **widget.Button,
	pGuiGriHeader **fyne.Container,
	pGuiBorConnection **fyne.Container,
) {
	conf := getConf()
	*pGuiLabConnection = widget.NewLabel("default text")
	*pGuiButConnection = widget.NewButtonWithIcon("Connect", theme.LoginIcon(), func() { // 连接按钮前移,解决sarch的空指针问题
		k := (*pGuiSelConnection).Selected
		_, ok := conf.MapConn.Get(k)
		if !ok {
			dialog.NewError(ErrorPathHasNoData, w).Show()
			return
		}
		if err := conf.SetSelect(k); err != nil {
			dialog.NewError(err, w).Show()
			return
		}
		(*pGuiButConnection).Disable()
		dialog.NewInformation("Information", fmt.Sprintf("%s now is Connecting to: %q", WINDOW_TITLE, k), w).Show()
	})
	*pGuiSelConnection = widget.NewSelect(conf.MapConn.M_sliKey, func(s string) {
		c, ok := conf.MapConn.Get(s)
		if !ok {
			return
		}
		(*pGuiLabConnection).SetText(phelp.ToJsonIndent(c))
		if conf.ConnSelect == s {
			(*pGuiButConnection).Disable()
		} else {
			(*pGuiButConnection).Enable()
		}
	})
	(*pGuiSelConnection).SetSelected(conf.ConnSelect)
	*pGuiButEdit = widget.NewButtonWithIcon(STR_EMPTY, theme.DocumentCreateIcon(), func() {
		k := (*pGuiSelConnection).Selected
		v, ok := conf.MapConn.Get(k)
		if !ok {
			dialog.NewError(ErrorPathHasNoData, w).Show()
			return
		}
		vOri := phelp.ToJsonIndent(v)
		binV := binding.NewString()
		binV.Set(vOri)
		guiEntV := widget.NewMultiLineEntry()
		guiEntV.Bind(binV)
		guiEntV.SetPlaceHolder("json format")
		guiEntV.SetMinRowsVisible(GUI_SETTING_EDIT_CONN_ENTRY_LINE_NUM)
		guiEntV.Validator = func(s string) error {
			if json.Valid([]byte(s)) {
				return nil
			}
			return ErrorInputNeedJsonFormat
		}
		d := dialog.NewForm("Edit", "OK", "Cancel",
			[]*widget.FormItem{widget.NewFormItem(k, guiEntV)}, func(b bool) {
				if !b {
					return
				}
				v, _ := binV.Get()
				if v == vOri {
					return
				}

				c := petcd.NewConfigEtccdConn()
				if err := json.Unmarshal([]byte(v), &c); err != nil {
					dialog.NewError(err, w).Show()
					return
				}

				if err := conf.ModConn(k, c); err != nil {
					dialog.NewError(err, w).Show()
					return
				}

				(*pGuiSelConnection).OnChanged(k)
			}, w)

		d.Resize(w.Canvas().Size())
		d.Show()
	})
	*pGuiButCreate = widget.NewButtonWithIcon(STR_EMPTY, theme.ContentAddIcon(), func() {
		binK := binding.NewString()
		guiEntK := widget.NewEntry()
		guiEntK.Bind(binK)
		guiEntK.Validator = func(s string) error {
			if s == "" {
				return ErrorConnectionNameEmpty
			}
			if _, ok := getConf().MapConn.Get(s); ok {
				return ErrorDuplicateEtcdConnectionName
			}
			return nil
		}

		binV := binding.NewString()
		binV.Set(GUI_SETTING_CREATE_CONN_PLACEHOLDER)
		guiEntV := widget.NewMultiLineEntry()
		guiEntV.Bind(binV)
		guiEntV.SetPlaceHolder("json format")
		guiEntV.SetMinRowsVisible(GUI_SETTING_EDIT_CONN_ENTRY_LINE_NUM)
		guiEntV.Validator = func(s string) error {
			if json.Valid([]byte(s)) {
				return nil
			}
			return ErrorInputNeedJsonFormat
		}
		d := dialog.NewForm("Add", "OK", "Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("ConnName", guiEntK),
				widget.NewFormItem("Conn", guiEntV),
			}, func(b bool) {
				if !b {
					return
				}
				k, _ := binK.Get()
				v, _ := binV.Get()
				if k == STR_EMPTY || v == STR_EMPTY {
					dialog.NewError(ErrorInputDataEmpty, w).Show()
					return
				}
				if _, ok := conf.MapConn.Get(k); ok {
					dialog.NewError(ErrorConnNameAlreadyExist, w).Show()
					return
				}

				c := petcd.NewConfigEtccdConn()
				if err := json.Unmarshal([]byte(v), &c); err != nil {
					dialog.NewError(err, w).Show()
					return
				}

				conf.AddConn(k, c)

				(*pGuiSelConnection).Options = conf.MapConn.M_sliKey
				(*pGuiSelConnection).Refresh()
				(*pGuiSelConnection).SetSelected(k)
			}, w)

		d.Resize(w.Canvas().Size())
		d.Show()
	})
	*pGuiButDelete = widget.NewButtonWithIcon(STR_EMPTY, theme.DeleteIcon(), func() {
		k := (*pGuiSelConnection).Selected
		_, ok := conf.MapConn.Get(k)
		if !ok {
			dialog.NewError(ErrorPathHasNoData, w).Show()
			return
		}
		d := dialog.NewConfirm("Delete", "Delete ConnName:"+k, func(b bool) {
			if !b {
				return
			}

			if err := conf.DelConn(k); err != nil {
				dialog.NewError(err, w).Show()
				return
			}

			(*pGuiSelConnection).Options = conf.MapConn.M_sliKey
			(*pGuiSelConnection).Refresh()
			(*pGuiSelConnection).ClearSelected()
			(*pGuiLabConnection).SetText("")
		}, w)

		d.Show()
	})

	*pGuiGriHeader = container.NewGridWithColumns(2, *pGuiSelConnection, container.NewHBox(*pGuiButConnection, *pGuiButEdit, *pGuiButCreate, *pGuiButDelete))
	*pGuiBorConnection = container.NewBorder(*pGuiGriHeader, nil, nil, nil, *pGuiLabConnection)
}
