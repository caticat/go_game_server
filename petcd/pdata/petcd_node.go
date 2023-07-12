package pdata

import (
	"path"
	"sort"
)

type PEtcdNode struct {
	m_key         string
	m_value       string
	m_mapChild    map[string]*PEtcdNode
	m_sliChildKey []string // 保证map遍历顺序,防止界面刷新出现闪动问题
}

func newPEtcdNode(prefix string) *PEtcdNode {
	return &PEtcdNode{
		m_key:         prefix,
		m_value:       "",
		m_mapChild:    make(map[string]*PEtcdNode),
		m_sliChildKey: nil,
	}
}

func (n *PEtcdNode) GetKey() string   { return n.m_key }
func (n *PEtcdNode) GetValue() string { return n.m_value }
func (n *PEtcdNode) IsBranch() bool   { return len(n.m_mapChild) > 0 }

func (n *PEtcdNode) ChildKeys() []string {
	keys := make([]string, 0, len(n.m_mapChild))
	for _, k := range n.m_sliChildKey {
		keys = append(keys, n.m_mapChild[k].GetKey())
	}

	return keys
}

func (n *PEtcdNode) Clear() {
	// 这里不会删除m_key,因为key的值是不变的

	n.m_value = ""
	n.m_sliChildKey = make([]string, 0)

	for k := range n.m_mapChild {
		delete(n.m_mapChild, k)
	}
}

func (t *PEtcdNode) get(key string) *PEtcdNode {
	pathPrefix, pathLeft := PopPath(key)
	if pathPrefix == "" {
		return t
	} else {
		if child := t.m_mapChild[pathPrefix]; child != nil {
			if IsCurPath(pathLeft) {
				return child
			} else {
				return child.get(pathLeft)
			}
		} else {
			return nil
		}
	}
}

func (t *PEtcdNode) set(value string) { t.m_value = value }

func (t *PEtcdNode) create(prefix, key string) *PEtcdNode {
	pathPrefix, pathLeft := PopPath(key)
	if pathPrefix == "" {
		return t
	} else {
		prefix = path.Join(prefix, pathPrefix)
		child := t.m_mapChild[pathPrefix]
		if child == nil {
			child = newPEtcdNode(prefix)
			t.m_mapChild[pathPrefix] = child
			t.m_sliChildKey = append(t.m_sliChildKey, pathPrefix)
			sort.Strings(t.m_sliChildKey)
		}
		if IsCurPath(pathLeft) {
			return child
		} else {
			return child.create(prefix, pathLeft)
		}
	}
}
