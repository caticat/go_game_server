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

	// conn
	lep := widget.NewLabel("default text")
	var search *widget.Select
	var butConnect *widget.Button
	butConnect = widget.NewButtonWithIcon("Connect", theme.LoginIcon(), func() { // 连接按钮前移,解决sarch的空指针问题
		k := search.Selected
		_, ok := conf.MapConn.Get(k)
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
	search = widget.NewSelect(conf.MapConn.M_sliKey, func(s string) {
		c, ok := conf.MapConn.Get(s)
		if !ok {
			return
		}
		lep.SetText(phelp.ToJsonIndent(c))
		if conf.ConnSelect == s {
			butConnect.Disable()
		} else {
			butConnect.Enable()
		}
	})
	search.SetSelected(conf.ConnSelect)
	butEdit := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		k := search.Selected
		v, ok := conf.MapConn.Get(k)
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

				c := petcd.NewConfigEtccdConn()
				if err := json.Unmarshal([]byte(v), &c); err != nil {
					dialog.NewError(err, w).Show()
					return
				}

				if err := conf.ModConn(k, c); err != nil {
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
				widget.NewFormItem("ConnName", ek),
				widget.NewFormItem("Conn", en),
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

				search.Options = conf.MapConn.M_sliKey
				search.Refresh()
				search.SetSelected(k)
			}, w)

		di.Resize(w.Canvas().Size())
		di.Show()
	})
	butDelete := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		k := search.Selected
		_, ok := conf.MapConn.Get(k)
		if !ok {
			dialog.NewError(ErrorPathHasNoData, w).Show()
			return
		}
		di := dialog.NewConfirm("Delete", "Delete ConnName:"+k, func(b bool) {
			if !b {
				return
			}

			if err := conf.DelConn(k); err != nil {
				dialog.NewError(err, w).Show()
				return
			}

			search.Options = conf.MapConn.M_sliKey
			search.Refresh()
			search.ClearSelected()
		}, w)

		di.Show()
	})

	header := container.NewGridWithColumns(2, search, container.NewHBox(butConnect, butEdit, butCreate, butDelete))
	cep := container.NewBorder(header, nil, nil, nil, lep)

	con := container.NewAppTabs(
		container.NewTabItem("Connection", cep),
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
