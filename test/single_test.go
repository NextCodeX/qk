package test

import (
	"fmt"
	"net/url"
	"path/filepath"
	"qk/core"
	"testing"
)

func Test_demo(t *testing.T) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		errorLog(err)
	// 	}
	// }()
	demo, _ := filepath.Abs("../examples/http_client.qk")
	// demo, _ := filepath.Abs("../examples/demo.qk")
	// core.DEBUG = true
	core.TestFlag = true
	core.ExecScriptFile(demo)
}

func Test_api(t *testing.T) {
	urlObj, err := url.Parse("https://127.0.0.1:8080/put/fast?qk=春杰&key=985school&num=雁行")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(urlObj.RawQuery)
	queryVals, err := url.ParseQuery(urlObj.RawQuery)
	fmt.Println(queryVals.Encode())
	fmt.Printf(`%v://%v%v?%v`, urlObj.Scheme, urlObj.Host, urlObj.Path, queryVals.Encode())
}
