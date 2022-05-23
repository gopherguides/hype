# Concrete Types Vs. Interfaces

Interfaces allow us to specify behavior. They are about doing, not being. Interfaces also allow us to abstract code to make it more reusable, extensible, and more testable.

To illustrate this, let's consider the concept of a performance venue. A performance venue should allow a variety of performers to perform at the venue.

An example of a this as a function might look like the following function.

<go sym="PerformAtVenue" src="src/concrete"></go>

The `PerformAtVenue` function takes a `Musician` as an argument and calls the `Perform` method on the musician. The `Musician` type is a concrete type.

<go sym="Musician" src="src/concrete"></go>

When we pass a `Musician` to the `PerformAtVenue` function, our code compiles and we get the expected output.

<go src="src/concrete" run="." code="main.go#example"></go>

Because the `PerformAtVenue` function takes a `Musician` as an argument, which is a concrete type, we are restricted as to who can perform at the venue. For example, if we were to try to pass a `Poet` to the `PerformAtVenue` function, we would get a compilation error.

<go sym="Poet" src="src/broken"></go>

<go src="src/broken" run="." code="main.go#example" exit="2"></go>

Interfaces allow us to solve this problem by specifying a common set of methods that are required by the `PerformAtVenue` function.

In this example, we can introduce a `Performer` interface. This interface specifies that a `Perform` method is required to implement the `Performer` interface.

<go sym="Performer" src="src/fixed"></go>

Both the `Musician` and `Poet` types have the `Perform` method. Therefore, we can implement the `Performer` interface on both of these types. By updating the `PerformAtVenue` function to take a `Performer` as an argument, we are now able to pass a `Musician` or a `Poet` to the `PerformAtVenue` function.

<go src="src/fixed" sym="PerformAtVenue"></go>
<go src="src/fixed" run="."></go>

By using an interface, instead of a concrete type, we are able to abstract the code and make it more flexible and expandable.
