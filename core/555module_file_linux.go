package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// 与Linux命令：mkdir, ls, rm, mv, cp提供相似的功能
// 创建目录
func (fns *InternalFunctionSet) Mkdir(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("mkdir()", err)
	}
}

// 当前目录文件查看
func (fns *InternalFunctionSet) Ls(dir string) []interface{} {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("failed to read path: %v, %v\n", dir, err.Error())
		return nil
	}
	var res []interface{}
	for _, f := range fs {
		res = append(res, f.Name())
	}
	return res
}

// 删除文件或目录
func (fns *InternalFunctionSet) Rm(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 移动目录，文件； 或重命名目录，文件
func (fns *InternalFunctionSet) Mv(raw, target string) {
	err := os.Rename(raw, target)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 功能与linux命令 cp 基本一致
func (fns *InternalFunctionSet) Cp(src, dst string) {
	srcLen := len(src)
	if srcLen < 1 || len(dst) < 1 {
		return
	}
	if strings.HasSuffix(src, "/*") {
		src = src[:srcLen-2]
		srcLen = len(src)
		var finfos []*FileInfo
		doScanForInfo(src, &finfos)
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
	if !srcStat.IsDir() && !fileExist(dst) {
		copyFile(src, dst)
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
