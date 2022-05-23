package demo

// Foo is a foo.
func Foo() string {
	return "foo"
}

type bar struct {
	Foo string
}

func (b *bar) Foo() string {
	return b.Foo
}
