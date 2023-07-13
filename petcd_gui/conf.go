package main

import (
	"encoding/json"
	"flag"
	"time"

	"github.com/caticat/go_game_server/petcd"
	"github.com/caticat/go_game_server/phelp"
	"github.com/caticat/go_game_server/plog"
)

type ConfigEtcdGUI struct {
	Base           *petcd.ConfigEtcdBase               `json:"base"`            // 连接基础信息
	EndPointSelect string                              `json:"endpoint-select"` // 当前选择的连接节点
	MapEndPoint    *phelp.PSortedMap[string, []string] `json:"endpoint-list"`   // 所有连接节点配置 <连接名, 配置>
	LogLevel       int                                 `json:"log-level"`       // 日志等级
}

func NewConfigEtcdGUI() *ConfigEtcdGUI {
	return &ConfigEtcdGUI{
		Base:           petcd.NewConfigEtcdBase(),
		EndPointSelect: PETCD_CFG_ENDPOINT_SELECT_DEFAULT,
		MapEndPoint:    phelp.NewPSortedMap[string, []string](),
		LogLevel:       int(plog.ELogLevel_Debug),
	}
}

func (t *ConfigEtcdGUI) Init() error {
	// 程序记录数据清理
	reset := flag.Bool("reset", false, "clean up app local data")
	flag.Parse()
	if *reset {
		t.clearPreferences()
		return ErrorAppResetDone
	}

	// 连接基础信息
	a := getApp()
	cfgBase := a.Preferences().StringWithFallback(PETCD_CFG_BASE, PETCD_CFG_BASE_DEFAULT)
	if err := json.Unmarshal([]byte(cfgBase), t.Base); err != nil {
		return err
	}
	t.Base.EnableReInit = true // GUI工具这里强制可重复开启

	// 连接配置
	cfgEndPointList := a.Preferences().StringWithFallback(PETCD_CFG_ENDPOINT_LIST, PETCD_CFG_ENDPOINT_LIST_DEFAULT)
	tm := make(map[string][]string)
	if err := json.Unmarshal([]byte(cfgEndPointList), &tm); err != nil {
		return err
	}
	t.MapEndPoint.InitByMap(tm)

	// 当前选择的连接节点
	t.EndPointSelect = a.Preferences().StringWithFallback(PETCD_CFG_ENDPOINT_SELECT, PETCD_CFG_ENDPOINT_SELECT_DEFAULT)
	if _, ok := t.MapEndPoint.Get(t.EndPointSelect); !ok {
		if t.MapEndPoint.Length() > 0 {
			if k, _, err := t.MapEndPoint.GetByIndex(0); err == nil {
				t.EndPointSelect = k
			}
		}
	}

	// 日志等级
	logLevel := a.Preferences().IntWithFallback(PETCD_CFG_LOG_LEVEL, int(PETCD_CFG_LOG_LEVEL_DEFAULT))
	t.SetLogLevel(plog.ELogLevel(logLevel))

	return nil
}

func (t *ConfigEtcdGUI) GetCfgETCD() *petcd.ConfigEtcd {
	c := petcd.NewConfigEtcd()

	// base
	c.SetBase(t.Base)

	// endpoint
	if v, ok := t.MapEndPoint.Get(t.EndPointSelect); ok {
		c.SetEndpoints(v)
	}

	return c
}

func (t *ConfigEtcdGUI) SetBase(c *petcd.ConfigEtcdBase) {
	t.Base = c

	s, err := toJsonIndent(t.Base)
	if err != nil {
		plog.ErrorLn(err)
		return
	}

	a := getApp()
	a.Preferences().SetString(PETCD_CFG_BASE, s)
}

func (t *ConfigEtcdGUI) SetSelect(s string) error {
	if _, ok := t.MapEndPoint.Get(s); !ok {
		return ErrorSelectEtcdEndPointNotFound
	}

	if t.EndPointSelect != s {
		t.EndPointSelect = s
		a := getApp()
		a.Preferences().SetString(PETCD_CFG_ENDPOINT_SELECT, string(s))
	}

	return t.reconnect()
}

func (t *ConfigEtcdGUI) AddEndPoint(n string, c []string) error {
	if _, ok := t.MapEndPoint.Get(n); ok {
		return ErrorDuplicateEtcdEndPointName
	}

	t.MapEndPoint.Set(n, c)

	return t.saveEndPoint()
}

func (t *ConfigEtcdGUI) ModEndPoint(n string, c []string) error {
	_, ok := t.MapEndPoint.Get(n)
	if !ok {
		return ErrorPathHasNoData
	}

	t.MapEndPoint.Set(n, c)

	if err := t.saveEndPoint(); err != nil {
		return err
	}

	if t.EndPointSelect == n {
		return t.reconnect()
	} else {
		return nil
	}
}

func (t *ConfigEtcdGUI) DelEndPoint(n string) error {
	if _, ok := t.MapEndPoint.Get(n); !ok {
		return ErrorPathHasNoData
	}

	if t.EndPointSelect == n { // 不能删除当前连接
		return ErrorDeleteSelectingEndPoint
	}

	t.MapEndPoint.Del(n)

	if err := t.saveEndPoint(); err != nil {
		return err
	}

	return nil
}

func (t *ConfigEtcdGUI) saveEndPoint() error {
	a := getApp()
	if s, err := toJsonIndent(t.MapEndPoint.GetMap()); err == nil {
		a.Preferences().SetString(PETCD_CFG_ENDPOINT_LIST, s)
		return nil
	} else {
		return err
	}
}

func (t *ConfigEtcdGUI) reconnect() error {
	// 关闭当前连接
	setConnected(false)
	getFunUpdateTitle()()
	close()
	time.Sleep(time.Second)

	// 重新连接
	if err := petcd.Init(t.GetCfgETCD()); err != nil {
		return err
	}

	// 刷新界面
	setConnected(true)
	getFunUpdateTitle()()
	getFunGUIRefresh()()

	return nil
}

func (t *ConfigEtcdGUI) SetLogLevel(logLevel plog.ELogLevel) {
	if t.LogLevel == int(logLevel) {
		return
	}

	t.LogLevel = int(logLevel)
	plog.SetLogLevel(logLevel)

	a := getApp()
	a.Preferences().SetInt(PETCD_CFG_LOG_LEVEL, t.LogLevel)
}

func (t *ConfigEtcdGUI) clearPreferences() {
	a := getApp()

	a.Preferences().RemoveValue(PETCD_CFG_BASE)
	a.Preferences().RemoveValue(PETCD_CFG_ENDPOINT_SELECT)
	a.Preferences().RemoveValue(PETCD_CFG_ENDPOINT_LIST)
	a.Preferences().RemoveValue(PETCD_CFG_LOG_LEVEL)

	time.Sleep(time.Second) // 暂停一下保证数据保存完毕
}
