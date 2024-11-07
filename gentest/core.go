package gentest

//go:generate iimpl --relPath .. --iface IFoo --type Bar
type Bar struct{}

// +iipml:Bar:IFoo:begin
func (b Bar) Foo(foo string, Boo int) (error, int) {
	panic("to be implemented")
}

func (b Bar) Roo(Boo string) {
	panic("to be implemented")
}

// +iipml:Bar:IFoo:end

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
