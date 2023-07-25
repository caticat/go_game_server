package main

import (
	"encoding/json"
	"fmt"
	"io"
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
		guiBorImport     *fyne.Container    = nil // 布局 导入
		guiBorExport     *fyne.Container    = nil // 布局 导出
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

	// import
	initGUISettingImport(w, &guiBorImport)

	// export
	initGUISettingExport(w, &guiBorExport)

	// 页签
	guiConTab = container.NewAppTabs(
		container.NewTabItem("Connection", guiBorConnection),
		container.NewTabItem("Base", guiForBase),
		container.NewTabItem("Import", guiBorImport),
		container.NewTabItem("Export", guiBorExport),
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
			plog.ErrorLn(ErrorPathHasNoData)
			return
		}
		if err := conf.SetSelect(k); err != nil {
			dialog.NewError(err, w).Show()
			plog.ErrorLn(err)
			return
		}
		(*pGuiButConnection).Disable()
		msg := fmt.Sprintf("%s now is Connecting to: %q", WINDOW_TITLE, k)
		dialog.NewInformation("Information", msg, w).Show()
		plog.InfoLn(msg)
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
			plog.ErrorLn(ErrorPathHasNoData)
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
					plog.ErrorLn(err)
					return
				}

				if err := conf.ModConn(k, c); err != nil {
					dialog.NewError(err, w).Show()
					plog.ErrorLn(err)
					return
				}

				(*pGuiSelConnection).OnChanged(k)
				plog.InfoF("mod conn:%q\n", k)
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
					plog.ErrorLn(ErrorInputDataEmpty)
					return
				}
				if _, ok := conf.MapConn.Get(k); ok {
					dialog.NewError(ErrorConnNameAlreadyExist, w).Show()
					plog.ErrorLn(ErrorConnNameAlreadyExist)
					return
				}

				c := petcd.NewConfigEtccdConn()
				if err := json.Unmarshal([]byte(v), &c); err != nil {
					dialog.NewError(err, w).Show()
					plog.ErrorLn(err)
					return
				}

				conf.AddConn(k, c)

				(*pGuiSelConnection).Options = conf.MapConn.M_sliKey
				(*pGuiSelConnection).Refresh()
				(*pGuiSelConnection).SetSelected(k)
				plog.InfoF("add conn:%q\n", k)
			}, w)

		d.Resize(w.Canvas().Size())
		d.Show()
	})
	*pGuiButDelete = widget.NewButtonWithIcon(STR_EMPTY, theme.DeleteIcon(), func() {
		k := (*pGuiSelConnection).Selected
		_, ok := conf.MapConn.Get(k)
		if !ok {
			dialog.NewError(ErrorPathHasNoData, w).Show()
			plog.ErrorLn(ErrorPathHasNoData)
			return
		}
		d := dialog.NewConfirm("Delete", "Delete ConnName:"+k, func(b bool) {
			if !b {
				return
			}

			if err := conf.DelConn(k); err != nil {
				dialog.NewError(err, w).Show()
				plog.ErrorLn(err)
				return
			}

			(*pGuiSelConnection).Options = conf.MapConn.M_sliKey
			(*pGuiSelConnection).Refresh()
			(*pGuiSelConnection).ClearSelected()
			(*pGuiLabConnection).SetText("")
			plog.InfoF("del conn:%q\n", k)
		}, w)

		d.Show()
	})

	*pGuiGriHeader = container.NewGridWithColumns(2, *pGuiSelConnection, container.NewHBox(*pGuiButConnection, *pGuiButEdit, *pGuiButCreate, *pGuiButDelete))
	*pGuiBorConnection = container.NewBorder(*pGuiGriHeader, nil, nil, nil, *pGuiLabConnection)
}

