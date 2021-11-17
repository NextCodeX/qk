package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type HttpRequest struct {
	headers map[string][]string
	args    map[string][]string
	files   map[string][]Value
	body    []byte
}

func newHttpRequest(req *http.Request, method string) Value {
	obj := &HttpRequest{}
	parseArgs(obj, req, method)
	obj.headers = req.Header

	return newClass("HttpRequest", &obj)
}

func parseArgs(obj *HttpRequest, req *http.Request, method string) {
	contentType := req.Header.Get("Content-Type")
	if method == "GET" || strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		_ = req.ParseForm()
		obj.args = req.Form
		return
	}

	obj.args = req.URL.Query()
	isPost := method == "POST"
	if isPost &&
		(strings.HasPrefix(contentType, "text/plain") ||
			strings.HasPrefix(contentType, "application/json")) {
		obj.body, _ = ioutil.ReadAll(req.Body)
		return
	}

	if isPost && strings.HasPrefix(contentType, "multipart/form-data") {
		reader, err := req.MultipartReader()
		if err != nil {
			log.Fatal(err)
		}

		obj.files = make(map[string][]Value)
		for {
			part, err := reader.NextPart()
			if err == io.EOF || part == nil {
				break
			}
			formName := part.FormName()
			fileName := part.FileName()
			if fileName != "" {
				fs := obj.files[formName]
				bs, _ := ioutil.ReadAll(part)
				fs = append(fs, newMultipartFile(fileName, bs))

				obj.files[formName] = fs
				continue
			}

			fieldVals := obj.args[formName]
			val, _ := ioutil.ReadAll(part)
			fieldVals = append(fieldVals, string(val))

			obj.args[formName] = fieldVals
		}
	}
}

func (req *HttpRequest) ShowHeaders() {
	fmt.Println("=============================")
	for k, v := range req.headers {
		fmt.Printf("%v: %v\n", k, v)
	}
	fmt.Println("=============================")
}

func (req *HttpRequest) Showhs() {
	req.ShowHeaders()
}

func (req *HttpRequest) Headers() JSONObject {
	return httpValsToJSONObject(req.headers)
}

func (req *HttpRequest) GetHeaderVal(name string) interface{} {
	ps := req.args[name]
	if len(ps) > 0 {
		return ps[0]
	}
	return nil
}

func (req *HttpRequest) H(name string) interface{} {
	return req.GetHeaderVal(name)
}

func (req *HttpRequest) GetHeaderVals(name string) []string {
	ps := req.args[name]
	if len(ps) > 0 {
		return ps
	}
	return nil
}

func (req *HttpRequest) Hs(name string) []string {
	return req.GetHeaderVals(name)
}

func (req *HttpRequest) Args() JSONObject {
	res := httpValsToJSONObject(req.args)
	obj := req.Json()
	if obj != nil {
		for k, v := range obj {
			res.put(k, newQKValue(v))
		}
	}
	return res
}

func (req *HttpRequest) Json() map[string]interface{} {
	var res interface{}
	_ = json.Unmarshal(req.body, &res)
	if res != nil {
		return res.(map[string]interface{})
	}
	return nil
}

func (req *HttpRequest) Body() string {
	return string(req.body)
}

func (req *HttpRequest) BodyBytes() []byte {
	return req.body
}

func (req *HttpRequest) Param(name string) interface{} {
	ps := req.args[name]
	if len(ps) > 0 {
		return ps[0]
	}
	return nil
}

func (req *HttpRequest) Params(name string) interface{} {
	ps := req.args[name]
	if len(ps) > 0 {
		return ps
	}
	return nil
}

func (req *HttpRequest) Files() Value {
	res := make(map[string]Value)
	for k, multiFiles := range req.files {
		var fs []Value
		for _, multiFile := range multiFiles {
			fs = append(fs, multiFile)
		}
		res[k] = array(fs)
	}
	return jsonObject(res)
}

func (req *HttpRequest) GetFiles(formName string) JSONArray {
	for k, multiFiles := range req.files {
		if k != formName {
			continue
		}
		var fs []Value
		for _, multiFile := range multiFiles {
			fs = append(fs, multiFile)
		}
		return array(fs)
	}
	return nil
}

func (req *HttpRequest) GetFile(formName string) Value {
	for k, multiFiles := range req.files {
		if k != formName {
			continue
		}
		for _, multiFile := range multiFiles {
			return multiFile
		}
	}
	return nil
}

type MultipartFile struct {
	name string
	data []byte
}

func newMultipartFile(name string, data []byte) Value {
	obj := &MultipartFile{name: name, data: data}
	return newClass("MultipartFile", &obj)
}

func (f *MultipartFile) Name() string {
	return f.name
}

func (f *MultipartFile) Size() int {
	return len(f.data)
}

func (f *MultipartFile) Bytes() []byte {
	return f.data
}

func (f *MultipartFile) Save(path string) int64 {
	dst, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		runtimeExcption(err)
	}
	size, err := io.Copy(dst, bytes.NewReader(f.data))
	if err != nil {
		runtimeExcption(err)
	}
	return size
}

func httpValsToJSONObject(vals map[string][]string) JSONObject {
	obj := make(map[string]Value)
	for k, vals := range vals {
		var arr []Value
		for _, val := range vals {
			arr = append(arr, newQKValue(val))
		}
		obj[k] = array(arr)
	}
	return jsonObject(obj)
}
