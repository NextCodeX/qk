package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"bufio"
	"encoding/json"
	"strings"
	"regexp"
	"path/filepath"
	"io"
)

func fileModuleInit()  {
	f := &File{}
	collectFunctionInfo(&f, "file")
}

var regBlankChar = regexp.MustCompile(`\s+`)

type File struct {

}

func (f *File) Bytes(filename string) []byte {
	bs, err := ioutil.ReadFile(filename)
	if err == nil {
		return bs
	}
	log.Fatal(fmt.Sprintf("failed to read %v file: %v", filename, err.Error()))
	return nil
}

func (f *File) Content(filename string) string {
	bs, err := ioutil.ReadFile(filename)
	if err == nil {
		return string(bs)
	}
	log.Fatal(fmt.Sprintf("failed to read %v file: %v", filename, err.Error()))
	return ""
}

func (f *File) Lines(filename string) []interface{} {
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

func (f *File) Json(filename string) map[string]interface{} {
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

func (f *File) Args(filename string) []interface{} {
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

func (f *File) Scan(path string) []interface{} {
	var res []interface{}
	doScan(path, &res)
	return res
}

func doScan(path string, res *[]interface{})  {
	if !isDir(path) {
		log.Fatal(path, "is not directory.")
		return
	}
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to read path: %v, %v", path, err.Error()))
		return
	}
	for _, f := range fs {
		nextPath := filepath.Join(path, f.Name())
		if f.IsDir() {
			*res = append(*res, nextPath)
			doScan(nextPath, res)
			continue
		}
		*res = append(*res, nextPath)
	}
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to read path: %v, %v", path, err.Error()))
		return false
	}
	return fileInfo.IsDir()
}

func (f *File) Out(path, content string) {
	err := ioutil.WriteFile(path, []byte(content), 0666)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to write content to file: %v, %v", path, err.Error()))
	}
}

func (f *File) Append(path, content string) {
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