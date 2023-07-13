package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/caticat/go_game_server/petcd/pdata"
)

var (
	g_conf           = NewConfigEtcdGUI()
	g_etcdData       = pdata.NewPEtcdRoot()
	g_etcdKey        = ""
	g_etcdValue      = binding.NewString()
	g_app            fyne.App
	g_funGUIRefresh  func()
	g_funUpdateTitle func()
	g_connected      bool = false
	g_logData             = binding.NewString()
)

func getConf() *ConfigEtcdGUI       { return g_conf }
func getEtcdData() *pdata.PEtcdRoot { return g_etcdData }
func getEtcdKey() string            { return g_etcdKey }
func setEtcdKey(key string)         { g_etcdKey = key }
func getEtcdValue() binding.String  { return g_etcdValue }
func getApp() fyne.App              { return g_app }
func getFunGUIRefresh() func()      { return g_funGUIRefresh }
func getFunUpdateTitle() func()     { return g_funUpdateTitle }
func setConnected(c bool)           { g_connected = c }
func getConnected() bool            { return g_connected }
func getLogData() binding.String    { return g_logData }
