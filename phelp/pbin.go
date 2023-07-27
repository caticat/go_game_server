package phelp

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type PBinFlags int

const (
	PBinFlag_None      PBinFlags = 0
	PBinFlag_All       PBinFlags = 1
	PBinFlag_Recursive           = 1 << 1
	PBinFlag_Force               = 1 << 2
)

// 获取文件夹下所有内容
func Ls(dir string, files map[string]bool, flags PBinFlags) error {
	dir = Format(dir)

	// 文件文件夹判断
	staFrom, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !staFrom.IsDir() { // 文件,直接返回
		files[dir] = false
		return nil
	}

	// 选项参数
	all := (flags & PBinFlag_All) > 0
	recursive := (flags & PBinFlag_Recursive) > 0

	// 逻辑
	rFiles, err := os.ReadDir(dir)
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
	dir = Format(dir)
	if dir == PATH_SEPARATOR {
		return ErrorCannotRmRoot
	}
	return os.RemoveAll(dir)
}

// 复制文件,返回成功复制的文件名列表(不包含文件夹)
func Cp(from, to string, flags PBinFlags) ([]string, error) {
	from = Format(from)
	to = Format(to)

	files := make(map[string]bool)
	if err := Ls(from, files, flags); err != nil {
		return nil, err
	}

	isOneFile := false // 只是一个文件的复制(特殊情况,直接复制内容到to文件,而不是to文件夹)
	if len(files) == 1 {
		if staFrom, err := os.Stat(from); err != nil {
			return nil, err
		} else {
			if !staFrom.IsDir() {
				if staTo, err := os.Stat(to); err != nil {
					isOneFile = true
				} else {
					if !staTo.IsDir() {
						isOneFile = true
					}
				}
			}
		}
	}

	force := (flags & PBinFlag_Force) > 0

	sliCpFiles := make([]string, 0, len(files))
	for fileFrom, isDir := range files {
		if isDir { // 空文件夹不处理
			continue
		}

		fileNameRelative, _ := strings.CutPrefix(fileFrom, from)
		if fileNameRelative == "" { // from本身是文件的情况
			fileNameRelative = path.Base(from)
		}
		fileTo := to
		if !isOneFile {
			fileTo = path.Join(to, fileNameRelative)
		}
		staTo, err := os.Stat(fileTo)
		if err == nil { // 文件已存在
			if !force { // 不覆盖
				continue
			}
			if staTo.IsDir() != isDir {
				return sliCpFiles, fmt.Errorf("from:%q to:%q fileType mismatch", fileFrom, fileTo)
			}
		}

		if err := os.MkdirAll(path.Dir(fileTo), 0655); err != nil { // 创建目标文件对应文件夹
			return sliCpFiles, err
		}
		fileFromOpen, err := os.Open(fileFrom)
		if err != nil {
			return sliCpFiles, err
		}
		defer fileFromOpen.Close()
		fileToOpen, err := os.OpenFile(fileTo, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return sliCpFiles, err
		}
		defer fileToOpen.Close()
		_, err = io.Copy(fileToOpen, fileFromOpen)
		if err != nil {
			return sliCpFiles, err
		}
		sliCpFiles = append(sliCpFiles, fileFrom)
	}

	return sliCpFiles, nil
}

func Mv(from, to string, flags PBinFlags) error {
	from = Format(from)
	to = Format(to)

	if sliCpFiles, err := Cp(from, to, flags); err != nil {
		return err
	} else {
		for _, file := range sliCpFiles {
			if err := Rm(file); err != nil {
				return err
			}
		}
	}

	return nil
}

func Format(pathFrom string) string {
	return strings.ReplaceAll(pathFrom, "\\", "/")
}
