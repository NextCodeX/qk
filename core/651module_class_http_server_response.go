package core

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ResponseType int
const (
	ResponseText ResponseType = 1 << iota
	ResponseHtml
	ResponseJson
	ResponseFile
	ResponseFileForDownload
	ResponseRedirect
)

type ServerResponse struct {
	t        ResponseType
	headers  map[string]string
	data     []byte
	str      string
	obj      JSONObject
	fileName string
}

func (sr *ServerResponse) Redirect(url string) {
	sr.t = ResponseRedirect
	sr.str = url
}

func (sr *ServerResponse) Txt(raw string) {
	sr.t = ResponseHtml
	sr.str = raw
}

func (sr *ServerResponse) Html(raw string) {
	sr.t = ResponseHtml
	sr.str = raw
}

func (sr *ServerResponse) Json(raw JSONObject) {
	sr.t = ResponseJson
	sr.obj = raw
}

func (sr *ServerResponse) File(bs []byte) {
	sr.t = ResponseFile
	sr.data = bs
}

func (sr *ServerResponse) FileForDownload(fileName string, bs []byte) {
	sr.t = ResponseFileForDownload
	sr.fileName = url.QueryEscape(fileName) // 防止中文乱码
	sr.data = bs
}

func (sr *ServerResponse) Dl(fileName string, bs []byte) {
	sr.FileForDownload(fileName, bs)
}

func (sr *ServerResponse) getData() []byte {
	return sr.data
}

func (sr *ServerResponse) SetHeader(name, value string) {
	if sr.headers == nil {
		sr.headers = make(map[string]string)
	}
	sr.headers[name] = value
}
func (sr *ServerResponse) Seth(name, value string) {
	sr.SetHeader(name, value)
}

func (sr *ServerResponse) SetContentType(value string) {
	if sr.headers == nil {
		sr.headers = make(map[string]string)
	}
	sr.headers["Content-Type"] = value
}
func (sr *ServerResponse) SetType(value string) {
	sr.SetContentType(value)
}

func (sr *ServerResponse) isRedirect() bool {
	return sr.t == ResponseRedirect
}

func (sr *ServerResponse) isText() bool {
	return sr.t == ResponseText
}

func (sr *ServerResponse) isHtml() bool {
	return sr.t == ResponseHtml
}

func (sr *ServerResponse) isJson() bool {
	return sr.t == ResponseJson
}

func (sr *ServerResponse) isFile() bool {
	return sr.t == ResponseFile
}

func (sr *ServerResponse) isFileForDownload() bool {
	return sr.t == ResponseFileForDownload
}






func assembleResponse(w http.ResponseWriter, resp *ServerResponse) {
	if resp.isText() {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = io.WriteString(w, resp.str)
	} else if resp.isHtml() {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, resp.str)
	} else if resp.isJson() {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = io.WriteString(w, resp.obj.String())
	} else if resp.isFile() {
		w.Header().Set("Content-Length", fmt.Sprint(len(resp.data)))
		_, _ = io.Copy(w, bytes.NewReader(resp.data))
	} else if resp.isFileForDownload() {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprint(len(resp.data)))
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%v"`, resp.fileName))
		_, _ = w.Write(resp.data)
	} else if resp.isRedirect() {
		redirect(w, resp.str)
	} else {}
}



func redirect(w http.ResponseWriter, newUrl string) {
	w.Header().Set("Location", newUrl)
	w.WriteHeader(http.StatusPermanentRedirect)
}

func setResponseBody(w http.ResponseWriter, dataType string, data []byte) {
	w.Header().Set("Content-Type", dataType)
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	_, _ = w.Write(data)
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = io.WriteString(w, "404 NOT FOUND")
}


