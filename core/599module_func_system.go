package core

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"time"
)

// 当前协程进入休眠
func (fns *InternalFunctionSet) Sleep(t int64)  {
	time.Sleep(time.Duration(t) * time.Millisecond)
}

// 命令行调用
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

// base64 编码
func (fns *InternalFunctionSet) Base64Encode(arg interface{}) string {
	if raw, ok := arg.([]byte); ok {
		return base64.StdEncoding.EncodeToString(raw)
	} else if raw, ok := arg.(string); ok {
		return base64.StdEncoding.EncodeToString([]byte(raw))
	} else{
		fmt.Println("base64Encode() the first parameter type must be ByteArray or String")
	}
	return ""
}
func (fns *InternalFunctionSet) Base64(arg interface{}) string {
	return fns.Base64Encode(arg)
}

// base64 解码
func (fns *InternalFunctionSet) Base64Decode(raw string) []byte {
	 data, err := base64.StdEncoding.DecodeString(raw)
	 if err != nil {
	 	fmt.Println(err)
	 }
	 return data
}
func (fns *InternalFunctionSet) Debase64(raw string) []byte {
	 return fns.Base64Decode(raw)
}

// gzip 解码 （解压缩）
func (fns *InternalFunctionSet) GzipDecode(data []byte) []byte {
	bytesReader := bytes.NewReader(data)
	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil {
		runtimeExcption(err)
	}
	res, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		runtimeExcption(err)
	}
	return res
}
func (fns *InternalFunctionSet) Degzip(data []byte) []byte {
	return fns.GzipDecode(data)
}

// gzip 编码 （压缩）
func (fns *InternalFunctionSet) GzipEncode(src []byte) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(src); err != nil {
		runtimeExcption(err)
	}
	if err := gz.Flush(); err != nil {
		runtimeExcption(err)
	}
	if err := gz.Close(); err != nil {
		runtimeExcption(err)
	}
	return buf.Bytes()
}
func (fns *InternalFunctionSet) Gzip(src []byte) []byte {
	return fns.GzipEncode(src)
}