func initGUISettingImport(w fyne.Window, pGuiBorImport **fyne.Container) {
	// 导入数据输入框内容
	tmpCB := petcd.NewConfigEtcdInitBase()
	tmpCB.EtcdKVList = append(tmpCB.EtcdKVList, petcd.NewEtcdKV("/exampleKey", "exampleValue"))
	binJsonData := binding.NewString()
	binJsonData.Set(phelp.ToJsonIndent(tmpCB))
	guiEntJsonData := widget.NewMultiLineEntry()
	guiEntJsonData.Bind(binJsonData)
	guiEntJsonData.SetMinRowsVisible(GUI_SETTING_EDIT_IMPORT_ENTRY_LINE_NUM)
	guiEntJsonData.Validator = func(s string) error {
		if json.Valid([]byte(s)) {
			return nil
		}
		return ErrorInputNeedJsonFormat
	}

	// 导入数据触发
	funImportJsonData := func() {
		jsonData, err := binJsonData.Get()
		if err != nil {
			dialog.NewError(err, w).Show()
			plog.ErrorLn(err)
			return
		}
		cb := petcd.NewConfigEtcdInitBase()
		if err = json.Unmarshal([]byte(jsonData), cb); err != nil {
			dialog.NewError(err, w).Show()
			plog.ErrorLn(err)
			return
		}
		if !getConnected() {
			dialog.NewError(ErrorEtcdNotConnected, w).Show()
			plog.ErrorLn(ErrorEtcdNotConnected)
			return
		}
		c := petcd.NewConfigEtcdInit()
		c.SetConfigEtcd(getConf().GetCfgETCD())
		c.SetBase(cb)
		diaCon := dialog.NewCustomConfirm("Run Init Etcd Data?", "Confirm", "Cancel", container.NewScroll(widget.NewLabel(phelp.ToJsonIndent(c))), func(b bool) {
			if !b {
				return
			}
			if err = petcd.ProcessEtcdInitWithoutConn(c); err != nil {
				dialog.NewError(err, w).Show()
				plog.ErrorLn(err)
				return
			}
			dialog.NewInformation("Infomation", "Init Ectd Data Done", w).Show()
			getFunGUIHomeRefresh()()
			plog.InfoLn("import ectd data done")
		}, w)
		// dialog.NewConfirm("Run Init Etcd Data?", phelp.ToJsonIndent(c), )
		diaCon.Resize(w.Canvas().Size())
		diaCon.Show()
	}

	// 界面
	*pGuiBorImport = container.NewMax(container.NewAppTabs(
		container.NewTabItem("From File", container.NewVBox(widget.NewButtonWithIcon("Open File(Need Json Format)", theme.FolderOpenIcon(), func() {
			diaFile := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
				if err != nil {
					plog.ErrorLn(err)
					return
				}
				if f == nil { // 主动关闭
					return
				}
				defer f.Close()
				bs, err := io.ReadAll(f)
				if err != nil {
					plog.ErrorLn(err)
					return
				}
				binJsonData.Set(string(bs))
				funImportJsonData()
			}, w)
			diaFile.Resize(w.Canvas().Size())
			diaFile.Show()
		}))),
		container.NewTabItem("From Editor", container.NewScroll(
			container.NewVBox(
				widget.NewLabel("Json Data:"),
				guiEntJsonData,
				widget.NewButtonWithIcon("Submit", theme.ConfirmIcon(), func() {
					funImportJsonData()
				}),
			),
		)),
	))
}

func initGUISettingExport(w fyne.Window, pGuiBorExport **fyne.Container) {
	*pGuiBorExport = container.NewVBox(widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), func() {
		diaFile := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
			if uc == nil {
				return
			}
			cb := petcd.NewConfigEtcdInitBase()
			root := getEtcdData()
			ks := root.AllKeys()
			cb.EtcdKVList = make([]*petcd.EtcdKV, 0, len(ks))
			for _, k := range ks {
				v, ok := root.GetValue(k)
				if !ok {
					plog.ErrorLn("key:", k, " has not found")
					continue
				}
				cb.EtcdKVList = append(cb.EtcdKVList, petcd.NewEtcdKV(k, v))
			}
			jsonData := phelp.ToJsonIndent(cb)
			l := len(jsonData)
			for l > 0 {
				n, err := io.WriteString(uc, jsonData)
				if err != nil {
					dialog.NewError(err, w).Show()
					plog.ErrorLn(err)
					break
				}
				l -= n
			}
			if l <= 0 {
				msg := fmt.Sprintf("export etcd eata to %q done", uc.URI())
				dialog.NewInformation("Export", msg, w).Show()
				plog.InfoLn(msg)
			}
		}, w)
		diaFile.Resize(w.Canvas().Size())
		diaFile.Show()
	}))
}
