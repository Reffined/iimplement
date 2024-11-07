package gentest

//go:generate iimpl --relPath .. --iface IFoo --type Bar
type Bar struct{}


type (
	IFoo interface {
		Foo(foo string, Boo int) (error, int)
		Roo(Boo string)
	}
	IBar interface {
		IFoo
		Bar()
	}
)
