package main

import (
	_ "embed"
	"fmt"
	"os"
	"qk/core"
	"runtime"
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

	core.Start()
}
