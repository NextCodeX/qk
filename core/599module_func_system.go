package core

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os"
	"os/exec"
	"runtime"
	"time"
)


// 获取当前系统名称
func (fns *InternalFunctionSet) Sys() string {
	return runtime.GOOS
}

// 获取cpu的逻辑核数
func (fns *InternalFunctionSet) CpuNum() int {
	return runtime.NumCPU()
}

// 设置工作路径
func (fns *InternalFunctionSet) Setpwd(pwd string) {
	err := os.Chdir(pwd)
	if err != nil {
		fmt.Println(err)
	} else {
		mainFunc.setVar("pwd", newQKValue(pwd))
	}
}

// 当前协程进入休眠
func (fns *InternalFunctionSet) Sleep(t int64)  {
	time.Sleep(time.Duration(t) * time.Millisecond)
}

// 命令行调用
func (fns *InternalFunctionSet) Cmd(command string) string {
	var cmder *exec.Cmd
	if runtime.GOOS == "windows" {
		cmder = exec.Command("cmd", "/C", command)
	} else {
		cmder = exec.Command("bash", "-c", command)
	}

	output, err := cmder.CombinedOutput()
	if err != nil {
		fmt.Println(err, cmdResult(output))
		return ""
	}
	return cmdResult(output)
}

func cmdResult(bs []byte) string {
	bs = bytes.TrimSpace(bs)
	if runtime.GOOS != "windows" {
		return string(bs)
	}

	resBytes, err := simplifiedchinese.GBK.NewDecoder().Bytes(bs)
	if err != nil {
		return string(bs)
	} else {
		return string(resBytes)
	}
}

// 用浏览器打开相应的url
func (fns *InternalFunctionSet) OpenBrowser(url string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "start", url)
	} else {
		cmd = exec.Command("xdg-open", url)
	}
	res, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err, string(res))
	}
}
func (fns *InternalFunctionSet) Ob(url string) {
	fns.OpenBrowser(url)
}

