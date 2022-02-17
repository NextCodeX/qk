package core

import (
	"fmt"
	"path/filepath"
)

// 指定路径的目录
func (fns *InternalFunctionSet) PathDir(path string) string {
	return filepath.Dir(path)
}

// 指定路径的文件名
func (fns *InternalFunctionSet) PathBase(path string) string {
	return filepath.Base(path)
}

// 路径拼接
func (fns *InternalFunctionSet) PathJoin(paths []string) string {
	return filepath.Join(paths...)
}

// 返回绝对路径
func (fns *InternalFunctionSet) PathAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("PathAbs:", err)
		return ""
	}
	return abs
}

// 文件名后缀
func (fns *InternalFunctionSet) PathExt(path string) string {
	return filepath.Ext(path)
}

// 路径计算
func (fns *InternalFunctionSet) PathClean(path string) string {
	return filepath.Clean(path)
}

// 路径匹配
func (fns *InternalFunctionSet) PathMatch(pattern, path string) bool {
	flag, err := filepath.Match(pattern, path)
	if err != nil {
		fmt.Println("PathMatch:", err)
		return false
	}
	return flag
}
