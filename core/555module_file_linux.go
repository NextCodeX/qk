package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 与Linux命令：mkdir, ls, rm, mv, cp提供相似的功能
// 创建目录
func (this *InternalFunctionSet) Mkdir(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("mkdir() exception!", err)
	}
}

// 当前目录文件查看
func (this *InternalFunctionSet) Ls(dir string) Value {
	fs, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("failed to read path: %v, %v\n", dir, err.Error())
		return nil
	}
	var res []Value
	for _, f := range fs {
		fi := newFileInfo(dir, f)
		res = append(res, newClass("FileInfo", &fi))
	}
	return array(res)
}

// 删除文件或目录
func (this *InternalFunctionSet) Rm(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 移动目录，文件； 或重命名目录，文件
// 文件 -> 文件(重命名或移动)
// 目录 -> 目录
// 文件/目录列表 -> 目录
func (this *InternalFunctionSet) Mv(src, dst string) {
	srcFiles := findSrcFilesForMove(src)
	qty := len(srcFiles)
	if qty < 1 {
		return
	}
	if qty == 1 {
		srcFile := srcFiles[0]
		dstFile := calcMoveTargetFile(srcFile, dst)
		doFileMove(srcFile, dstFile)
		return
	}

	if fileExist(dst) {
		// 目标路径是一个已存在的文件(非目录), 直接返回
		fmt.Println(dst, " is the existed file, not directory!")
		return
	}
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		os.MkdirAll(dst, 0777)
	}

	for _, srcFile := range srcFiles {
		dstFile := calcMoveTargetFile(srcFile, dst)
		doFileMove(srcFile, dstFile)
	}
}

func findSrcFilesForMove(src string) []string {
	var res []string
	if _, err := os.Stat(src); os.IsNotExist(err) {
		if !strings.Contains(src, "*") {
			fmt.Println(src, "is not found!")
			return res
		}
		if _, err = os.Stat(filepath.Dir(src)); os.IsNotExist(err) {
			fmt.Println(src, "is not found!")
			return res
		}
	} else {
		res = append(res, src)
		return res
	}

	if strings.Contains(src, "*") {
		name := filepath.Base(src)
		symbolIndex := strings.Index(name, "*")

		allFlag := name == "*"
		firstFlag := symbolIndex == 0
		endFlag := symbolIndex == len(src)-1
		var prefix, suffix string
		if !allFlag {
			split := strings.Split(name, "*")
			if firstFlag {
				suffix = split[1]
			} else if endFlag {
				prefix = split[0]
			} else {
				prefix, suffix = split[0], split[1]
			}
		}

		dir := filepath.Dir(src)
		infos, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Println(err)
			return res
		}
		for _, info := range infos {
			fname := info.Name()
			if allFlag {
				res = append(res, filepath.Join(dir, fname))
				continue
			}

			if firstFlag && strings.HasSuffix(fname, suffix) ||
				endFlag && strings.HasPrefix(fname, prefix) ||
				strings.HasPrefix(fname, prefix) && strings.HasSuffix(fname, suffix) {
				res = append(res, filepath.Join(dir, fname))
			}
		}
	}
	return res
}

func doFileMove(src string, dst string) {
	err := os.Rename(src, dst)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func calcMoveTargetFile(srcFile string, dst string) string {
	if isDir(dst) {
		return filepath.Join(dst, filepath.Base(srcFile))
	} else if isDir(filepath.Dir(dst)) {
		return dst
	} else {
		parentDir := filepath.Dir(dst)
		os.MkdirAll(parentDir, 0777)
		return dst
	}
}

// 功能与linux命令 cp 基本一致
func (this *InternalFunctionSet) Cp(src, dst string) {
	srcFiles := findSrcFiles(src)
	qty := len(srcFiles)
	if qty < 1 {
		return
	}
	if qty == 1 {
		srcFile := srcFiles[0]
		_, dstFile := calcDstFile(src, dst, srcFile)
		os.MkdirAll(filepath.Dir(dstFile), 0777)
		copyFile(srcFile, dstFile)
		return
	}

	if fileExist(dst) {
		// 目标路径是一个已存在的文件(非目录), 直接返回
		fmt.Println(dst, " is the existed file, not directory!")
		return
	}
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		os.MkdirAll(dst, 0777)
	}

	for _, srcFile := range srcFiles {
		dstDir, dstFile := calcDstFile(src, dst, srcFile)
		os.MkdirAll(dstDir, 0777)
		copyFile(srcFile, dstFile)
	}
}

func calcDstFile(src string, dst string, srcFile string) (string, string) {
	if sta, err := os.Stat(dst); os.IsNotExist(err) || (sta != nil && !sta.IsDir()) {
		// 文件已存在(且非目录)或文件不存在
		return "", dst
	}

	srcDir := filepath.Dir(src)
	fname := srcFile[len(srcDir):]
	dstPath := filepath.Join(dst, fname)
	return filepath.Dir(dstPath), dstPath
}

func findSrcFiles(src string) []string {
	var res []string
	if _, err := os.Stat(src); os.IsNotExist(err) {
		if !strings.Contains(src, "*") {
			fmt.Println(src, "is not found!")
			return res
		}
		if _, err = os.Stat(filepath.Dir(src)); os.IsNotExist(err) {
			fmt.Println(src, "is not found!")
			return res
		}
	} else if !isDir(src) && !strings.Contains(src, "*") {
		res = append(res, src)
		return res
	} else {
	}

	if isDir(src) || filepath.Base(src) == "*" {
		dir := src
		if filepath.Base(src) == "*" {
			dir = filepath.Dir(src)
		}
		var tmp []*FileInfo
		doScanForInfo(dir, &tmp)
		for _, item := range tmp {
			if item.IsDir() {
				continue
			}
			res = append(res, item.Path())
		}
		return res
	}
	if strings.Contains(src, "*") {
		name := filepath.Base(src)
		symbolIndex := strings.Index(name, "*")

		firstFlag := symbolIndex == 0
		endFlag := symbolIndex == len(src)-1
		var prefix, suffix string
		split := strings.Split(name, "*")
		if firstFlag {
			suffix = split[1]
		} else if endFlag {
			prefix = split[0]
		} else {
			prefix, suffix = split[0], split[1]
		}

		dir := filepath.Dir(src)
		infos, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Println(err)
			return res
		}
		for _, info := range infos {
			fname := info.Name()
			if firstFlag && strings.HasSuffix(fname, suffix) ||
				endFlag && strings.HasPrefix(fname, prefix) ||
				strings.HasPrefix(fname, prefix) && strings.HasSuffix(fname, suffix) {
				path := filepath.Join(dir, fname)
				if !isDir(path) {
					res = append(res, path)
					continue
				}
				var tmp []*FileInfo
				doScanForInfo(path, &tmp)
				for _, item := range tmp {
					if item.IsDir() {
						continue
					}
					res = append(res, item.Path())
				}
			}
		}
	}
	return res
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
