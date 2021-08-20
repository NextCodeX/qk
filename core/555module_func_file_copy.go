package core

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

// 功能与linux命令 cp 基本一致
func (fns *InternalFunctionSet) Cp(src, dst string)  {
	srcLen := len(src)
	if srcLen < 1 || len(dst) < 1 {
		return
	}
	if strings.HasSuffix(src, "/*") {
		src = src[:srcLen-2]
		srcLen = len(src)
		var finfos []*FileInfo
		doScanForInfo(src,  &finfos)
		for _, finfo := range finfos {
			srcRelativePath := finfo.absolutePath[len(src):]
			dstPath := pathJoin(dst, srcRelativePath)
			if finfo.IsDir() {
				err := os.MkdirAll(dstPath, os.ModePerm)
				if err != nil {
					panic(err)
				}
			} else {
				copyFile(finfo.absolutePath, dstPath)
			}
		}
		return
	}

	srcStat, err := os.Stat(src)
	if err != nil {
		fmt.Println(err)
		return
	}
	dstStat, err := os.Stat(dst)
	if err != nil {
		fmt.Println(err)
		return
	}
	if srcStat.IsDir() && !dstStat.IsDir() {
		runtimeExcption("error: src is directory, dst is file")
	}
}

func copyFile(src string, dst string) {
	fsrc, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	fdst, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(fdst, fsrc)
	if err != nil {
		panic(err)
	}
}

type FileInfo struct {
	absolutePath string
	info fs.FileInfo
}

func newFileInfo(absPath string, info fs.FileInfo) *FileInfo {
	return &FileInfo{absolutePath: absPath, info: info}
}

func (fi *FileInfo) IsDir() bool {
	return fi.info.IsDir()
}

func (fi *FileInfo) Name() string {
	return fi.info.Name()
}

func (fi *FileInfo) AbsolutePath() string {
	return fi.absolutePath
}