package main

import (
	"fmt"
	"os"

	"github.com/ikngtty/gocof/pkg/gocof"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("not specified dstPath and srcPath")
		os.Exit(1)
	}
	gocof.Execute(os.Args[1], os.Args[2])
}
