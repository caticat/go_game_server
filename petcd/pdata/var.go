package pdata

import "regexp"

var (
	g_regPath = regexp.MustCompile("(^/[0-9A-Za-z_]+?)(/.*)")
)

func getRegPath() *regexp.Regexp { return g_regPath }
