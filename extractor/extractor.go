package extractor

import (
	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"
)

type Extractor struct {
	*parser.Builder
	Universe   types.Universe
	Interfaces map[string]*types.Type
	TargetType *types.Type
}

func NewExtractor(root string, targetType string) *Extractor {
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

	return ex
}