package core

import (
	"fmt"
	"os"
	"path/filepath"
)

var internalVars = make(map[string]Value)

// 用于加载内置变量与内置函数
func initInternalVars() {
	// 添加命令行参数
	rawCmdArgs := os.Args
	if len(rawCmdArgs) > 2 {
		var arr []Value
		for _, rawCmdArg := range rawCmdArgs[2:] {
			arr = append(arr, newQKValue(rawCmdArg))
		}
		internalVars["cmdArgs"] = array(arr)
	} else {
		internalVars["cmdArgs"] = emptyArray()
	}

	// 提供qk执行文件所在的目录
	if executable, err := os.Executable(); err == nil {
		rootDir := filepath.Dir(executable)
		internalVars["qkDir"] = newQKValue(rootDir)
	} else {
		fmt.Println(err)
	}

	// 当前命令行所在的路径，与`pwd`等同(工作路径)
	if cwd, err := os.Getwd(); err == nil {
		internalVars["pwd"] = newQKValue(cwd)
	} else {
		fmt.Println(err)
	}

	// 添加POST请求常用的Content-Type
	mimes := make(map[string]Value)
	mimes["txt"] = newQKValue("text/plain")
	mimes["json"] = newQKValue("application/json")
	mimes["form"] = newQKValue("application/x-www-form-urlencoded")
	mimes["data"] = newQKValue("multipart/form-data")
	internalVars["mime"] = jsonObject(mimes)

	// 将内部函数注册到主函数(main)的内部栈中
	fns := &InternalFunctionSet{}
	fmap := collectFunctionInfo(&fns)
	for fname, moduleFunc := range fmap {
		internalVars[fname] = newInternalFunc(fname, moduleFunc)
	}
}

func setRootDir(scriptFileName string) {
	// 提供当前脚本文件所在的目录
	if dir, err := filepath.Abs(filepath.Dir(scriptFileName)); err == nil {
		internalVars["root"] = newQKValue(dir)
	} else {
		fmt.Println(err)
	}
}
