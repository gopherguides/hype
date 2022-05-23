# Asserting Interface Implementation

Often, especially while implementing an interface, it can be useful to assert your type conforms to all of the interfaces you are trying to implement. One way to do this is to declare a new variable, whose type is the interface you are implementing, and try to assign your type to it. Using the `_` character tells the compiler to do the assignment to the variable and then throw away the result. These assertions are usually done at the package level.

```go
package main

var _ io.Writer = &Scribe{}
var _ fmt.Stringer = Scribe{}
```

The compiler will keep failing until the `Scribe` type implements both the `io.Writer` and `fmt.Stringer` interfaces.
