package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/Reffined/iimplement/extractor"
)

var root string

func main() {
	flag.StringVar(&root, "relPath", "", "reletive path to project's root")
	flag.Parse()
	if root == "" {
		fmt.Println("--relPath is required")
		return
	}

	_, ok := os.LookupEnv("GOFILE")
	if !ok {
		fmt.Println("GOFILE not found")
		return
	}
	rootFs := os.DirFS(root)
	files, _ := fs.Glob(rootFs, "go.mod")
	if len(files) == 0 {
		println("go.mod not found")
		return
	}
	ex := extractor.NewExtractor(root + "/gentest")
	for _, v := range ex.Universe {
		for _, v1 := range v.Types {
			if v1.Kind == "Interface" {
				fmt.Println(v1.Name)
				for _, m := range v1.Methods {
					fmt.Println(m.Signature.Parameters)
				}
			}
		}
	}
}
