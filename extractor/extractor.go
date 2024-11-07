package extractor

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"
)

type Extractor struct {
	*parser.Builder
	Universe          types.Universe
	Interfaces        map[string]*types.Type
	TargetType        *types.Type
	TargetTypeMethods map[string]string
}

func NewExtractor(root string, targetType string, fileName string) *Extractor {
	ex := new(Extractor)
	ex.Builder = parser.New()
	ex.AddDirRecursive(root)
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
	for _, v := range res {
		fmt.Printf("%s:%s:%s\n", v[0], v[1], v[2])
	}
	return nil
}
