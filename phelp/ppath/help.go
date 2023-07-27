package ppath

import (
	"path/filepath"
	"strings"
)

func IsSubDir(absBase, absSub string) bool {
	rel, err := filepath.Rel(absBase, absSub)
	if err != nil {
		return false
	}

	return (!strings.Contains(rel, ".."))
}
