package pgit

import "errors"

var (
	ErrDuplicateOpenRepository = errors.New("duplicate open repository")
	ErrConfigNotFound          = errors.New("config not found")
	ErrRepositoryNotOpen       = errors.New("repository not open")

	NotErrRefFound = errors.New("not error, ref found")
)
