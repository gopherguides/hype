# Problems with String Keys

As is mentioned in the <godoc>context</godoc> documentation using string keys is not recommended. As we just saw when <godoc>context#Context.Value</godoc> tries to resolve a key it finds the first, if any, <godoc>context#Context</godoc> that contains the key and returns that value.

<figure id="context.doc" type="listing">

<go doc="context.WithValue"></go>

<figcaption>Using strings as context keys is not recommended per the documentation.</figcaption>
</figure>

When we use the <godoc>context#Context.Value</godoc> function, we get the last value that was set for the given key. Each time we use <godoc>context#WithValue</godoc> to wrap a <godoc>context#Context</godoc> with a new <godoc>context#Context</godoc>, the new <godoc>context#Context</godoc> will have, essentially, replaced the previous value for the given key.

<img src="assets/string-keys.svg" alt="string-keys">

## Key Collisions

Consider the following example. We wrap a <godoc>context#Context</godoc> multiple times, each time with a different value, but the same key, `request_id`, which is of type `string`.

<code src="src/string-keys/main.go#example"></code>

When we try to log both the `request_id` for both `A` and `A` we see that they are both set to the same value.

<go run="main.go" src="src/string-keys"></go>

One way to solve this problem would be try and "namespace" your `string` keys, `myapp.request_id`. While you may never get into a collision scenario, the possibility of someone else using the same key is there.

## Custom String Key Types

Because Go is a typed language, we can leverage the type system to solve the problem of key collisions. We can create a new type based on `string` that we can use as the key.

<code src="src/custom-keys/main.go#types"></code>

<code src="src/custom-keys/main.go#example"></code>

<code src="src/custom-keys/main.go#logger"></code>

The `Logger` is now properly able to retrieve the two different `request_id` values because they are no longer of the same type.

<go run="main.go" src="src/custom-keys"></go>

This code can be further cleaned up by using constants for the keys that our package, or application, uses. This allows for cleaner code and makes it easier to document the potential keys that may be in a <godoc>context#Context</godoc>.

<code src="src/custom-const/main.go#consts"></code>

<code src="src/custom-const/main.go#logger"></code>
