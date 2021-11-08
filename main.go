package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"qk/core"
	"runtime"
	"strings"
	"time"
)

func main() {
	startupTime := time.Now().UnixNano()

	if len(os.Args) > 1 {
		if arg := os.Args[1]; arg == "-v" {
			fmt.Println("Quick version:", version)
			return
		}
	}

	qkfile := "examples/demo.qk"
	//qkfile := getScriptFile()

	bs, err := ioutil.ReadFile(qkfile)
	if err != nil {
		log.Fatalf("failed to read %v; error info: %v", qkfile, err)
	}
	core.SetRootDir(qkfile)
	core.Run(bs)

	if core.IsCost {
		duration := time.Now().UnixNano() - startupTime
		fmt.Printf("\n\nspend: %vns, %.3fms, %.3fs  \n", duration, float64(duration)/1e6, float64(duration)/1e9)
	}
}

// find qk script file to run
func getScriptFile() string {
	cmdDir := getCmdDir()
	if len(os.Args) > 1 {
		arg := os.Args[1]
		// optimize script running
		if !strings.HasSuffix(arg, ".qk") {
			arg = arg + ".qk"
		}
		if fileExist(arg) {
			// absolute path or available relative path
			return arg
		}

		// workspace path
		if currentDirFile := filepath.Join(cmdDir, arg); fileExist(currentDirFile) {
			return currentDirFile
		}

		// environment path
		qkHome := os.Getenv("QK_HOME")
		if qkHome != "" {
			if envDirFile := filepath.Join(qkHome, arg); fileExist(envDirFile) {
				return envDirFile
			}
		}

		log.Fatal(arg + " is not found!")
		return ""
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
	if len(fnames) < 1 {
		fmt.Println("qk script file is not found in current directory!")
		os.Exit(5)
	}
	if len(fnames) > 1 {
		fmt.Println("Multiple qk script file exist in current directory!")
		os.Exit(5)
	}
	return filepath.Join(cmdDir, fnames[0])
}

// 判断文件是否存在
func fileExist(path string) bool {
	sta, err := os.Stat(path)
	return sta != nil && !sta.IsDir() && !os.IsNotExist(err)
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
