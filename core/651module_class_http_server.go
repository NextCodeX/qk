package core

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

var mimes = map[string]string{
	".css":  "text/css; charset=utf-8",
	".gif":  "image/gif",
	".htm":  "text/html; charset=utf-8",
	".html": "text/html; charset=utf-8",
	".jpeg": "image/jpeg",
	".jpg":  "image/jpeg",
	".js":   "text/javascript; charset=utf-8",
	".json": "application/json",
	".mjs":  "text/javascript; charset=utf-8",
	".pdf":  "application/pdf",
	".png":  "image/png",
	".svg":  "image/svg+xml",
	".wasm": "application/wasm",
	".webp": "image/webp",
	".xml":  "text/xml; charset=utf-8",
}

func (fns *InternalFunctionSet) HttpServer(port int) Value {
	obj := &HttpServer{port: port}
	obj.services = make(map[string]Function)
	obj.serviceMethods = make(map[string]string)
	return newClass("HttpServer", &obj)
}

type HttpServer struct {
	serverName string
	services map[string]Function
	serviceMethods map[string]string
	staticResourcePath string
	staticResourceDir string
	port int
}

func (srv *HttpServer) ServerName(name string) {

	srv.serverName = name
}

func (srv *HttpServer) StaticDir(path , localDir string) {
	srv.staticResourcePath = path
	srv.staticResourceDir = localDir
}

func (srv *HttpServer) Get(path string, service Function) {
	if path == "" {
		path = "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	srv.services[path] = service
	srv.serviceMethods[path] = "GET"
}

func (srv *HttpServer) Post(path string, service Function) {
	if path == "" {
		path = "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	srv.services[path] = service
	srv.serviceMethods[path] = "POST"
}

func (srv *HttpServer) Startup() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			// handle runtime exception
			if err := recover(); err != nil {
				fmt.Printf("%s\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = io.WriteString(w, "server is in error state")
			}
		}()

		mt := req.Method
		pt := req.URL.Path

		if srv.staticResourcePath != "" {
			// 返回静态资源

		}

		if srv.serverName != "" {
			pt = "/" + srv.serverName + pt
		}
		method, ok := srv.serviceMethods[pt]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, "404 NOT FOUND")
			return
		}
		if method != mt {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = io.WriteString(w, "405 Method Not Allowed")
			return
		}


		resp := &ServerResponse{}
		request := newHttpRequest(req, mt)
		response := newClass("ServerResponse", &resp)

		// run service
		service := srv.services[pt]
		var args []Value
		args = append(args, request)
		args = append(args, response)
		service.setArgs(args)
		service.execute()

		if resp.headers != nil {
			// set response header
			for k, v := range resp.headers {
				w.Header().Set(k, v)
			}
		}

		// set response body
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
			w.Header().Set("Location", resp.str)
			w.WriteHeader(http.StatusTemporaryRedirect)
		}

	})
	addr := fmt.Sprintf(":%v", srv.port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

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

