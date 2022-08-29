package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
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
var httpClients *http.Client

func getHttpClient() *http.Client {
	if httpClients == nil {
		httpClients = &http.Client{}
	}
	return httpClients
}

func (this *InternalFunctionSet) HttpGet(args []interface{}) Value {
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
	req, err := http.NewRequest(http.MethodGet, urlencoded(url), nil)
	if err != nil {
		runtimeExcption(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req.WithContext(ctx)
	if headers != nil {
		for key, val := range headers.mapVal() {
			req.Header.Set(key, val.String())
		}
	}
	resp, err := getHttpClient().Do(req)
	if err != nil {
		runtimeExcption(err)
	}

	return newHttpResponse(resp)
}

func (this *InternalFunctionSet) HttpPost(url string) Value {
	req, err := http.NewRequest(http.MethodPost, urlencoded(url), nil)
	if err != nil {
		runtimeExcption(err)
	}
	resp, err := getHttpClient().Do(req)
	if err != nil {
		runtimeExcption(err)
	}

	return newHttpResponse(resp)
}

func (this *InternalFunctionSet) HttpPostUrlencoded(url string, body JSONObject) Value {
	var content string
	for key, value := range body.mapVal() {
		if content == "" {
			content += fmt.Sprintf("%v=%v", key, value.String())
		} else {
			content += fmt.Sprintf("&%v=%v", key, value.String())
		}
	}
	req, err := http.NewRequest(http.MethodPost, urlencoded(url), strings.NewReader(content))
	if err != nil {
		runtimeExcption(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(content)))
	resp, err := getHttpClient().Do(req)
	if err != nil {
		runtimeExcption(err)
	}
	return newHttpResponse(resp)
}

func simpleRequestWithJson(method string, url string, body JSONObject) Value {
	content := body.toJSONObjectString()
	req, err := http.NewRequest(method, urlencoded(url), strings.NewReader(content))
	if err != nil {
		runtimeExcption(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(content)))
	resp, err := getHttpClient().Do(req)
	if err != nil {
		runtimeExcption(err)
	}
	return newHttpResponse(resp)
}

func (this *InternalFunctionSet) HttpPostJson(url string, body JSONObject) Value {
	return simpleRequestWithJson(http.MethodPost, url, body)
}

func readFileData(path string) []byte {
	bs, err := os.ReadFile(path)
	if err != nil {
		runtimeExcption(err)
	}
	return bs
}

func (this *InternalFunctionSet) HttpPostData(url string, body JSONObject) Value {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	for key, value := range body.mapVal() {
		if !toBoolean(value) {
			continue
		}
		if strings.HasPrefix(key, "@") { // file url
			fieldName := ifElse(key[1:] == "", "file", key[1:])
			fAttrs := parsePostFileAttrs(value)
			for _, fAttr := range fAttrs {
				fileName := fAttr.name
				fileData := readFileData(fAttr.path)
				fileWriter, _ := bodyWriter.CreateFormFile(fieldName, fileName)
				_, _ = io.Copy(fileWriter, bytes.NewReader(fileData))
			}
		} else if strings.HasPrefix(key, "#") { // file data
			fieldName := ifElse(key[1:] == "", "file", key[1:])
			fileName, fileData := "file", goBytes(value)
			fileWriter, _ := bodyWriter.CreateFormFile(fieldName, fileName)
			_, _ = io.Copy(fileWriter, bytes.NewReader(fileData))
		} else {
			_ = bodyWriter.WriteField(key, value.String())
		}
	}
	content, contentType, contentLen := bodyBuffer, bodyWriter.FormDataContentType(), bodyBuffer.Len()
	err := bodyWriter.Close()
	assert(err != nil, "failed to populate post args", err)

	req, err := http.NewRequest(http.MethodPost, urlencoded(url), content)
	if err != nil {
		runtimeExcption(err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", strconv.Itoa(contentLen))
	resp, err := getHttpClient().Do(req)
	if err != nil {
		runtimeExcption(err)
	}
	return newHttpResponse(resp)
}
func parsePostFileAttrs(raw Value) []PostFileAttr {
	var res []PostFileAttr
	if raw.isString() {
		v := raw.String()
		if strings.TrimSpace(v) == "" {
			return res
		}
		attr := PostFileAttr{
			name: filepath.Base(v),
			path: v,
		}
		res = append(res, attr)
	} else if raw.isJsonArray() {
		for _, item := range goArr(raw).values() {
			v := item.String()
			if strings.TrimSpace(v) == "" {
				continue
			}
			attr := PostFileAttr{
				name: filepath.Base(v),
				path: v,
			}
			res = append(res, attr)
		}
	} else {
	}
	return res
}

type PostFileAttr struct {
	name string
	path string
}

func (this *InternalFunctionSet) HttpPut(url string, body JSONObject) Value {
	return simpleRequestWithJson(http.MethodPut, url, body)
}

func (this *InternalFunctionSet) HttpPatch(url string, body JSONObject) Value {
	return simpleRequestWithJson(http.MethodPatch, url, body)
}

func (this *InternalFunctionSet) HttpDelete(url string, body JSONObject) Value {
	return simpleRequestWithJson(http.MethodDelete, url, body)
}

func (this *InternalFunctionSet) HttpHead(url string) Value {
	req, err := http.NewRequest(http.MethodHead, urlencoded(url), nil)
	if err != nil {
		runtimeExcption(err)
	}
	resp, err := getHttpClient().Do(req)
	if err != nil {
		runtimeExcption(err)
	}

	return newHttpResponse(resp)
}

func urlencoded(raw string) string {
	if !strings.Contains(raw, "?") {
		return raw
	}
	urlObj, err := url.Parse(raw)
	if err != nil {
		runtimeExcption(err)
	}
	queryVals, err := url.ParseQuery(urlObj.RawQuery)
	if err != nil {
		runtimeExcption(err)
	}
	return fmt.Sprintf(`%v://%v%v?%v`, urlObj.Scheme, urlObj.Host, urlObj.Path, queryVals.Encode())
}