package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	srcFset := token.NewFileSet()
	const srcPath = "/Users/ikngtty/Projects/src/github.com/ikngtty/go-contestlib/math/modular"
	srcPkgs, err := parser.ParseDir(srcFset, srcPath, func(fi fs.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	const dstPath = "/Users/ikngtty/Projects/src/github.com/ikngtty/contest-template/AtCoder/Golang/main.go"
	dstFileIO, err := os.OpenFile(dstPath, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer dstFileIO.Close()

	for _, srcPkg := range srcPkgs {
		for _, srcFile := range srcPkg.Files {
			astutil.Apply(srcFile, nil, func(cur *astutil.Cursor) bool {
				decl, ok := cur.Node().(*ast.GenDecl)
				if !ok {
					return true
				}
				if decl.Tok == token.IMPORT {
					cur.Delete()
				}
				return true
			})

			err := format.Node(dstFileIO, srcFset, srcFile.Decls)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			_, err = dstFileIO.Write([]byte("\n\n"))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
