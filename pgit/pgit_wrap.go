package pgit

import (
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

func (t *PGit) getConf() *PGitConf         { return t.m_conf }
func (t *PGit) getRep() *git.Repository    { return t.m_rep }
func (t *PGit) setRep(rep *git.Repository) { t.m_rep = rep }

func (t *PGit) openIfNilClone() error {
	// 初始校验
	rep := t.getRep()
	if rep != nil {
		return ErrDuplicateOpenRepository
	}
	conf := t.getConf()
	if conf == nil {
		return ErrConfigNotFound
	}

	// 打开本地仓库
	rep, err := git.PlainOpen(conf.Local)
	if err == nil {
		// 本地远程仓库地址不同的校验
		sliRemote, err := rep.Remotes()
		if err != nil {
			return err
		}
		isSameRepo := false
		for _, remote := range sliRemote {
			if strings.Contains(remote.String(), conf.Repository) {
				isSameRepo = true
				break
			}
		}
		if !isSameRepo {
			return ErrRepositoryRemoteLocalConflict
		}
		t.setRep(rep)
		return nil
	}
	if (err != git.ErrRepositoryNotExists) && (err != transport.ErrEmptyRemoteRepository) {
		return err
	}

	// 本地仓库没有,拉取远程仓库
	rep, err = git.PlainClone(conf.Local, false, &git.CloneOptions{
		URL: conf.Repository,
		// Progress: os.Stdout,
		Auth: conf.Auth,
	})
	if err != nil {
		return err
	} else {
		t.setRep(rep)
		return nil
	}
}

func (p *PGit) push(force bool) error {
	rep := p.getRep()
	if rep == nil {
		return ErrRepositoryNotOpen
	}
	conf := p.getConf()
	if conf == nil {
		return ErrConfigNotFound
	}

	err := rep.Push(&git.PushOptions{
		Force: force,
		Auth:  conf.Auth,
	})
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}
