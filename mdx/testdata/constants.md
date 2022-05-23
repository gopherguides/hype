# Constants

Constants are like variables, except they can't be modified once they have been declared.

Constants can only be a character, string, boolean, or numeric value.

```go
const gopher = "Gopher"
fmt.Println(gopher)
```

Output:

```text
Gopher
```

If you try to modify a constant after it was declared, you will get a compile time error:

<code src="src/constants/const-err/main.go" snippet="main"></code>
Output:

<code src="src/constants/const-err/main.go" snippet="output"></code>

## Untyped

Constants can be `untyped`. This can be useful when working with numbers such as integer type data. If the constant is `untyped`, it is explicitly converted, where `typed` constants are not.

<code src="src/constants/const/main.go"></code>

## Typed

If you declare a constant with a type, it will be that exact type. `leapYear` was defined as data type `int32`. This means it can only operate with `int32` data types. `year` was declared with no type, so it is considered `untyped`. Because of this, you can use it with any integer data type.

<code src="src/constants/const_type/main.go"></code>

If you try to use a `typed` constant with anything other than it's type, Go will throw a compile time error.

## Type Inference

Remember that the `untyped` `const` or `var` will be converted to the `type` it is combined with for any mathematical operation:

<code src="src/constants/const-infer/main.go" snippet="main"></code>

Output:

<code src="src/constants/const-infer/main.go" snippet="output"></code>
