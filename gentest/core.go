package gentest

//go:generate iimpl --relPath .. --iface IBar --type Bar
type Bar struct{}

// +iipml:Bar:IBar:begin
func (b Bar) Boo() {
	panic("to be implemented")
}
func (b Bar) Foo(foo string, Boo int) (int, error) {
	return 0, nil
}
func (b Bar) Goo(n int) {
	panic("to be implemented")
}
func (b Bar) Roo(Boo string) {
	panic("to be implemented")
}
// +iipml:Bar:IBar:end



































type (
	IFoo interface {
		Foo(foo string, Boo int) (int, error)
		Roo(Boo string)
	}
	IBar interface {
		IFoo
		Boo()
		Goo(n int)
	}
)
