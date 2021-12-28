package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 查找待执行的脚本文件。
// 查找成功，返回文件路径；
// 查找失败，直接报错退出程序。
func findScriptFile() string {
	cmdDir := getCmdDir()
	if len(os.Args) > 1 {
		arg := os.Args[1]
		// 允许运行脚本文件时，不指定脚本文件名。
		if !strings.HasSuffix(arg, ".qk") {
			arg = arg + ".qk"
		}
		if fileExist(arg) {
			// absolute path or available relative path
			return arg
		}

		// workspace path
		// 在当前目录下查找
		if currentDirFile := filepath.Join(cmdDir, arg); fileExist(currentDirFile) {
			return currentDirFile
		}

		// environment path
		// 在配置好的环境变量目录下查找
		qkHome := os.Getenv("QK_HOME")
		if qkHome != "" {
			if envDirFile := filepath.Join(qkHome, arg); fileExist(envDirFile) {
				return envDirFile
			}
		}

		runtimeExcption(arg, " is not found!")
		return ""
	}

	// 当前目录只有一个qk文件可以忽略文件名，没有或有多个时程序报错退出
	fs, err := ioutil.ReadDir(cmdDir)
	if err != nil {
		errorf("failed to get script file: %v", err.Error())
	}
	var targetScript string
	for _, f := range fs {
		fname := f.Name()
		if strings.HasSuffix(fname, ".qk") {
			if targetScript != "" {
				runtimeExcption("Multiple qk script file exist in current directory!")
			}
			targetScript = fname
		}
	}
	if targetScript == "" {
		runtimeExcption("qk script file is not found in current directory!")
	}
	return filepath.Join(cmdDir, targetScript)
}

// 获取命令所在的路径
func getCmdDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		errorf("failed to get current work directory: %v \n", err.Error())
	}
	return cwd
}
