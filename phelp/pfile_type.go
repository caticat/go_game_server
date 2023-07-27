package phelp

import "strings"

// image formats and magic numbers
var g_mapImageMagicTable = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

func IsImage(incipit []byte) bool {
	for magic, _ := range g_mapImageMagicTable {
		if strings.HasPrefix(string(incipit), magic) {
			return true
		}
	}

	return false
}
