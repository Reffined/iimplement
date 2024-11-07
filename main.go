package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
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

	ex := extractor.NewExtractor(root + "/gentest")
	i, ok := ex.Interfaces[iface]
	if !ok {
		fmt.Printf("iface %s not found\n", iface)
		return
	}
	txt := strings.Builder{}
	for n, v := range i.Methods {
		toRunes := []rune(strings.ToLower(t))
		recver := fmt.Sprintf("func(%s %s)%s", string(toRunes[0]), t, n)
		args := strings.Builder{}
		args.WriteRune('(')
		for ii := 0; ii < len(v.Signature.ParameterNames); ii++ {
			args.WriteString(v.Signature.ParameterNames[ii])
			args.WriteString(" ")
			args.WriteString(v.Signature.Parameters[ii].String())
			if ii != len(v.Signature.ParameterNames)-1 {
				args.WriteRune(',')
			}
		}
		args.WriteString(")")
		result := strings.Builder{}
		resLen := len(v.Signature.Results)
		if resLen == 1 {
			result.WriteString(v.Signature.Results[0].Name.String())
		} else if resLen > 1 {
			result.WriteRune('(')
			for ii := 0; ii < resLen; ii++ {
				result.WriteString(v.Signature.Results[ii].String())
				if ii != resLen-1 {
					result.WriteRune(',')
				}
			}
			result.WriteRune(')')
		}

		txt.WriteString(fmt.Sprintf("%s%s%s{\n  panic(\"to be implemented\")\n}\n", recver, args.String(), result.String()))
	}
	err := appender.DeleteLastAppend(goFile, t, iface)
	if err != nil {
		panic(err)
	}
}
