package ppath

import (
	"path"
	"regexp"
	"sort"

	"github.com/caticat/go_game_server/phelp"
)

type PPath struct {
	m_mapFiles *phelp.PSortedMap[string, bool] // 所有文件 <路径, 是否是文件夹>
	m_pathBase string
}

func NewPPath(pathBase string) *PPath {
	p := &PPath{
		m_mapFiles: phelp.NewPSortedMap[string, bool](),
		m_pathBase: path.Join(path.Dir(pathBase), path.Base(pathBase)),
	}

	return p
}

func (p *PPath) Refresh() error {
	files := make(map[string]bool)
	if err := phelp.Ls(p.m_pathBase, files, phelp.PBinFlag_Recursive); err != nil {
		return err
	}

	files[p.m_pathBase] = true // 追加根目录
	p.m_mapFiles.InitByMap(files)
	sort.Strings(p.m_mapFiles.M_sliKey)

	return nil
}

func (p *PPath) KeysAll(filter string) []string {
	if filter == "" {
		return p.m_mapFiles.M_sliKey
	} else {
		sliRet := make([]string, 0, len(p.m_mapFiles.M_sliKey))
		for _, key := range p.m_mapFiles.M_sliKey {
			match, _ := regexp.MatchString(filter, key)
			if !match {
				continue
			}
			sliRet = append(sliRet, key)
		}

		return sliRet
	}
}
func (p *PPath) Keys(basePath string) []string {
	basePath = p.FixPath(basePath)
	if _, ok := p.m_mapFiles.Get(basePath); !ok {
		return []string{}
	}

	ret := make([]string, 0, 1)
	for _, v := range p.m_mapFiles.M_sliKey {
		if path.Dir(v) != basePath {
			continue
		}
		ret = append(ret, v)
	}

	return ret
}

func (p *PPath) Has(basePath string) bool {
	basePath = p.FixPath(basePath)
	_, ok := p.m_mapFiles.Get(basePath)
	return ok
}

func (p *PPath) IsDir(basePath string) (bool, bool) {
	basePath = p.FixPath(basePath)
	return p.m_mapFiles.Get(basePath)
}

func (p *PPath) FixPath(pa string) string {
	if (pa == "") || (pa == ".") {
		return p.m_pathBase
	} else {
		return pa
	}
}

func (p *PPath) Debug() *phelp.PSortedMap[string, bool] {
	return p.m_mapFiles
}
