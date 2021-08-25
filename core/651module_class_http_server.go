package core

import (
	"fmt"
	"io"
	"log"
	"net/http"
)


func (fns *InternalFunctionSet) HttpServer(port int) Value {
	obj := &HttpServer{port: port}
	obj.services = make(map[string]Function)
	obj.serviceMethods = make(map[string]string)
	return newClass("HttpServer", &obj)
}

type HttpServer struct {
	serverName     string
	services       map[string]Function
	serviceMethods map[string]string
	fileServer     *HttpFileServer
	port           int
}

func (srv *HttpServer) ServerName(name string) {
	srv.serverName = name
}

func (srv *HttpServer) Debug(args []interface{}) {
	if srv.fileServer != nil {
		srv.fileServer.debugFlag = true
		if len(args) > 0 && args[0].(bool) {
			srv.fileServer.listenDirFlag = true
		}
	}
}

func (srv *HttpServer) StaticDir(basePath , localDir string) {
	srv.fileServer = newHttpFileServer(basePath, localDir, srv.port)
}
func (srv *HttpServer) Static(basePath , localDir string) {
	srv.StaticDir(basePath, localDir)
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
	if srv.fileServer != nil && srv.fileServer.debugFlag {
		srv.fileServer.initFileInfos()
		http.HandleFunc("/listenFileState", srv.fileServer.FileStateListener)
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

		if srv.fileServer != nil {
			// 返回静态资源
			if srv.fileServer.handleStatic(w, req) {
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
		assembleResponse(w, resp)
	})
	addr := fmt.Sprintf(":%v", srv.port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}



