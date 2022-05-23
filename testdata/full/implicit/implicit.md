# Implicit Interface Implementation

In Go interfaces are implemented _implicitly_. This means we don't have indicate to Go that we are implementing an interface. Given a `Performer` interface, a type would need to implement the `Perform` method to be considered a performer.

<go sym="Performer" src="src/implicit"></go>

By adding a `Perform` method, that matches the signature of the `Performer` interface, the `Musician` type is now implicitly implementing the `Performer` interface.

<go sym="-all Musician" src="src/implicit"></go>

Provided a type implements all behaviors specified in the interface, it can be said to implement that interface. The compiler will check to make sure a type is acceptable and will report an error if it does not. Sometimes this is called **Duck Typing**, but since it happens at compile-time in Go, it is called **Structural typing**.

<go run="." src="src/implicit" code="main.go#example"></go>

Structural typing has a handle of useful side-effects.

- The concrete type does not need to know about your interface
- You are able to write interfaces for concrete types that already exist
- You can write interfaces for other people's types, or types that appear in other packages

> Suggested Reading: [Duck typing](https://en.wikipedia.org/wiki/Duck_typing)

> Suggested Reading: [Structural typing](https://en.wikipedia.org/wiki/Structural_type_system)
