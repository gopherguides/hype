# Complex Interfaces

## Embedding Interfaces

Interfaces in Go can embed one or more interfaces within itself. This can be used to great effect to combine behaviors into more complex behaviors. The [`io`](https://pkg.go.dev/io) package defines many interfaces, including interfaces that are composed of other interfaces, such as [`io.ReadWriter`](https://pkg.go.dev/io#ReadWriter) and [`io.ReadWriteCloser`](https://pkg.go.dev/io#ReadWriteCloser).

<go doc="io.ReadWriteCloser"></go>

The alternative to embedding other interfaces would to re-declare those same methods in the combined interface.

```go
package io

// ReadWriteCloser is the interface that groups the basic Read, Write and Close methods.
type ReadWriteCloser interface {
  Read(p []byte) (n int, err error)
  Write(p []byte) (n int, err error)
  Close() error
}
```

This is, however, the wrong thing to do. If the intention, as is with `io.ReadWriter`, is implement the `io.Reader` interface, and the `io.Reader` interface changes, then it would no longer implement the correct interface. Embedding the desired interfaces allows us to keep our interfaces cleaner and more resilient.

## Defining an Validatable Interface

Since the act of inserting a model is different than the act of updating a model, we can define an interface to ensure that only types that are **both** a `Model` and have a `Validate() error` method can be inserted.

<go src="src/store/validatable" sym="Validatable"> </go>

The `Validatable` interface embeds the `Model` interface and introduces a new method, `Validate() error`, that must be implemented in addition to the method requirements of the `Model` interface. The `Validate() error` method allows the data model to validate itself before insertion.

<code src="src/store/validatable/store.go" snippet="func"></code>

<go src="src/store/validatable" sym="User.Validate"></go>
