package core

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

type HttpFileServer struct {
	serverPort int
	staticResourcePath string
	staticResourceDir string
	dirInitSize       int64
	fileInitSizes     map[string]int64
	mux sync.Mutex
	debugFlag bool
	listenDirFlag bool
}

func newHttpFileServer(basePath, localDir string, serverPort int) *HttpFileServer {
	log.Printf("static resource mapping: %v => %v\n", basePath, localDir)
	return &HttpFileServer{serverPort: serverPort, staticResourceDir: localDir, staticResourcePath: basePath}
}

func (hfs *HttpFileServer) handleStatic(w http.ResponseWriter, req *http.Request) bool {
	targetUri, srcPath, srcDir := req.URL.Path, hfs.staticResourcePath, hfs.staticResourceDir
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
			notFound(w)
			return true
		}

		bs, _ := os.ReadFile(uri)
		bs = hfs.addListenScript(ftype, bs) // 为debug添加监听脚本
		setResponseBody(w, ftype, bs)
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
			redirect(w, realUrl)
			return true
		}

		bs, _ := os.ReadFile(uri)
		bs = hfs.addListenScript("text/html", bs) // 为debug添加监听脚本
		setResponseBody(w,"text/html; charset=utf-8", bs)
		return true
	}

	return false
}

func (hfs *HttpFileServer) localFilePath(urlPath string) string {
	srcPath, srcDir :=  hfs.staticResourcePath, hfs.staticResourceDir
	startIndex := strings.Index(urlPath, srcPath)
	var localPath string
	if startIndex < 0 {
		return ""
	}

	startIndex = startIndex + len(srcPath)
	localPath = urlPath[startIndex:]
	localPath = pathJoin(srcDir, localPath)
	if fileExist(localPath) {
		return localPath
	}

	localPath = pathJoin(localPath, "index.html")
	if fileExist(localPath) {
		return localPath
	}
	return ""
}

func (hfs *HttpFileServer) recordFilesInitSize() {
	hfs.dirInitSize = fileSize(hfs.staticResourceDir)

	hfs.fileInitSizes = make(map[string]int64)
	var fpaths []string
	doScan(hfs.staticResourceDir, false, &fpaths)
	for _, fpath := range fpaths {
		if strings.HasSuffix(fpath, ".html") {
			hfs.fileInitSizes[fpath] = fileSize(fpath)
		}
	}
}

func (hfs *HttpFileServer) resetDirSize(fsize int64) {
	hfs.mux.Lock()
	hfs.dirInitSize = fsize
	hfs.mux.Unlock()
}

func (hfs *HttpFileServer) resetFileSize(fname string, fsize int64) {
	hfs.mux.Lock()
	hfs.fileInitSizes[fname] = fsize
	hfs.mux.Unlock()
}

func (hfs *HttpFileServer) addListenScript(ftype string, bs []byte) []byte {
	if !hfs.debugFlag || !strings.HasPrefix(ftype, "text/html") {
		return bs
	}

	script := `
<script>
    var xhr = new XMLHttpRequest();
    setInterval(function(){
        try {
            var path = window.location.href
            url = "%v?page=" + path
            xhr.open('GET', url);
            xhr.send();
        }catch(e) {
            console.log(e)
        }
    }, 1000);
    var count = 0
    xhr.onreadystatechange = function(){
    　　if ( xhr.readyState == 4 && xhr.status == 200 ) {
            count ++
    　　　　　console.log(count, xhr.responseText)
            if (xhr.responseText == "changed") {
                window.location.reload()
            }
    　　}
    };
</script>
`
	addr := fmt.Sprintf("http://localhost:%v/listenFileState", hfs.serverPort)
	script = fmt.Sprintf(script, addr)
	scriptBytes := []byte(script)
	bs = append(bs, scriptBytes...)
	return bs
}


func (hfs *HttpFileServer) FileStateListener(w http.ResponseWriter, req *http.Request) {
	if hfs.listenDirFlag {
		dirSize := fileSize(hfs.staticResourceDir)
		if hfs.dirInitSize != dirSize {
			setResponseBody(w, "text/plain", []byte("changed"))
			hfs.resetDirSize(dirSize)
		} else {
			setResponseBody(w, "text/plain", []byte("nothing is changed"))
		}
		return
	}

	err := req.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	href, ok := req.Form["page"]
	if !ok {
		setResponseBody(w, "text/plain", []byte("page url is required"))
		return
	}

	parse, err := url.Parse(href[0])
	if err != nil {
		log.Println(err)
		setResponseBody(w, "text/plain", []byte(err.Error()))
		return
	}
	localPath := hfs.localFilePath(parse.Path)
	//fmt.Println("url ->", href, parse.Path, localPath)
	if localPath == "" {
		setResponseBody(w, "text/plain", []byte("page file is not found"))
		return
	}

	size := fileSize(localPath)
	if hfs.fileInitSizes[localPath] != size {
		setResponseBody(w, "text/plain", []byte("changed"))
		hfs.resetFileSize(localPath, size)
	} else {
		setResponseBody(w, "text/plain", []byte("nothing is changed"))
	}
}

