package core

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// 获取当前系统名称
func (this *InternalFunctionSet) Sys() string {
	return runtime.GOOS
}

// 获取cpu的逻辑核数
func (this *InternalFunctionSet) CpuNum() int {
	return runtime.NumCPU()
}

// 活跃协程数量
func (this *InternalFunctionSet) RoutineNum() int {
	return runtime.NumGoroutine()
}

// 设置工作路径
func (this *InternalFunctionSet) Setpwd(pwd string) {
	var err error
	pwd, err = filepath.Abs(pwd)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.Chdir(pwd)
	if err != nil {
		fmt.Println(err)
	} else {
		this.owner.main.setVar("pwd", newQKValue(pwd))
	}
}
func (this *InternalFunctionSet) Cd(pwd string) {
	this.Setpwd(pwd)
}

// 获取所有环境变量
func (this *InternalFunctionSet) Envs() []string {
	return os.Environ()
}

// 获取单个环境变量
func (this *InternalFunctionSet) Env(key string) string {
	return os.Getenv(key)
}

// 设置环境变量(仅对当前脚本有效)
func (this *InternalFunctionSet) Setenv(key, val string) {
	err := os.Setenv(key, val)
	if err != nil {
		fmt.Println(err)
	}
}

// 获取并打印环境变量QK_HOME
func (this *InternalFunctionSet) Home() string {
	home := os.Getenv("QK_HOME")
	if home == "" {
		fmt.Println("QK_HOME is not configured")
	} else {
		fmt.Println("QK_HOME:", home)
	}
	return home
}

// 统计脚本耗时
var IsCost = false

func (this *InternalFunctionSet) Cost() {
	IsCost = true
}

func (this *InternalFunctionSet) Assert(flag bool) {
	if !flag {
		//获取的是 CallerA函数的调用者的调用栈
		pc, file, lineNo, ok := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc).Name()
		panic(fmt.Sprintf("%v; %v; %v; %v; %v\n", pc, file, lineNo, ok, funcName))
	}
}

// 当前协程进入休眠
func (this *InternalFunctionSet) Sleep(t int64) {
	time.Sleep(time.Duration(t) * time.Millisecond)
}

// 命令行调用
func (this *InternalFunctionSet) Cmd(command string) string {
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
func (this *InternalFunctionSet) OpenBrowser(url string) {
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
func (this *InternalFunctionSet) Ob(url string) {
	this.OpenBrowser(url)
}
