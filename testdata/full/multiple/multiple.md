# Multiple Interfaces

Because interfaces are implemented implicitly it means that types can implement many interfaces at once, without explicit declaration. In addition to implementing <godoc#a>io#Writer</godoc#a> the `Scribe` type also implements the <godoc#a>fmt#Stringer</godoc#a> interface.

<go doc="fmt.Stringer"></go>

The <godoc#a>fmt#Stringer</godoc#a> interface is used to convert a value to a string. By implementing a `String() string` method on the `Scribe` type, the `Scribe` now implements both the <godoc#a>fmt#Stringer</godoc#a> and <godoc#a>io#Writer</godoc#a> interfaces.

<go src="src/stringer" sym="Scribe.String"></go>
