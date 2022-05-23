# Defining Interfaces

You can create a new interface in the Go type system by using the `type` keyword, giving the new type a name, and then basing that new type on the `interface` type.

```go
type MyInterface interface {}
```

Interfaces define behavior, therefore they are only a collection of methods. Interfaces can have zero, one, or many methods.

> The larger the interface, the weaker the abstraction. -- Rob Pike

It is considered to be non-idiomatic to have large interfaces. Keep the number of methods per interface as small as possible. Small interfaces allow for easier interface implementations, especially when testing. Small interfaces also help us keep our functions and methods small in scope making them more maintainable and testable.

```go
type MyInterface interface {
	Method1()
	Method2() error
	Method3() (string, error)
}
```

It is important to note that interfaces are a collection of methods, not fields. In Go only structs have fields, however, any type in the system can have methods. This is why interfaces are limited to methods only.

```go
// valid
type Writer interface {
	Write(p []byte) (int, error)
}

// invalid
type Emailer interface {
	Email string
}
```

## Defining a Model Interface

Consider, <ref id="insert"></ref>, the `Insert` method for our data store. The method takes two arguments. The first argument is the ID of the model to be stored. The second argument, should be one of our data models, however, because we are using an empty interface, any type from `int` to `nil` may be passed in.

<figure id="insert">

<code src="src/store/start/store.go" snippet="func"></code>

</figure>

<go src="src/store/start" test="-v" code="store_test.go#example"></go>

To prevent types, such as a function definition, that aren't an expected data model, we can define an interface to solve this problem. Since the `Insert` function needs an ID for insertion, we can use that as the basis for an interface.

<go src="src/store/model" sym="Model"></go>

To implement the `Model` interface a type must have a `ID() int` method. We can cleanup the `Insert` method's definition by accepting a single argument, the `Model` interface.

<code src="src/store/model/store.go" snippet="func"></code>

Now, the compiler and/or runtime will reject any type that, such as `string`, `[]byte`, and `func()`, that doesn't have a `ID() int` method.

<go src="src/store/model" test="-v" code="store_test.go#example" exit="2"></go>

## Implementing the Interface

Finally, let's create a new type, `User`, that implements the `Model` interface.

<go src="src/store/user" sym="User"></go>

When we update the tests to use the `User` type, our tests now pass.

<go src="src/store/user" test="-v" code="store_test.go#example"></go>

<ref id="insert"></ref>
