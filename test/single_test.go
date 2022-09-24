package test

import (
	"fmt"
	"math"
	"path/filepath"
	"qk/core"
	"runtime"
	"testing"
)

func Test_demo(t *testing.T) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		errorLog(err)
	// 	}
	// }()
	demo, _ := filepath.Abs("../examples/type_str_templ.qk")
	//demo, _ := filepath.Abs("../examples/demo.qk")
	// core.DEBUG = true
	core.TestFlag = true
	core.ExecScriptFile(demo)
}

func Test_api(t *testing.T) {
	fmt.Println(runtime.Version())
	fmt.Println(math.Max(13, 9))
}
