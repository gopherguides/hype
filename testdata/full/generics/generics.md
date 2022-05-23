# Generics

[Generics](https://en.wikipedia.org/wiki/Generic_programming) were first introduced to Go with the release of [Go 1.18](https://go.dev/blog/intro-generics). Go 1.18 was release in March of 2022, during the writing of this {{.Binder.Ident}}. We, like the Go team, have tried our best to present the current idioms and thoughts on the how, what, when, where, and why questions about generics in Go.

## What are Generics?

Generic programming is a programming paradigm that allows us to stub out the implementation of a function with a type that will be provided later. This has benefits for both writing, and using, generic functions. With generics we can write functions that can work with multiple types directly, without having to write the same function multiple times, once for each type. When using generic functions, we can continue use our types as concrete types, instead of interface representations.

## The Problem with Interfaces

Interfaces in Go are a powerful tool that has allowed developers to create some very powerful tools in Go. Interfaces allow us to define a set of methods that describe the behavior of a type. Any type that implements those methods, and behaviors, is considered to implement that interface.

We have already discussed the benefits and drawbacks of interfaces earlier in this {{.Binder.Ident}} so we don't have to re-iterate the benefits of interfaces, but let's discuss some problems with interfaces. For example, let's consider the problem of how to write a function that will return the keys for a given map.

<go src="src/keys/any" sym="Keys"></go>

Go is a statically typed language and so we have to specify the type of the map that we want to get the keys from. A map needs to have both its key and value types specified. We also need to specify the type of slice this function will be returning. In order for this function to support all map types we need to use the `any`, or empty interface, type which will match any type.

While this means we can write a function that we return a list of keys from a map, it also means that this function is difficult to use. Consider a test that tries use a map that isn't of type `map[any]any`. This code fails to compile because the type of map in the test is not compatible with the type of map required by the function.

<go src="src/keys/any" test="-v" code="keys_test.go#example" exit="2"></go>

To solve this problem we would need to create a new, interstitial map of the correct type, and copy all of the keys from the original map into the new map. The same is true of trying to handle the results. We need to loop through the returned slice of keys, asserts the keys are of the correct type, and then copy those values into a new slice of the correct type.

<code src="src/keys/any/keys.go"></code>

<go src="src/keys/fixed" test="-v" code="keys_test.go#example"></go>

While this fixes the tests, it is a very cumbersome way to work with a function such as this. Generics, were designed to help solve exactly this sort of problem.

## Type Constraints

Generics in Go introduced a new concept to the language, called Type Constraints. Type Constraints allow us to specify that a type fits within a certain set of constraints. This is useful when we want to write a function that can work with multiple types, but we want to be able to specify that the function can only work with a specific type.

For example, so far have been using an `int` for the key type in a map, and `string` for the value type. This is fine, but we can use generics to make this more flexible. We may want to use an `int32` or a `float64` for the key type, and `any` value for the value type.

Generics allows us to specify those types as constraints when defining a function or a type. Constraints are added with `[]` after the name of the function or type, but before any parameters.

```go
func Name[constraints](parameters) (returns) {
	// ...
}
```

For example, we can define an `Slicer` function that defines a constraint, type `T`, which can be of `any` type. That new `T` type can then be used in the function signature. In this, the `Slicer` function will return a slice of `T` values.

<go src="src/slicer" sym="Slicer"></go>

When calling the `Slicer` function we can pass any type, and it returns a slice of that same type back.

<go src="src/slicer" test="-v" code="slicer_test.go#example"></go>

In our tests we passed a `string` type to the `Slicer` function. At compile time, sees that we are calling the `Slicer` function with a `string` type and then inserts a function with the appropriate typed signature. For example, by passing a `string` type, the compiler generates a function like the following.

<go src="src/slicer-static" sym="Slicer"></go>

## Multiple Generic Types

With an understanding of the basics of generics, let's go back to our `Keys` function and update it to support generics.

<code src="src/keys/fixed/keys.go#def"></code>

A map has both a key and a value type. We can use generics to specify which types are allowed to be used for both. For example, we can specify that the key type, `K`, must of a type `int`, but the value type, `V`, can be of any type.

<code src="src/keys/generic/start/keys.go#def"></code>

With this change we can pass a map of key type `int` and a value type of `string` to the `Keys` function and it will return a slice of `int` values.

<go src="src/keys/generic/start" test="-v" code="keys_test.go#example"></go>

This however, doesn't work if we want to use a map key of type `string` or `float64`. To do this we will need to specify a bigger set of constraints for the key type.

## Instantiating Generic Functions

When calling a generic function, or creating a new value of a generic type, the Go compiler needs to know which types are being provided for the generic parameters. So far, we have been letting the Go compiler infer the types of the generic parameters based on the types of the values passed in. However, if we were to, instead of call a generic function directly, create a variable pointing to the generic function, the Go compiler will not know the types of the generic parameters and the code will fail to compile.

<go src="src/instantiation/broken" test="-v" code="keys_test.go#example" exit="2"></go>

In these situations we need to provide the compiler with the types of the generic parameters.

<go src="src/instantiation/fixed" test="-v" code="keys_test.go#example"></go>

## Defining Constraints

So far we have been using pretty simple types, such as `int` and `any` for the key and value types. But what if we wanted to use more types than just these? To specify which types can be used for a generic parameter, we can use constraints. Constraints are defined in a similar way to interfaces, but instead of specifying a set of methods, we specify a set of types.

As a start we can define a constraint that requires the type to be an `int`.

<go src="src/constraints/defining" sym="MapKey"></go>

With a `MapKey` constraint defined we can update the `Keys` function to use it instead of `int`.

<code src="src/constraints/defining/keys.go#def"></code>

## Multiple Type Constraints

Currently, the `MapKey` constraint only allows an `int` to be used for the key. If we were to try and use the `Keys` function with a map using a key type of `float64` we would get a compile error.

<go src="src/constraints/floats" test="-v" code="keys_test.go#example" exit="2"></go>

When defining constraints we can use the `|` operator to create an intersection of constraints. For example, we can define a constraint that requires the key type to be either `int` or `float64`.

<go src="src/constraints/or" sym="MapKey"></go>

With this change to the `MapKey` constraint we can now use the `Keys` function with a map using a key type of `float64`.

<go src="src/constraints/or" test="-v"></go>

## Underlying Type Constraints

In Go, we are allowed to create new types based on other types. For example, we can create a new type, `MyInt`, that is based on the `int` type.

<go src="src/constraints/underlying/broken" sym="MyInt"></go>

However, when we try to use the `Keys` function with a map using a key type of `MyInt` we will get a compile error.

<go src="src/constraints/underlying/broken" test="-v" code="keys_test.go#example" exit="2"></go>

The reason for this compilation is that the type `MyInt`, while based on `int`, does not satisfy the `MapKey` constraint because it is **not** an `int` itself. When writing constraints we, usually, are interested in the underlying type, not the type that is wrapped by the type. To express this in when defining a constraint we can use the `~` operator.

<go src="src/constraints/underlying/fixed" sym="MapKey"></go>

By updating the constraint to use the `~` operator we can now use the `Keys` function with a map using a key type of `MyInt`.

<go src="src/constraints/underlying/fixed" test="-v"></go>

## The Constraints Package

When generics were released in Go 1.18, the Go team, decided to be cautious and not update the standard library immediately to use them. They wanted to see how generics were being used before deciding to update the standard library. As a result of this, the Go team have create a series of packages in the <godoc#a>golang.org/x/exp</godoc#a> namespace to experiment with generics. One of these packages is the <godoc#a>golang.org/x/exp/constraints</godoc#a> package. The <godoc#a>golang.org/x/exp/constraints</godoc#a> package defines a set of constraints for all of the numerical, and comparable types in the language.

<go src="src/constraints/pkg" doc="golang.org/x/exp/constraints"></go>

For example, consider the <godoc#a>golang.org/x/exp/constraints#Signed</godoc#a> constraint. This constraint requires that the type be any of the signed integer types defined in the Go language, and any types based on those types.

<go src="src/constraints/pkg" doc="golang.org/x/exp/constraints.Signed"></go>

### The Ordered Constraint

One of the most useful constraints defined in the <godoc#a>golang.org/x/exp/constraints</godoc#a> package is the <godoc#a>golang.org/x/exp/constraints#Ordered</godoc#a> constraint. This constraint list all of the comparable types in the language, and any types based on those types. The <godoc#a>golang.org/x/exp/constraints#Ordered</godoc#a> constraint covers all numerical types and strings.

<go src="src/constraints/pkg" doc="golang.org/x/exp/constraints.Ordered"></go>

The <godoc#a>golang.org/x/exp/constraints#Ordered</godoc#a> constraint is perfect for map keys because all of the types defined in the constraint are comparable. By updating the `Keys` function to use the <godoc#a>golang.org/x/exp/constraints#Ordered</godoc#a> constraint we can now use the `Keys` function with a map using a key type of `string`, or any other type that is comparable.

<code src="src/constraints/pkg/keys.go#def"></code>

<go src="src/constraints/pkg" test="-v" code="keys_test.go#example"></go>

## Type Assertions

When using constraints that are based on types, and not on methods like interfaces, type assertions are not allowed. For example, inside the `Keys` functions we might want to print the key out to the console, but only if it implements the <godoc#a>fmt#Stringer</godoc#a> interface.

<go src="src/assertions/broken" sym="Keys"></go>

With method based interfaces this is possible, but with constraints we can't make this sort of assertion.

<go src="src/assertions/broken" test="-v" exit="2"></go>

As mentioned earlier, at compile time, generic function calls are replaced with their concrete types instead. The result is a `Keys` function that takes a map of `string` to `int` and returns a `[]string`.

<go src="src/assertions/static" sym="Keys"></go>

When looking at the compilation error for "concrete" representation of the `Keys` function the error is a little more clear.

<go src="src/assertions/static" test="-v" exit="2"></go>

In Go type assertions, such as this, against concrete types is not allowed. This is no reason to assert if `string` or `User` or another types implements the interface, because the compiler already if it can be done.

## Mixing Method and Type Constraints

When defining constraints we have to choose between type based constraints and method based constraints. For example, can can't define a constraint that is either <godoc#a>golang.org/x/exp/constraints#Ordered</godoc#a> or <godoc#a>fmt#Stringer</godoc#a>.

<go src="src/assertions/mixed" sym="MapKey"></go>

<go src="src/assertions/mixed" test="-v" exit="2" code="keys.go#def"></go>

## Generic Types

In addition to functions, types can also be generic. If we return to the earlier store example, we could update the `Model` interface to use generics allowing us to use a generic type for the `ID()` method to return.

<go src="src/store" sym="Model"></go>

Now, in order to implement the `Model` interface a type needs a `ID()` method that returns a type listed in the <godoc#a>golang.org/x/exp/constraints#Ordered</godoc#a> constraint.

<go src="src/store" sym="User"></go>

We can update the `Store` type as well to use two constraints, one for the type of map key to be used and the other the `Model` constraint.

<go src="src/store" sym="Store"></go>

When defining methods on types that use generics, the receiver of the method needs to be instantiated with the appropriate concrete type or types. Consider the `Find` method on the `Store` type.

<go src="src/store" sym="Store.Find"></go>

The receiver, `(s Store[K, M])`, is instantiated with the concrete types that the `Store` type was instantiated with. Those types can also be used to define arguments and return values for these methods.

<go src="src/store" test="-v" code="store_test.go#example"></go>
