package test

import (
	"fmt"
	"os"
	"path/filepath"
	"qk/core"
	"testing"
)

const COLOR_NORMAL = "\033[0m"
const COLOR_GREEN = "\033[1;32m"
const COLOR_YELLOW = "\033[1;33m"
const COLOR_RED = "\033[1;31m"
const COLOR_GREY = "\033[1;30m"

func TestAllExample(t *testing.T) {
	absDir, err := filepath.Abs("../examples")
	if err != nil {
		return
	}
	entrys, err := os.ReadDir(absDir)
	if err != nil {
		return
	}
	for _, entry := range entrys {
		res := singleFileTest(filepath.Join(absDir, entry.Name()))
		if !res {
			errorLog("failed to test !!!!")
			break
		}
	}
}

func singleFileTest(f string) (res bool) {
	defer func() {
		if err := recover(); err != nil {
			errorLog(fmt.Sprintf("failed to exec: %s\n", f))
			res = false
			panic(err)
		} else {
			fmt.Printf("exec: %s successfully!\n", f)
		}
	}()

	core.TestFlag = true
	core.ExecScriptFile(f)
	res = true
	return
}

func errorLog(format string, args ...any) {
	fmt.Printf("%v%v%v\n", COLOR_RED, fmt.Sprintf(format, args...), COLOR_NORMAL)
}
