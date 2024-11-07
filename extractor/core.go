package extractor

import (
	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"
)

type Extractor struct {
	*parser.Builder
	Universe   types.Universe
	Interfaces map[string]*types.Type
}

func NewExtractor(root string) *Extractor {
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
		}
	}

	return ex
}
