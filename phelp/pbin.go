package phelp

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type PBinFlags int

const (
	PBinFlag_None      PBinFlags = 0
	PBinFlag_All       PBinFlags = 1
	PBinFlag_Recursive           = 1 << 1
)

// 获取文件夹下所有内容
func Ls(dir string, files map[string]bool, flags PBinFlags) error {
	// 选项参数
	all := (flags & PBinFlag_All) > 0
	recursive := (flags & PBinFlag_Recursive) > 0

	// 逻辑
	rFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, rFile := range rFiles {
		name := rFile.Name()
		if !all { // 隐藏文件/文件夹
			if strings.HasPrefix(name, PATH_CURRENT) {
				continue
			}
		}
		name = path.Join(dir, name)
		isDir := rFile.IsDir()
		files[name] = isDir
		if isDir && recursive {
			if err = Ls(name, files, flags); err != nil {
				return err
			}
		}
	}

	return nil
}

// 删除文件/文件夹下所有内容 递归
func Rm(dir string) error {
	return os.RemoveAll(dir)
	// rFiles, err := ioutil.ReadDir(dir)
	// if err != nil {
	// 	return err
	// }

	// for _, rFile := range rFiles {
	// 	name := path.Join(dir, rFile.Name())
	// 	if rFile.IsDir() {
	// 		return Rm(name)
	// 	} else {
	// 		os.Remove(name)
	// 	}
	// }

	// return nil
}
