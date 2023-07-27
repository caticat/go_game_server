package phelp

import "errors"

var (
	ErrorIndexOutofRange = errors.New("index out of range")
	ErrorCannotRmRoot    = errors.New("can not rm /")
)
