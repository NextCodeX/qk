package core

import (
	"fmt"
	"path/filepath"
)

// 指定路径的目录
func (this *InternalFunctionSet) PathDir(path string) string {
	return filepath.Dir(path)
}

// 指定路径的文件名
func (this *InternalFunctionSet) PathBase(path string) string {
	return filepath.Base(path)
}

// 路径拼接
func (this *InternalFunctionSet) PathJoin(paths []string) string {
	return filepath.Join(paths...)
}

// 返回绝对路径
func (this *InternalFunctionSet) PathAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("PathAbs:", err)
		return ""
	}
	return abs
}

// 文件名后缀
func (this *InternalFunctionSet) PathExt(path string) string {
	return filepath.Ext(path)
}

// 路径计算
func (this *InternalFunctionSet) PathClean(path string) string {
	return filepath.Clean(path)
}

// 路径匹配
func (this *InternalFunctionSet) PathMatch(pattern, path string) bool {
	flag, err := filepath.Match(pattern, path)
	if err != nil {
		fmt.Println("PathMatch:", err)
		return false
	}
	return flag
}
