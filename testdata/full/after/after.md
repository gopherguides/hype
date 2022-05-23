# Using Interfaces

One of the most well known interfaces in Go is the <godoc#a>io#Writer</godoc#a> interface. The <godoc#a>io#Writer</godoc#a> interface requires the implementation of a `Write` method that matches the signature of `Write(p []byte) (n int, err error)`.

<go doc="io.Writer"></go>

Implementations of the <godoc#a>io#Writer</godoc#a> interface can be found all of over the standard library, as well as third party packages. A few of the most common implementations of the <godoc#a>io#Writer</godoc#a> interface are: <godoc#a>os#File</godoc#a>, <godoc#a>bytes#Buffer</godoc#a>, and <godoc#a>strings#Builder</godoc#a>.

Knowing that the only portion of <godoc#a>os#File</godoc#a> we are using matches the <godoc#a>io#Writer</godoc#a> interface we can modify the `WriteData` to use the <godoc#a>io#Writer</godoc#a>, and improve the compatibility and testability of the method.

<go src="src/writer-after" sym="WriteData"></go>

The usage of the `WriteData` function does not change.

<go src="src/writer-after" sym="main"></go>

Testing the `WriteData` function also becomes easier now that we can substitute the implementation with an easier to test implementation.

<go src="src/writer-after" test="-v" code="main_test.go#test"></go>
