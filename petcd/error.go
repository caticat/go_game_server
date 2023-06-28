package petcd

import "errors"

var (
	ErrorNilClient        = errors.New("client is nil")
	ErrorNilConfig        = errors.New("config is nil")
	ErrorInvalidOperation = errors.New("invalid operation")
)
