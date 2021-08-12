package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpResponse struct {
	headers    map[string][]string
	body       []byte
	cookies    []*http.Cookie
	status     string
	statusCode int
}

func newHttpResponse(resp *http.Response) Value {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		runtimeExcption(err)
	}
	defer resp.Body.Close()
	obj := &HttpResponse{body: body, cookies: resp.Cookies()}
	obj.status = resp.Status
	obj.statusCode = resp.StatusCode
	obj.headers = resp.Header
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

func (resp *HttpResponse) Save(path string) {
	err := ioutil.WriteFile(path, resp.body, 0666)
	if err != nil {
		runtimeExcption(err)
	}
}

func (resp *HttpResponse) String() string {
	return string(resp.body)
}

func (resp *HttpResponse) Bytes() []byte {
	return resp.body
}

func (resp *HttpResponse) Json() map[string]interface{} {
	var tmp interface{}
	_ = json.Unmarshal(resp.body, &tmp)
	if tmp == nil {
		return nil
	}
	return tmp.(map[string]interface{})
}


