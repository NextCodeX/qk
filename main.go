package main

import (
	"fmt"
	"os"
	"qk/core"
)

func main() {
	if len(os.Args) > 1 {
		if arg := os.Args[1]; arg == "-v" {
			fmt.Println("Quick version:", version)
			return
		}
	}

	core.Start()
}
