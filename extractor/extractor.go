package extractor

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"regexp"
	"strings"

	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"
)

type Extractor struct {
	*parser.Builder
	Universe          types.Universe
	Interfaces        map[string]*types.Type
	TargetType        *types.Type
	TargetTypeMethods [][]string
}

func NewExtractor(root string, modCache string, targetType string, fileName string) *Extractor {
	ex := new(Extractor)
	ex.Builder = parser.New()
	ex.AddDirRecursive(root)
	ex.AddDirRecursive(modCache)
	ex.gatherModCacheIfaces(modCache)
	u, err := ex.FindTypes()
	if err != nil {
		panic(err)
	}
	ex.Universe = u
	ex.Interfaces = make(map[string]*types.Type)
	for _, i := range u {
		for k, ts := range i.Types {
			if ts.Kind == "Interface" {
				ex.Interfaces[k] = ts
			}
			if ts.Kind == "Struct" {
				if ts.Name.Name == targetType {
					ex.TargetType = ts
				}
			}
		}
	}
	ex.printAllIfaces()
	ex.extractTargetMethods(fileName)
	return ex
}

func (e *Extractor) extractTargetMethods(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	pattern := `(?m)func[\s]*\(.*?[\s](?P<recver>.*?)\)[\s]*(?P<mName>.*?)\(.*?\)[\s]*\{[\s]*.*[\s]*\}`
	reg := regexp.MustCompile(pattern)
	res := reg.FindAllStringSubmatch(string(content), -1)

	e.TargetTypeMethods = res
	return nil
}

func (e *Extractor) printAllIfaces() {
	for _, v := range e.Universe {
		for _, v1 := range v.Types {
			if v1.Kind == "Interface" {
				fmt.Println(v1.Name)
			}
		}
	}
}

func (e *Extractor) gatherModCacheIfaces(path string) {
	dir := os.DirFS(path)
	fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		dirPathParts := strings.Split(path, "/")
		dirPath := strings.Join(dirPathParts[:len(dirPathParts)-1], "/")
		fmt.Println(dirPath)
		return nil
	})
}
