# Booleans

The Boolean data type (`bool`) can be one of two values, either **true** or **false**. We use Booleans in programming to make comparisons and to control the flow of the program.

Many operations in math give us answers that evaluate to either True or False:

```go
var ok bool
fmt.Println(ok)
```

Output:

```plain
false
```
In this example we defined a **variable** called `ok` with the data type of `bool`.  When it was printed, we saw the output was `false`.  In Go, all variables have a `zero` value.  In the case of the `bool` data type, the zero value is `false`.

To declare a variable and initialize the value, the `:=` operator can be used:

```go
found := true
fmt.Println(found)
```

Output:

```plain
true
```

The `bool` data type is useful for making logic decisions in your code.
