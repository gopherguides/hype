# Type Assertion

With a concrete type, like an `int` or a `struct`, the Go compiler and/or runtime knows exactly what the capabilities of that type are. Interfaces, however, can be backed by any type that matches that interface. This means it is possible that the concrete type backing a particular interface provides additional functionality beyond the scope of the interface.

Go allows us to test an interface to see if its concrete implementation is of a certain type. In Go, this is called type assertion.

In this example we are asserting that variable `i`, of type `interface{}` (empty interface), also implements the [`io.Writer`](https://pkg.go.dev/io#Writer) interface. The result of this assertion is assigned to the variable `w`. The variable `w` is of type `io.Writer` and can be used as such.

<code src="src/assertion/bad/assert.go" snippet="def"></code>

<code src="src/assertion/bad/assert.go" snippet="good-assert"></code>

What happens, however, when someone passes a type, such as an `int` or `nil`, that does **not** implement `io.Writer`?

<code src="src/assertion/bad/assert.go" snippet="bad-assert"></code>

The result of a bad assertion is a `panic`.

<go src="src/assertion/bad" run="." exit="1"></go>

These panics can, and will, crash your applications and need to be protected against.

## Asserting Assertion

To prevent a runtime `panic` when a type assertion fails, we can capture a second argument during the assertion. This second variable is of type `bool` and will be `true` if the type assertion succeeded and false if it does not.

<code src="src/assertion/good/assert.go" snippet="def"></code>

You should **ALWAYS** check this boolean to prevent panics and to help keep your applications from crashing.

## Asserting Concrete Types

In addition to asserting that one interface implements another interface, we can use type assertion to get the concrete type underneath.

In this example we trying to assert that variable `w`, of type `io.Writer`, to the type [`*bytes.Buffer`](https://pkg.go.dev/bytes#Buffer). If the assertion is successful, `ok == true`, then variable `bb` will be of type `*bytes.Buffer` and we can now access any publicly exported fields and methods on `*bytes.Buffer`.

<code src="src/assertion/concrete/assert.go" snippet="def"></code>

---

# Assertions Through Switch

When we want to assert an interface for a variety of different types, we can use the `switch` statement in lieu of a lot of `if` statements. Using `switch` statements when doing type assertions also prevents the type assertion panics we saw earlier with individual type assertions.

<code src="src/assertion/switch-basic/assert.go" snippet="def"></code>

## Capturing Switch Type

While just switching on a type can be useful it is often much more useful to capture the result of the type assertion to a variable.

In the following example the result of the type assertion in the switch is assigned to the variable `t := i.(type)`.

<code src="src/assertion/switch-good/assert.go" snippet="def"></code>

In the case of `i` being of type `*bytes.Buffer` then the variable `t` will also be of type `*bytes.Buffer` and all publicly exported fields and methods of `*bytes.Buffer` can now be used.

## Beware of Case Order

The `case` clauses in a `switch` statement are checked in the order that they are listed. A poorly organized `switch` statement can lead to incorrect matches.

In the following example since both `*bytes.Buffer` and `io.WriteStringer` implement `io.Writer`. The first `case` clause matches against `io.Writer` which will match both of those types and prevent the correct clause from being run.

<code src="src/assertion/switch-bad/assert.go" snippet="def"></code>

The [`go-staticcheck`](https://staticcheck.io) tool can be used to check `switch` statements for poor `case` clause organization.

<cmd src="src/assertion/switch-bad" exec="staticcheck" exit="1"></cmd>

---

# Using Assertions

Assertion doesn't just work with empty interfaces. Any interface can be asserted against to see if it implements another interface. We can use this in our data store to add callback hooks; "before insert" and "after insert".

<code src="src/store/validatable/store.go#insert"></code>

## Defining the Callback Interfaces

We can define two new interface types in our system to support before and after insert callbacks.

<code src="src/store/callbacks/callbacks.go" snippet="callbacks"></code>

The `Insert` function can be updated to check for these new interfaces at the appropriate time in the workflow.

<go src="src/store/callbacks" sym="Store.Insert"></go>

These new interfaces allow a type that implements `Validatable` to opt-in to additional functionality.

## Breaking it Down

Let's look at how we are using these interfaces in the `Insert` method.

<code src="src/store/callbacks/store.go" snippet="before"></code>

If the `m` variable, of type `Validatable` (interface), can be asserted to the `BeforeInsertable` interface, then the `bi` variable will be of type `BeforeInsertable` and `ok` will be `true`. The `BeforeInsert` method will be called, the error it returns will be checked, and the application will continue or will return the error. If, however, `m` does not implement `BeforeInsertable` then `ok` will return `false` and the `BeforeInsert` method will never be called.

We check the `AfterInsertable` interface in the same way.

<code src="src/store/callbacks/store.go" snippet="after"></code>
