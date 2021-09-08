package core

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// 添加 Quick 系统内部函数
func addInternalFunc(name string, internalFunc func([]interface{})interface{}) {
	mainFunc.addFunc(name, newInternalFunc(name, internalFunc))
}
func addModuleFunc(name string, moduleFunc *FunctionExecutor) {
	mainFunc.addFunc(name, newModuleFunc(name, moduleFunc))
}

func init() {
	// 添加命令行参数
	rawCmdArgs := os.Args
	if len(rawCmdArgs) > 2 {
		var arr []Value
		for _, rawCmdArg := range rawCmdArgs[2:] {
			arr = append(arr, newQKValue(rawCmdArg))
		}
		mainFunc.setPreVar("cmdArgs", array(arr))
	} else {
		mainFunc.setPreVar("cmdArgs", emptyArray())
	}

	// 提供qk执行文件所在的目录
	if executable, err := os.Executable(); err == nil {
		rootDir := path.Dir(executable)
		mainFunc.setPreVar("qkDir", newQKValue(rootDir))
	} else {
		fmt.Println(err)
	}

	// 当前命令行所在的路径，与`pwd`等同
	if cwd, err := os.Getwd(); err==nil {
		mainFunc.setPreVar("pwd", newQKValue(cwd))
	}else{
		fmt.Println(err)
	}
}

func SetRootDir(scriptFileName string) {
	// 提供当前脚本文件所在的目录
	if dir, err := filepath.Abs(filepath.Dir(scriptFileName)); err == nil {
		mainFunc.setPreVar("root", newQKValue(dir))
	} else {
		fmt.Println(err)
	}
}
