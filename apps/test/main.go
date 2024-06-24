package main

import (
	"fmt"
	"github.com/bytengine-d/go-d/lang"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	fi, err := lang.RealFileInfo("apps/test/main.go")
	if err != nil {
		panic(err)
	}
	fmt.Println(fi.Name())
	fileExt := path.Ext(fi.Name())
	fmt.Println(strings.Replace(fi.Name(), fileExt, "", -1))
	absPath, err := filepath.Abs(fi.Name())
	if err != nil {
		panic(err)
	}
	fmt.Println(path.Dir(absPath))
}
