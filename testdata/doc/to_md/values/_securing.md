# Securing Context Keys and Values

If we export, make public, the types, and names, of the <godoc>context#Context</godoc> keys our package or application uses, we run the risk of a malicious agent stealing, or modifying our values. For example, in a web request we might set a `request_id` at the beginning of the request, but a piece of middleware later in the chain might modify that value to something else.

<code src="src/malicious/foo/foo.go" snippet="types"></code>

<code src="src/malicious/bar/bar.go" snippet="example"></code>

<code src="src/malicious/main.go" snippet="example"></code>
<code src="src/malicious/foo/foo.go" snippet="example"></code>

<go run="main.go" src="src/malicious"></go>

## Securing by Not Exporting

The best way to secure your that your key/value pairs aren't maliciously overwritten, or accessed, is by not exporting the types, and any constants, used for keys.

<code src="src/secured/foo/foo.go" snippet="types"></code>

Now, you are in control of what values from the context you wish to make public. For example, we can add a helper function to allow others to get access to the `request_id` value.

Because the return value from <godoc>context#Context.Value</godoc> is an empty interface, `interface{}`, we can use these helper functions to, not just retrieve access to the value, but also type assert the value to the type we want, or return an error if it doesn't.

<code src="src/secured/foo/foo.go" snippet="example"></code>

Our application can be updated to use the new helper function to print the `request_id` or exit if there was a problem getting the value.

<code src="src/secured/main.go" snippet="example"></code>

The malicious `bar` package can no longer set, or retrieve, the `request_id` value set by the `foo` package. The `bar` package does not have the ability to create a new type of value `foo.ctxKey` because the type is un-exported can be accessed outside of the `foo` package.

<code src="src/secured/bar/bar.go" snippet="example"></code>

As a result of securing our <godoc>context#Context</godoc> values, the application now correctly retrieves the `request_id` value set by the `foo` package.

<go run="main.go" src="src/secured"></go>
