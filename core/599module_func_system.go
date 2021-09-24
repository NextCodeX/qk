package core

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

// 当前协程进入休眠
func (fns *InternalFunctionSet) Sleep(t int64) {
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
		fmt.Println(err, string(output))
		return ""
	}
	return string(bytes.TrimSpace(output))
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
