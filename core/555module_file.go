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

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var regBlankChar = regexp.MustCompile(`\s+`)

// 读取文件所有内容，返回字节数组
func (this *InternalFunctionSet) Fbytes(filename string) []byte {
	bs, err := ioutil.ReadFile(filename)
	if err == nil {
		return bs
	}
	log.Fatal(fmt.Sprintf("failed to read %v file: %v", filename, err.Error()))
	return nil
}

// 读取文件所有内容，返回字节数组
func (this *InternalFunctionSet) Fbs(filename string) []byte {
	return this.Fbytes(filename)
}

// 读取文件所有内容，返回字符串
func (this *InternalFunctionSet) Fstr(filename string) string {
	bs, err := ioutil.ReadFile(filename)
	if err == nil {
		return string(bs)
	}
	log.Fatal(fmt.Sprintf("failed to read %v file: %v", filename, err.Error()))
	return ""
}

// 逐行读取文件，返回一个字符串数组
func (this *InternalFunctionSet) Flines(filename string, readGBK bool) []string {
	if readGBK {
		return readGBKByLine(filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to read %v file: %v", filename, err.Error()))
		return nil
	}
	var res []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	defer func() {
		if err = scanner.Err(); err != nil {
			runtimeExcption(err)
		}
		if err = file.Close(); err != nil {
			runtimeExcption(err)
		}
	}()
	return res
}

func readGBKByLine(filename string) []string {
	var enc = simplifiedchinese.GBK
	// Read UTF-8 from a GBK encoded file.
	f, err := os.Open(filename)
	if err != nil {
		runtimeExcption(err)
	}
	r := transform.NewReader(f, enc.NewDecoder())
	// Read converted UTF-8 from `r` as needed.
	// As an example we'll read line-by-line showing what was read:
	var res []string
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		res = append(res, sc.Text())
	}
	defer func() {
		if err = sc.Err(); err != nil {
			runtimeExcption(err)
		}

		if err = f.Close(); err != nil {
			runtimeExcption(err)
		}
	}()
	return res
}

// 读取文件内容，返回一个JSONObject
func (this *InternalFunctionSet) Fjson(filename string) map[string]interface{} {
	bs, err := os.ReadFile(filename)
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
func (this *InternalFunctionSet) Fprops(filename string) map[string]interface{} {
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
func (this *InternalFunctionSet) Fargs(filename string) []interface{} {
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

// 保存一个字符串或字节数组至指定文件（文件若已存在，清空其内容再保存）
func (this *InternalFunctionSet) Fsave(path string, content interface{}) {
	var data []byte
	if bs, ok := content.([]byte); ok {
		data = bs
	} else if str, ok := content.(string); ok {
		data = []byte(str)
	} else {
		runtimeExcption("function fsave(path, content): the parameter content must be type String/ByteArray")
	}

	fileSave(path, data)
}
func (this *InternalFunctionSet) Fsv(path string, content interface{}) {
	this.Fsave(path, content)
}
func fileSave(path string, data []byte) {
	err := os.WriteFile(path, data, 0666)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to write content to file: %v, %v", path, err.Error()))
	}
}

// 文件内容追加
func (this *InternalFunctionSet) Fappend(path string, content interface{}) {
	var data []byte
	if bs, ok := content.([]byte); ok {
		data = bs
	} else if str, ok := content.(string); ok {
		data = []byte(str)
	} else {
		runtimeExcption("function Fappend(path, content): the parameter content must be type String/ByteArray")
	}

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

// 获取当前路径下的所有子文件路径，不包含子目录
func (this *InternalFunctionSet) Fscan(path string) Value {
	var tmp []*FileInfo
	doScanForInfo(path, &tmp)
	var res []Value
	for _, item := range tmp {
		fi := item
		res = append(res, newClass("FileInfo", &fi))
	}
	return array(res)
}
