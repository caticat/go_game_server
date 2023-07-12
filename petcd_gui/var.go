package main

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/caticat/go_game_server/petcd/pdata"
)

var (
	g_conf      = NewConfigEtcdGUI()
	g_etcdData  = pdata.NewPEtcdRoot()
	g_etcdKey   = ""
	g_etcdValue = binding.NewString()
)

func getConf() *ConfigEtcdGUI       { return g_conf }
func getEtcdData() *pdata.PEtcdRoot { return g_etcdData }
func getEtcdKey() string            { return g_etcdKey }
func setEtcdKey(key string)         { g_etcdKey = key }
func getEtcdValue() binding.String  { return g_etcdValue }
