# Assignment: Week 4

<assignment number="04">

Given the following `Venue` type and interfaces implement an `Entertain` method on `*Venue` that takes the number of audience members, `int`, and a list of `Entertainer`. For each `Entertainer` call its `Perform` method passing in the `Venue`. The `Venue` should check each `Entertainer` to see if it implements the `Setuper` or `Teardowner` interfaces and call them accordingly. The `Venue` should log all `Setup`, `Perform` and `Teardown` calls. Logging should be written to the `Venue.Log` field and use the following formats:

* Setup - `"%s has completed setup.\n"`
* Perform - `"%s has performed for %d people.\n"`
* Teardown - `"%s has completed teardown.\n"`

Write test cases, including error cases, for the provided interfaces by implementing the appropriate interfaces, calling the `Venue.Entertain` method, and checking the logged messages.

You will need to create at least two implementations of the `Entertainer` interfaces. No type should implement more than two of the three interfaces. No type should implement **all** of the interfaces.

<code src="src/venue.go" snippet="venue"></code>

<code src="src/interfaces.go"></code>

</assignment>