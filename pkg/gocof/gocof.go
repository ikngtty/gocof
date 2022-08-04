package gocof

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

func SetSourceDir(dirPath string) {
	err := os.MkdirAll(getGocofDirPath(), 0777)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.WriteFile(getSourceDirFilePath(), []byte(dirPath), 0664)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Execute(dstPath, pkgName string) {
	sourceDir, err := os.ReadFile(getSourceDirFilePath())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	srcPath := path.Join(string(sourceDir), pkgName)

	srcFset := token.NewFileSet()
	srcPkgs, err := parser.ParseDir(srcFset, srcPath, func(fi fs.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dstFileIO, err := os.OpenFile(dstPath, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer dstFileIO.Close()

	for _, srcPkg := range srcPkgs {
		_, err := dstFileIO.Write([]byte(fmt.Sprintf("\n// package %s\n\n", pkgName)))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, srcFile := range srcPkg.Files {
			astutil.Apply(srcFile, nil, func(cur *astutil.Cursor) bool {
				deleteImport(cur)
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

func getGocofDirPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return path.Join(home, ".gocof")
}

func getSourceDirFilePath() string {
	return path.Join(getGocofDirPath(), "source-dir")
}

func deleteImport(cur *astutil.Cursor) bool {
	decl, ok := cur.Node().(*ast.GenDecl)
	if !ok {
		return false
	}
	if decl.Tok != token.IMPORT {
		return false
	}
	cur.Delete()
	return true
}
