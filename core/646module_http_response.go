package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type HttpResponse struct {
	headers    map[string][]string
	body       []byte
	cookies    []*http.Cookie
	status     string
	statusCode int
}

func newHttpResponse(resp *http.Response) Value {
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			runtimeExcption(err)
		}
	}()

	obj := &HttpResponse{cookies: resp.Cookies()}
	obj.status = resp.Status
	obj.statusCode = resp.StatusCode
	obj.headers = resp.Header
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(resp.StatusCode, err)
		obj.statusCode = 504
	}
	obj.body = body

	return newClass("HttpResponse", &obj)
}

func (resp *HttpResponse) ShowCookies() {
	fmt.Println("+++++++++++++++++++++++++++++")
	for _, ck := range resp.cookies {
		fmt.Printf("%v=%v; Domain=%v; Path=%v; Expires=%v; MaxAge=%v; HttpOnly=%v; Secure=%v; SameSite=%v\n", ck.Name, ck.Value, ck.Domain, ck.Path, ck.Expires, ck.MaxAge, ck.HttpOnly, ck.Secure, ck.SameSite)
	}
	fmt.Println("+++++++++++++++++++++++++++++")
}

func (resp *HttpResponse) Showck() {
	resp.ShowCookies()
}

func (resp *HttpResponse) Cookies() map[string]string {
	res := make(map[string]string)
	for _, ck := range resp.cookies {
		res[ck.Name] = ck.Value
	}
	return res
}

func (resp *HttpResponse) ShowHeaders() {
	fmt.Println("=============================")
	for k, v := range resp.headers {
		fmt.Printf("%v: %v\n", k, v)
	}
	fmt.Println("=============================")
}

func (resp *HttpResponse) Showhs() {
	resp.ShowHeaders()
}

func (resp *HttpResponse) Headers() JSONObject {
	return httpValsToJSONObject(resp.headers)
}

func (resp *HttpResponse) Hs() JSONObject {
	return resp.Headers()
}

func (resp *HttpResponse) Status() string {
	return resp.status
}

func (resp *HttpResponse) StatusCode() int {
	return resp.statusCode
}

func (resp *HttpResponse) Code() int {
	return resp.StatusCode()
}

func (resp *HttpResponse) Ok() bool {
	return resp.StatusCode() == 200
}

func (resp *HttpResponse) Save(path string) {
	err := os.WriteFile(path, resp.body, 0666)
	if err != nil {
		runtimeExcption(err)
	}
}

func (resp *HttpResponse) String() string {
	return string(resp.body)
}
func (resp *HttpResponse) Str() string {
	return resp.String()
}

func (resp *HttpResponse) Bytes() []byte {
	return resp.body
}
func (resp *HttpResponse) Bs() []byte {
	return resp.Bytes()
}

func (resp *HttpResponse) Pretty() {
	var out bytes.Buffer
	err := json.Indent(&out, resp.body, "", "  ")
	if err != nil {
		fmt.Printf("%s \n", resp.body)
	} else {
		fmt.Println(out.String())
	}
}
func (resp *HttpResponse) Pr() {
	resp.Pretty()
}

func (resp *HttpResponse) Json() interface{} {
	var tmp interface{}
	err := json.Unmarshal(resp.body, &tmp)
	if err != nil {
		return nil
	}

	if obj, ok := tmp.(map[string]interface{}); ok {
		return obj
	} else if arr, ok := tmp.([]interface{}); ok {
		return arr
	} else {
		return nil
	}
}
