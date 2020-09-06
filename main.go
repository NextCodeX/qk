package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"qk/core"
	"runtime"
	"strings"
)

func main() {
	//qkfile := "examples/string-module-test2.qk"
	//qkfile := "examples/db-test.qk"
	qkfile := getScriptFile()
	changeWorkDirectory()

	bs, _ := ioutil.ReadFile(qkfile)
	core.Run(bs)
}

// change work dirctory to current command directory
func changeWorkDirectory() {
	cmdDir := getCmdDir()
	err := os.Chdir(cmdDir)
	if err != nil {
		fmt.Printf("failed to change work dirctory to current command directory: %v", err.Error())
		os.Exit(5)
	}
}

// find qk script file for run
func getScriptFile() string {
	cmdDir := getCmdDir()
	if len(os.Args)>1 {
		arg := os.Args[1]
		if strings.HasPrefix(arg, "abs=") {
			return arg[4:]
		}
		return filepath.Join(cmdDir, arg)
	}

	fs, err := ioutil.ReadDir(cmdDir)
	if err != nil {
		fmt.Printf("failed to get script file: %v", err.Error())
		os.Exit(5)
	}
	var fnames []string
	for _, f := range fs {
		fname := f.Name()
		if strings.HasSuffix(fname, ".qk") {
			fnames = append(fnames, fname)
		}
	}
	if len(fnames)<1 {
		fmt.Println("qk script file is not found in current directory!")
		os.Exit(5)
	}
	if len(fnames)>1 {
		fmt.Println("Multiple qk script file exist in current directory!")
		os.Exit(5)
	}
	return filepath.Join(cmdDir, fnames[0])
}

// 获取命令所在的路径
func getCmdDir() string {
	cmd := exec.Command("cmd", "/c", "cd")
	if runtime.GOOS != "windows" {
		cmd = exec.Command("pwd")
	}
	d, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("failed to get personal work directory: %v \n", err.Error())
		os.Exit(5)
	}
	pwd := strings.TrimSpace(string(d))
	return pwd
}

