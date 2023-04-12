# Context Values

As we have seen one feature of the <godoc>context</godoc> package is that it allows you to pass request specific values to the next function in the chain.

This provide a lot of useful benefits, such as passing request or session specific values, such as the request id, user id of the requestor, etc. to the next function in the chain.

Using values, however, has its disadvantages, as we will see shortly.

## Understanding Context Values

The <godoc>context#WithValue</godoc> function can be used to wrap a given <godoc>context#Context</godoc> with a new <godoc>context#Context</godoc> that contains the given key/value pair.

<go doc="context.WithValue"></go>

The <godoc>context#WithValue</godoc> function takes a <godoc>context#Context</godoc> as its first argument, and a key and a value as its second and third arguments.

Both the key and value are `any` values. While this may seem like you can use any type for the key, this is not the case. Like maps, keys must be comparable, so complex types like maps or functions are not allowed.

<go run="main.go" src="src/keys" exit="1" code="main.go#example"></go>

## Key Resolution

When we ask for a key through the <godoc>context#Context.Value</godoc> function, the <godoc>context#Context</godoc> will first check if the key is present in the current <godoc>context#Context</godoc>. If the key is present, the value is returned. If the key is not present, the <godoc>context#Context</godoc> will then check if the key is present in the parent <godoc>context#Context</godoc>. If the key is present, the value is returned. If the key is not present, the <godoc>context#Context</godoc> will then check if the key is present in the <godoc>context#Context</godoc>'s parent's parent, and so on.

Consider the following example. We wrap a <godoc>context#Context</godoc> multiple times with different key/values.

<code src="src/resolution/main.go#example"></code>

From the output we see that the final <godoc>context#Context</godoc> has a parentage that includes all of the values added with <godoc>context#WithValue</godoc>. We can also see that we are able to find all of the keys, including the very first one that we set.

<go run="main.go" src="src/resolution"></go>

<include src="_strings.md" ></include>

<include src="_securing.md" ></include>
