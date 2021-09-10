# Declaring Variables

In Go, there are several ways to declare a variable, and in some cases, more than one way to declare the exact same variable and value.

First, let's declare a variable called `i` of data type `int` without initialization.  This means we will declare a space to put a value, but not give it an initial value.

```go
var i int
```

You now have a variable declared as `i` of data type `int`.

You can also initialize the value by using the equal (`=`) operator:

```go
var i int = 1
```

Both of the above forms are declaration are called `long` variable declarations in Go.

You can also use `short` variable declaration:

```go
i := 1
```

In this case, you have a variable called `i`, and a data type of `int`.  When you don't specify a data type, Go will infer the data type.

With there being three ways to declare variables, the Go community has adopted the following idioms:

Only use long form when you are not initializing the variable:

```go
var i int
```

Use short form when declaring and initializing:

```go
i := 1
```

If you did not desire Go to infer your data type, but you still want to use short variable declaration, you can wrap your value in your desired type.

```go
i := int64(1)
```

The long form is seldom used and not considered idiomatic in Go when you are also initializing the value:

```go
var i int = 1
```

## Zero Values

All built-in types have a zero value. Any allocated variable is usable even if it never has a value assigned. We can see the zero values for the following types with the following code:

<code src="src/variables/zero/main.go">TODO</code>

We used the `%T` verb in the `fmt.Printf` statement. This tells the function to print the `data type` for the variable.

In Go, because all values have a `zero` value, you can't have `undefined` values like some other languages.

For instance, a `boolean` in some languages could be `undefined`, `true`, or `false`, this allowing for `three` states to the variable.  In Go, you can't have more than `two` states for a boolean value.

## Nil

Another type in Go is the `nil` type.  This can represent many concepts.

One of the core concepts for `nil` is that it is the default value for many common types:

* maps
* slices
* functions
* channels
* interfaces
* errors

We'll cover `nil` as it pertains to each type in later chapters.

## Naming Rules

The naming of variables is quite flexible, but there are some rules you need to keep in mind:

* Variable names must only be one word (as in no spaces)
* Variable names must be made up of only letters, numbers, and underscore (`_`)
* Variable names cannot begin with a number

Following the rules above, let’s look at both valid and invalid variable names:

| Valid    | Invalid   | Why Invalid                  |
| -------- | --------- | ---------------------------- |
| userName | user-name | Hyphens are not permitted    |
| name4    | 4name     | Cannot begin with a number   |
| user     | $user     | Cannot use symbols           |
| userName | user name | Cannot be more than one word |

Something else to keep in mind when naming variables, is that they are case-sensitive, meaning that `userName`, `USERNAME`, `UserName`, and `uSERnAME` are all completely different variables. You should avoid using similar variable names within a program to ensure that both you and your current and future collaborators can keep your variables straight.

While variables are case sensitive, the case of the first letter of a variable has special meaning in Go.  If a variable starts with an uppercase letter, then that variable is accessible outside the package it was declared in (or `exported`).  If a variable starts with a lowercase letter, then it is only available within the package it is declared in.

```go
var Email string
var password string
```

`Email` starts with a uppercase letter and can be accessed by other packages.  `password` starts with a lowercase letter, and is only accessible inside the package it is declared in.

## Naming Style

It is common in Go to use very terse (or short) variable names. Given the choice between using `userName` and `user` for a variable, it would be idiomatic to choose `user`.

Scope also plays a role in the terseness of the variable name. The rule is that the smaller the scope the variable exists in, the smaller the variable name.

```go
names := []string{"Mary", "John", "Bob", "Anna"}
for i, n := range names {
	fmt.Printf("index: %d = %q\n", i, n)
}
```

The variable `names` is used in a larger scope, so it would be common to give it a more meaningful name to help remember what it means in the program.  However, the variables `i` and `n` are used immediately in the next line of code, and never used again.  Because of this, someone reading the code will not be confused about where they are used, or what they mean.

Finally, some notes about style. The style is to used `MixedCaps` or `mixedCaps` rather than underscores for multiword names.

| Conventional Style | Unconventional Style | Why Unconventional                       |
| ------------------ | -------------------- | ---------------------------------------- |
| userName           | user_name            | Underscores are not conventional         |
| i                  | index                | prefer `i` over `index` as it is shorter |
| serveHTTP          | serveHttp            | acronyms should be capitalized           |
| userID             | UserId               | acronyms should be capitalized           |

The most important thing about style is to be consistent, and that the team you work on agrees to the style.

## Multiple Assignment

Go also allows you to assign several values to several variables within the same line. Each of these values can be of a different data type:

<code src="src/variables/multiple/main.go" snippet="main">TODO</code>

In the example above, the variable `j` was assigned to the string "gopher", the variable `k` was assigned to the float 2.05, and the variable `l` was assigned to the integer `15`.

This approach to assigning multiple variables to multiple values in one line can keep your lines of code down, but make sure you are not compromising readability for fewer lines of code.
