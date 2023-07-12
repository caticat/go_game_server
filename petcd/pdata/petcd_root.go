package pdata

import (
	"strings"

	"github.com/caticat/go_game_server/phelp"
	"github.com/caticat/go_game_server/plog"
)

type PEtcdRoot struct {
	*PEtcdNode
	m_mapData map[string]string
}

func NewPEtcdRoot() *PEtcdRoot {
	return &PEtcdRoot{
		PEtcdNode: newPEtcdNode(PDATA_PREFIX),
		m_mapData: make(map[string]string),
	}
}

func (t *PEtcdRoot) GetNode() *PEtcdNode { return t.PEtcdNode }

func (t *PEtcdRoot) Get(key string) *PEtcdNode {
	if key == "" {
		key = PDATA_PREFIX
	}
	if !strings.HasPrefix(key, PDATA_PREFIX) {
		plog.ErrorLn(ErrorInvalidPath)
		return nil
	}

	node := t.GetNode()
	if IsCurPath(key) {
		return node
	} else {
		return node.get(key)
	}
}

func (t *PEtcdRoot) GetValue(key string) (value string, ok bool) {
	value, ok = t.m_mapData[key]
	return
}

func (t *PEtcdRoot) Set(key, value string) {
	if !strings.HasPrefix(key, PDATA_PREFIX) {
		plog.ErrorLn(ErrorInvalidPath)
		return
	}
	node := t.Get(key)
	if node == nil {
		node = t.createNode(key)
		if node == nil {
			return
		}
	}
	node.set(value)
}

func (t *PEtcdRoot) SetAll(mapData map[string]string) {
	t.m_mapData = mapData
	for k, v := range mapData {
		t.Set(k, v)
	}
}

func (t *PEtcdRoot) AllKeys() []string {
	return phelp.Keys(t.m_mapData)
}

func (t *PEtcdRoot) createNode(key string) *PEtcdNode {
	if !strings.HasPrefix(key, PDATA_PREFIX) {
		plog.ErrorLn(ErrorInvalidPath)
		return nil
	}

	node := t.GetNode()
	if IsCurPath(key) {
		return node
	} else {
		return node.create(PDATA_PREFIX, key)
	}
}
