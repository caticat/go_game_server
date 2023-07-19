package pgit

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type PGitConf struct {
	Repository string
	Local      string
	Auth       *http.BasicAuth   // 认证
	Author     *object.Signature // 作者
}

func newPGitConf(repository, local, username, password, authorName, authorEMail string) *PGitConf {
	c := &PGitConf{
		Repository: repository,
		Local:      local,
		Auth: &http.BasicAuth{
			Username: username,
			Password: password,
		},
		Author: &object.Signature{
			Name:  authorName,
			Email: authorEMail,
		},
	}

	if len(c.Auth.Username) == 0 {
		c.Auth.Username = STR_NIL
	}

	return c
}
