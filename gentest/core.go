package gentest

//go:generate iimpl --relPath ..
type Foo struct{}

type (
	IFoo interface {
		Foo(foo string)
	}
	IBar interface {
		IFoo
		Bar()
	}
)
