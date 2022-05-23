# Before Interfaces

Consider the following function definition. This function takes a pointer to [`os.File`](https://golang.org/pkg/os/#File), along with a slice of bytes. The function then calls the [`os.File.Write`](https://golang.org/pkg/os/#File.Write) function with the data passed in.

<go src="src/writer-before" sym="WriteData"></go>

In order to call this function we must have an `*os.File`, which is a concrete type in the system. In order to call this function we either need to retrieve, or create, a file on the filesystem, or we can use [`os.Stdout`](https://golang.org/pkg/os/#pkg-variables) which is an `*os.File`.

<go src="src/writer-before" sym="main"></go>

Testing this function involves significant setup. We need to create a new file, call the `WriteData` function, close the file, re-open the file, read the file, and then compare the contents. We need to do all of this work in order to be able to call one function, `Write`, on `*os.File`.

<go src="src/writer-before" test="-v" code="main_test.go#test"></go>

The `WriteData` function is a prime candidate to be refactored using interfaces.
