package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ikngtty/gocof/pkg/gocof"
)

func main() {
	setSourceDir := flag.String("setSourceDir", "", "set the source directory path")
	flag.Parse()

	if *setSourceDir != "" {
		gocof.SetSourceDir(*setSourceDir)
		return
	}

	if len(os.Args) < 3 {
		fmt.Println("not specified dstPath and pkgName")
		os.Exit(1)
	}
	gocof.Execute(os.Args[1], os.Args[2])
}
