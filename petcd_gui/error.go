package main

import "errors"

var (
	ErrorNoPathSelect  = errors.New("No Path selected")
	ErrorEmptyPath     = errors.New("No Key Entered")
	ErrorBadPathPrefix = errors.New("Path should start with '/'")
	ErrorPathHasNoData = errors.New("Path has No Data")
)
