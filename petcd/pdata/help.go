package pdata

func IsCurPath(path string) bool {
	return ((path == "") ||
		(path == "/"))
}

func PopPath(path string) (pathPrefix, pathLeft string) {
	pathPrefix = path
	pathLeft = ""
	if path == "" {
		return
	}

	sliPath := getRegPath().FindStringSubmatch(path)
	if len(sliPath) != 3 {
		return
	}

	pathPrefix = sliPath[1]
	pathLeft = sliPath[2]

	return
}
