package core

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

func (fns *InternalFunctionSet) Sleep(t int64)  {
	time.Sleep(time.Duration(t) * time.Millisecond)
}

func (fns *InternalFunctionSet) Cmd(command string) string {
	if command == "" {
		return ""
	}
	var executor string
	var args []string
	var collectFlag bool
	var tmpBytes []byte
	rawBytes := []byte(command)
	for _, b := range rawBytes {
		if !collectFlag && b != ' ' {
			collectFlag = true
		}
		if collectFlag && b != ' ' {
			tmpBytes = append(tmpBytes, b)
		}
		if collectFlag && b == ' ' {
			collectFlag = false
			if executor == "" {
				executor = string(tmpBytes)
			} else {
				args = append(args, string(tmpBytes))
			}
			tmpBytes = nil
		}
	}
	if collectFlag {
		if executor == "" {
			executor = string(tmpBytes)
		} else {
			args = append(args, string(tmpBytes))
		}
		tmpBytes = nil
	}
	//fmt.Println("++++++++++++++++++++")
	//fmt.Println("cmd: ", command)
	//fmt.Println("executor:", executor)
	//fmt.Println("args:", args, len(args))
	res := exec.Command(executor, args...)
	output, err := res.CombinedOutput()
	if err != nil {
		fmt.Println(err, string(output))
		return ""
	}
	return string(bytes.TrimSpace(output))
}