package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

/*
常见的媒体格式类型如下：
text/html ： HTML格式
text/plain ：纯文本格式
text/xml ： XML格式
image/gif ：gif图片格式
image/jpeg ：jpg图片格式
image/png：png图片格式

以application开头的媒体格式类型：
application/xhtml+xml ：XHTML格式
application/xml： XML数据格式
application/atom+xml ：Atom XML聚合格式
application/json： JSON数据格式
application/pdf：pdf格式
application/msword ： Word文档格式
application/octet-stream ： 二进制流数据（如常见的文件下载）
application/x-www-form-urlencoded ： <form encType="">中默认的encType，form表单数据被编码为key/value格式发送到服务器（表单默认的提交数据的格式）
multipart/form-data ： 需要在表单中进行文件上传时，就需要使用该格式
*/
// cookie, set-cookie

func (fns *InternalFunctionSet) HttpGet(args []interface{}) Value {
	if len(args) < 1 {
		runtimeExcption("HttpGet() url is required")
	}
	url, ok := args[0].(string)
	if !ok {
		runtimeExcption("HttpGet() url must be string type")
	}
	var headers JSONObject
	if len(args) > 1 {
		headers = args[1].(JSONObject)
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		runtimeExcption(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	req.WithContext(ctx)
	if headers != nil {
		for key, val := range headers.mapVal() {
			req.Header.Set(key, val.String())
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		runtimeExcption(err)
	}

	return newHttpResponse(resp)
}

func (fns *InternalFunctionSet) HttpPost(args []interface{}) Value {
	if len(args) < 1 {
		runtimeExcption("HttpPost() url is required")
	}
	url, ok := args[0].(string)
	if !ok {
		runtimeExcption("HttpPost() url must be string type")
	}
	var config JSONObject
	if len(args) > 1 {
		config = args[1].(JSONObject)
	}
	if config == nil {
		runtimeExcption("HttpPost() config is required and must be string type")
	}
	contentType := config.get("type").String()
	content, contentType := parseBody(contentType, config.get("action"))

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, content)
	if err != nil {
		runtimeExcption(err)
	}
	req.Header.Set("Content-Type", contentType)
	headersRaw := config.get("headers")
	if !headersRaw.isNULL() {
		headers := headersRaw.(JSONObject)
		for key, val := range headers.mapVal() {
			req.Header.Set(key, val.String())
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		runtimeExcption(err)
	}

	return newHttpResponse(resp)
}

func (fns *InternalFunctionSet) HttpHead(url string) Value {
	client := &http.Client{}
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		runtimeExcption(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		runtimeExcption(err)
	}

	return newHttpResponse(resp)
}

func parseBody(contentType string, val Value) (io.Reader, string) {
	if contentType == "application/x-www-form-urlencoded" {
		obj := val.(JSONObject)
		var res string
		for key, value := range obj.mapVal() {
			if res == "" {
				res += fmt.Sprintf("%v=%v", key, value.String())
			} else {
				res += fmt.Sprintf("&%v=%v", key, value.String())
			}
		}
		return strings.NewReader(res), contentType
	}
	if contentType == "multipart/form-data" {
		obj := val.(JSONObject)
		bodyBuffer := &bytes.Buffer{}
		bodyWriter := multipart.NewWriter(bodyBuffer)

		for key, value := range obj.mapVal() {
			if key == "files" {
				arr := value.(JSONArray)
				for _, elem := range arr.values() {
					info := elem.(JSONObject)
					fileNameVal := info.get("name")
					fieldNameVal := info.get("field")
					pathVal := info.get("path")
					dataVal := info.get("data")
					var fileName, fieldName, path string
					var data []byte
					if !fieldNameVal.isNULL() {
						fieldName = fieldNameVal.String()
					} else {
						fieldName = "files"
					}
					if !fileNameVal.isNULL() {
						fileName = fileNameVal.String()
					} else {
						runtimeExcption("fileName is be required")
					}
					if !dataVal.isNULL() {
						data = goBytes(dataVal)
					}
					if data == nil && !pathVal.isNULL() {
						path = pathVal.String()
						bs, err := ioutil.ReadFile(path)
						if err != nil {
							runtimeExcption(err)
						}
						data = bs
					}
					if data == nil {
						runtimeExcption("file data or path is be required")
					}
					fileWriter, _ := bodyWriter.CreateFormFile(fieldName, fileName)
					_, _ = io.Copy(fileWriter, bytes.NewReader(data))
				}
				continue
			}
			_ = bodyWriter.WriteField(key, value.String())
		}

		return bodyBuffer, bodyWriter.FormDataContentType()
	}

	return strings.NewReader(val.String()), contentType
}
