package test

import (
	"path/filepath"
	"qk/core"
	"testing"
)

func Test_demo(t *testing.T) {
	demo, _ := filepath.Abs("../examples/type_byte_array.qk")
	core.DEBUG = true
	core.TestFlag = true
	core.ExecScriptFile(demo)
}
