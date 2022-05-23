# Implementing io.Writer

Now that `WriteData` uses the <godoc#a>io#Writer</godoc#a> we can, not only use implementations from the standard library like <godoc#a>os#File</godoc#a> and <godoc#a>bytes#Buffer</godoc#a>, we can create our own implementation of <godoc#a>io#Writer</godoc#a>.

<go src="src/func-custom" sym="-all Scribe"></go>

By implementing the `Write` function, with the proper signature, we don't have to implicitly declare our type as an <godoc#a>io#Writer</godoc#a>. The compiler is able to determine whether or not the type being passed in implements the interface being requested.

<go src="src/func-custom-bad" run="." exit="2" code="main.go#bad"></go>

The `*Scribe` type can also be used to test `WriteData` like we did with <godoc#a>bytes#Buffer</godoc#a>.

<go src="src/func-custom" test="-v" code="main_test.go#test"></go>
