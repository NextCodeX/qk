package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)


var regBlankChar = regexp.MustCompile(`\s+`)

// 读取文件所有内容，返回字节数组
func (fns *InternalFunctionSet) Fbytes(filename string) []byte {
	bs, err := ioutil.ReadFile(filename)
	if err == nil {
		return bs
	}
	log.Fatal(fmt.Sprintf("failed to read %v file: %v", filename, err.Error()))
	return nil
}
// 读取文件所有内容，返回字节数组
func (fns *InternalFunctionSet) Fbs(filename string) []byte {
	return fns.Fbytes(filename)
}

// 读取文件所有内容，返回字符串
func (fns *InternalFunctionSet) Fstr(filename string) string {
	bs, err := ioutil.ReadFile(filename)
	if err == nil {
		return string(bs)
	}
	log.Fatal(fmt.Sprintf("failed to read %v file: %v", filename, err.Error()))
	return ""
}

// 逐行读取文件，返回一个字符串数组
func (fns *InternalFunctionSet) Flines(filename string) []interface{} {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to read %v file: %v", filename, err.Error()))
		return nil
	}
	var res []interface{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	return res
}

// 读取文件内容，返回一个JSONObject
func (fns *InternalFunctionSet) Fjson(filename string) map[string]interface{} {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to read %v file: %v", filename, err.Error()))
		return nil
	}
	var tmp interface{}
	err = json.Unmarshal(bs, &tmp)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to parse json[%v]: %v", string(bs), err.Error()))
		return nil
	}
	res := tmp.(map[string]interface{})
	return res
}

// 读取*.properties文件，返回一个JSONObject
func (fns *InternalFunctionSet) Fprops(filename string) map[string]interface{} {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to read %v file: %v", filename, err.Error()))
		return nil
	}
	res := make(map[string]interface{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if "" == line || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		arr := strings.Split(line, "=")
		key := strings.TrimSpace(arr[0])
		val := strings.TrimSpace(arr[1])
		res[key] = val
	}
	return res
}

// 读取一个参数文件，返回二维数组
func (fns *InternalFunctionSet) Fargs(filename string) []interface{} {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to read %v file: %v", filename, err.Error()))
		return nil
	}
	var res []interface{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		var args []interface{}
		var subList []string
		if strings.Contains(line, ",") {
			subList = strings.Split(line, ",")
		} else {
			line = regBlankChar.ReplaceAllString(line, ",")
			subList = strings.Split(line, ",")
		}
		for _, sub := range subList {
			arg := strings.TrimSpace(sub)
			args = append(args, arg)
		}

		res = append(res, args)
	}
	return res
}

func (fns *InternalFunctionSet) Fscan(path string) []interface{} {
	var res []interface{}
	doScan(path, false, &res)
	return res
}

func (fns *InternalFunctionSet) FscanAll(path string) []interface{} {
	var res []interface{}
	doScan(path, true, &res)
	return res
}

func (fns *InternalFunctionSet) Fout(path, content string) {
	err := ioutil.WriteFile(path, []byte(content), 0666)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to write content to file: %v, %v", path, err.Error()))
	}
}

func (fns *InternalFunctionSet) Fappend(path, content string) {
	data := []byte(content)
	fobj, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to open file: %v, %v", path, err.Error()))
	}
	n, err := fobj.Write(data)
	if err == nil && n < len(data) {
		log.Fatal(fmt.Sprintf("failed to write content to file: %v, %v", path, io.ErrShortWrite.Error()))
	}
	if err1 := fobj.Close(); err == nil && err1 != nil {
		log.Fatal(fmt.Sprintf("failed to close file: %v, %v", path, err1.Error()))
	}
}

func (fns *InternalFunctionSet) Fsave(path string, bytes []byte) {
	err := ioutil.WriteFile(path, bytes, 0666)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to write content to file: %v, %v", path, err.Error()))
	}
}

func (fns *InternalFunctionSet) FappendBytes(path string, bytes []byte) {
	fobj, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to open file: %v, %v", path, err.Error()))
	}
	n, err := fobj.Write(bytes)
	if err == nil && n < len(bytes) {
		log.Fatal(fmt.Sprintf("failed to write content to file: %v, %v", path, io.ErrShortWrite.Error()))
	}
	if err1 := fobj.Close(); err == nil && err1 != nil {
		log.Fatal(fmt.Sprintf("failed to close file: %v, %v", path, err1.Error()))
	}
}
