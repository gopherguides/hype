# The Empty Interface

All the interfaces we've seen so far have declared one or more methods. In Go there is no **minimum** method count for an interface. That means it is possible to have, what is called, an empty interface. If you declare an interface with zero methods, then every type in the system will be considered to have implemented it.

In Go we use the empty interface to represent "anything".

```go
// generic empty interface:
interface{} // aliased to "any"

// a named empty interface:
type foo interface{}
```

An `int`, for example, has no methods and because of that an `int` will match an interface with no methods.

## The "any" Keyword

In `Go1.18`, [generics](https://go.dev/doc/tutorial/generics) were added to the language. As part of this a new keyword, `any`, was added to the language. This keyword is an alias for `interface{}`.

Using `any` over `interface{}` is a good idea because it is more explicit and it is easier to read.

```go
// Go 1.x:
func foo(x interface{}) {
    // ...
}

// Go 1.18:
func foo(x any) {
    // ...
}
```

If using `Go1.18`, or greater, then you can use the `any` keyword instead of `interface{}`. Using `any` instead of `interface{}` is considered to be idiomatic.

## The Problem With Empty Interfaces

> interface{} says nothing. -- Rob Pike

It is considered bad practice in Go to overuse the empty interface. You should always try to accept either a concrete type or a non-empty interface.

While there are valid reasons to use an empty interface, the downsides should be considered first:

- No type information
- Runtime panics are _very_ possible
- Difficult code (to test, understand, document, etc...)

## Using an Empty Interface

Consider we are writing a data store, similar to a database. We might have an `Insert` method that takes id and the value we want to store. This `Insert` method should be able to our data models. These models might represent users, widgets, orders, etc.

We can use the empty interface to accept all of our models and insert them into the data store.

<code src="src/store/start/store.go" snippet="func"></code>

Unfortunately, this means that in addition to our data models anybody may pass any type to our data store. This is, clearly, not our desire. We could try an set up an elaborate set of `if/else` or `switch` statements, but, this becomes untenable and unmanageable over time. Interfaces allow us to filter out unwanted types and only allow through types that we want.
