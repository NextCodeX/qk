package core

import (
	"bytes"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)



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
	initSize int64
	mux sync.Mutex
	debugFlag bool
	port int
}

func (srv *HttpServer) ServerName(name string) {
	srv.serverName = name
}

func (srv *HttpServer) Debug() {
	srv.debugFlag = true
}

func (srv *HttpServer) StaticDir(path , localDir string) {
	log.Printf("static resource mapping: %v => %v\n", path, localDir)
	srv.staticResourcePath = path
	srv.staticResourceDir = localDir
}

func formatPath(path string) string {
	if path == "" {
		path = "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	return path
}

func (srv *HttpServer) Get(path string, service Function) {
	path = formatPath(path)
	srv.services[path] = service
	srv.serviceMethods[path] = "GET"
}

func (srv *HttpServer) Post(path string, service Function) {
	path = formatPath(path)
	srv.services[path] = service
	srv.serviceMethods[path] = "POST"
}

func (srv *HttpServer) Any(path string, service Function) {
	path = formatPath(path)
	srv.services[path] = service
	srv.serviceMethods[path] = "ANY"
}

func (srv *HttpServer) Startup() {
	if srv.debugFlag && srv.staticResourceDir != "" {
		srv.initSize = fileSize(srv.staticResourceDir)
		http.Handle("/listenFileState", websocket.Handler(srv.FileStateListener))
		log.Println("Debug Mode is running...")
	}
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			// handle runtime exception
			if err := recover(); err != nil {
				log.Printf("%s\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = io.WriteString(w, "server is in error state")
			}
		}()

		mt := req.Method
		pt := req.URL.Path

		if srv.staticResourcePath != "" {
			// 返回静态资源
			if srv.handleStatic(w, req) {
				return 
			}
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
		if method != "ANY" && method != mt {
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
		srv.assembleResponse(w, resp)

	})
	addr := fmt.Sprintf(":%v", srv.port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (srv *HttpServer) FileStateListener(ws *websocket.Conn) {
	defer ws.Close()
	tiker := time.NewTicker(time.Second)
	for range tiker.C {
		fsize := fileSize(srv.staticResourceDir)
		if fsize != srv.initSize {
			fmt.Printf(srv.staticResourceDir + " -> pre: %v; now: %v\n", srv.initSize, fsize)

			srv.resetFileSize(fsize)
			err := websocket.Message.Send(ws, "changed")
			if err != nil {
				log.Println(err)
			}
			return
		} else {
			err := websocket.Message.Send(ws, "nothing is changed")
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (srv *HttpServer) assembleResponse(w http.ResponseWriter, resp *ServerResponse) {
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
		srv.redirect(w, resp.str)
	} else {}
}

func (srv *HttpServer) handleStatic(w http.ResponseWriter, req *http.Request) bool {
	targetUri, srcPath, srcDir := req.URL.Path, srv.staticResourcePath, srv.staticResourceDir
	startIndex := strings.Index(targetUri, srcPath)
	var uri string
	if startIndex < 0 {
		return false
	} else {
		startIndex = startIndex + len(srcPath)
		uri = targetUri[startIndex:]
		uri = pathJoin(srcDir, uri)
	}

	if fileExist(uri) {
		fileName := fileName(uri)
		ftype := fileType(fileName)
		if ftype == "" {
			srv.notFound(w)
			return true
		}

		bs, _ := os.ReadFile(uri)
		bs = srv.addListenScript(ftype, bs) // 为debug添加监听脚本
		srv.setResponseBody(w, ftype, bs)
		return true
	}

	var dirUrlFlag bool
	var realUrl string
	if strings.LastIndex(targetUri, "/") != len(targetUri) - 1 {
		dirUrlFlag = true
		realUrl = targetUri + "/"
	}

	uri = pathJoin(uri, "index.html")
	if fileExist(uri) {
		if dirUrlFlag {
			// 当url指向的是一个目录，但目录下有index.html文件
			// 我们需要让浏览器重定向至 url+"/",
			// 表示当前页面在当前目录中，防止index.html中使用了相对路径的文件引用出错
			srv.redirect(w, realUrl)
			return true
		}

		bs, _ := os.ReadFile(uri)
		bs = srv.addListenScript("text/html", bs) // 为debug添加监听脚本
		srv.setResponseBody(w,"text/html; charset=utf-8", bs)
		return true
	}
	
	return false
}

func (srv *HttpServer) resetFileSize(fsize int64) {
	srv.mux.Lock()
	srv.initSize = fsize
	srv.mux.Unlock()
}

func (srv *HttpServer) addListenScript(ftype string, bs []byte) []byte {
	if !srv.debugFlag || !strings.HasPrefix(ftype, "text/html") {
		return bs
	}

	script := `
<script>
  var ws = new WebSocket("%v");
  ws.onopen = function(evt) {
    console.log("%v")
  };
  ws.onmessage = function(evt) {
    console.log("Received Message: " + evt.data);
    if (evt.data == "changed") {
        window.document.location.reload()
    }
  };  
</script>
`
	url := fmt.Sprintf("ws://localhost:%v/listenFileState", srv.port)
	msg := fmt.Sprintf(`connetion %v successfully!`, url)
	script = fmt.Sprintf(script, url, msg)
	scriptBytes := []byte(script)
	bs = append(bs, scriptBytes...)
	return bs
}

func (srv *HttpServer) redirect(w http.ResponseWriter, newUrl string) {
	w.Header().Set("Location", newUrl)
	w.WriteHeader(http.StatusPermanentRedirect)
}

func (srv *HttpServer) setResponseBody(w http.ResponseWriter, dataType string, data []byte) {
	w.Header().Set("Content-Type", dataType)
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	w.Write(data)
}

func (srv *HttpServer) notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = io.WriteString(w, "404 NOT FOUND")
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

