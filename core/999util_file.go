package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// 文件修改时间(纳秒)
func fileModTime(fpath string) int64 {
	stat, err := os.Stat(fpath)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return stat.ModTime().UnixNano()
}

// 获取文件的mime type
func fileType(name string) string {
	index := strings.LastIndex(name, ".")
	if index > -1 {
		suffix := name[index:]
		return mimes[suffix]
	} else {
		return ""
	}
}

// 判断文件是否存在(路径指向的是目录也会返回false)
func fileExist(pt string) bool {
	st, err := os.Stat(pt)
	return st != nil && !st.IsDir() && !os.IsNotExist(err)
}

// 路径拼接
func pathJoin(base, uri string) string {
	return path.Join(base, uri)
}

// 从路径中获取文件名
func fileName(uri string) string {
	return path.Base(uri)
}


func doScan(path string, scanAll bool, res *[]string)  {
	if !isDir(path) {
		log.Fatal(path, "is not directory.")
		return
	}
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to read path: %v, %v", path, err.Error()))
		return
	}
	for _, f := range fs {
		nextPath := filepath.Join(path, f.Name())
		if f.IsDir() {
			if scanAll {
				*res = append(*res, nextPath)
			}
			doScan(nextPath, scanAll, res)
			continue
		}
		*res = append(*res, nextPath)
	}
}

func doScanForInfo(path string, res *[]*FileInfo)  {
	if !isDir(path) {
		log.Fatal(path, "is not directory.")
		return
	}
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to read path: %v, %v", path, err.Error()))
		return
	}
	for _, f := range fs {
		nextPath := filepath.Join(path, f.Name())
		finfo := newFileInfo(nextPath, f)
		*res = append(*res, finfo)

		if f.IsDir() {
			doScanForInfo(nextPath, res)
		}
	}
}

// 判断路径指向的是否为目录
func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}


