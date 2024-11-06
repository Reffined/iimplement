package extractor

import (
	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"
)

type Extractor struct {
	*parser.Builder
	Universe types.Universe
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
	return ex
}
