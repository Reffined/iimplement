package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"github.com/Reffined/iimplement/appender"
	"github.com/Reffined/iimplement/extractor"
)

var (
	root  string
	iface string
	t     string
	as    string
)

func main() {
	cmd := exec.Command("go", "env", "GOMODCACHE")
	mod, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	modTrimed := strings.TrimSpace(string(mod))

	flag.StringVar(&root, "relPath", "", "reletive path to project's root")
	flag.StringVar(&iface, "iface", "", "interface to implement")
	flag.StringVar(&t, "type", "", "type to implement iface for")
	flag.StringVar(&as, "as", "pointer", "as 'pointer' receiver or as 'value' receiver")
	flag.Parse()
	if root == "" {
		fmt.Println("--relPath is required")
		return
	}

	goFile, ok := os.LookupEnv("GOFILE")
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
	ex := extractor.NewExtractor(root+"/gentest", modTrimed, t, goFile)
	i, ok := ex.Interfaces[iface]
	if !ok {
		fmt.Printf("iface %s not found\n", iface)
		return
	}
	app := appender.New(i.Methods, ex.TargetTypeMethods)
	err = app.DeleteLastAppend(goFile, t, iface)
	if err != nil {
		panic(err)
	}
	n, err := app.FindEndOfType(goFile, t)
	if err != nil {
		panic(err)
	}
	err = app.Append(goFile, n, t, iface)
	if err != nil {
		panic(err)
	}
}
