package pgit

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type PGit struct {
	// 配置参数
	m_conf *PGitConf

	// 逻辑参数
	m_rep *git.Repository
}

// repository 仓库地址
// local 本地目录
// username 用户名,可以为空
// password 密码或者token
// authorName 作者名
// authorEMail 作者邮箱
func NewPGit(repository, local, username, password, authorName, authorEMail string) (*PGit, error) {
	p := &PGit{
		m_conf: newPGitConf(repository, local, username, password, authorName, authorEMail),
	}

	if err := p.openIfNilClone(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *PGit) Commit(comment string) error {
	rep := p.getRep()
	if rep == nil {
		return ErrRepositoryNotOpen
	}
	conf := p.getConf()
	if conf == nil {
		return ErrConfigNotFound
	}

	wor, err := rep.Worktree()
	if err != nil {
		return err
	}

	err = wor.AddWithOptions(&git.AddOptions{
		All: true,
	})
	if err != nil {
		return err
	}

	author := conf.Author
	author.When = time.Now()
	_, err = wor.Commit(comment, &git.CommitOptions{
		Author: author,
	})

	if err == git.ErrEmptyCommit {
		return nil
	}

	return err
}

func (p *PGit) Pull() error {
	rep := p.getRep()
	if rep == nil {
		return ErrRepositoryNotOpen
	}
	conf := p.getConf()
	if conf == nil {
		return ErrConfigNotFound
	}

	wor, err := rep.Worktree()
	if err != nil {
		return err
	}

	err = wor.Pull(&git.PullOptions{
		Auth: conf.Auth,
	})
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}

func (p *PGit) Push() error {
	return p.push(false)
}

func (p *PGit) PushForce() error {
	return p.push(true)
}

func (p *PGit) ResetToRemote() error {
	rep := p.getRep()
	if rep == nil {
		return ErrRepositoryNotOpen
	}
	conf := p.getConf()
	if conf == nil {
		return ErrConfigNotFound
	}

	err := rep.Fetch(&git.FetchOptions{})
	if err != git.NoErrAlreadyUpToDate {
		return err
	}

	refs, err := rep.References()
	if err != nil {
		return err
	}
	var remoteHash plumbing.Hash
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			if ref.Name().String() == REF_REMOTE {
				remoteHash = ref.Hash()
				return NotErrRefFound
			}
		}
		return nil
	})
	if err != NotErrRefFound {
		return err
	}

	wor, err := rep.Worktree()
	if err != nil {
		return err
	}

	return wor.Reset(&git.ResetOptions{
		Commit: remoteHash,
		Mode:   git.HardReset,
	})
}

// 先拉取后推送
func (p *PGit) Sync() error {
	if err := p.Pull(); err != nil {
		return err
	}

	return p.Push()
}
