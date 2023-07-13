package phelp

import "errors"

var (
	ErrorIndexOutofRange = errors.New("index out of range")
	ErrorNilSortFunction = errors.New("sort function is nil")
)
