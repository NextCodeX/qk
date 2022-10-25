package main

import (
	_ "embed"
	"fmt"
	"os"
	"qk/core"
	"runtime"
	"time"
)

//go:embed version
var version string

func main() {
	if len(os.Args) > 1 {
		if arg := os.Args[1]; arg == "-v" {
			fmt.Println("Quick version:", version)
			fmt.Println("Build by", runtime.Version())
			return
		}
	}

	// 时区固定为 东八区
	time.Local = time.FixedZone("CST", 8*3600) // 东八

	core.Start()
}
