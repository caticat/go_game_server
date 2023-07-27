package pgit

import "errors"

var (
	ErrDuplicateOpenRepository       = errors.New("duplicate open repository")
	ErrConfigNotFound                = errors.New("config not found")
	ErrRepositoryNotOpen             = errors.New("repository not open")
	ErrRepositoryRemoteLocalConflict = errors.New("repository remote local conflict")

	NotErrRefFound = errors.New("not error, ref found")
	NotErrNumLimit = errors.New("not error, number limit")
)
