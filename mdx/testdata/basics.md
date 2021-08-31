# The Go Language Basics - TODO

## Overview of the Go Language - TODO

## Overview of Modules and Packages - TODO

## Basics of Running a Go Program - TODO

## Keywords

The following [keywords](https://golang.org/ref/spec#Keywords) in Go are reserved and may not be used as identifiers.

```text
break
case
chan
const
continue
```

## Operators And Delimiters

The following character sequences represent [operators](https://golang.org/ref/spec#Operators) (including [assignment operators](https://golang.org/ref/spec#assign_op)) and punctuation:

```text
+    &     +=    &=     &&    ==    !=    (    )
```

## Statically Typed

Go is a statically typed language.  Statically typed means that each statement in the program is checked at compile time. It also means that the data type is bound to the variable, whereas in Dynamically linked languages, the data type is bound to the value.

For example, in Go, the type is declared when declaring a variable:

```go
var pi float64 = 3.14
var week int = 7
```

As apposed to a language like PHP, where the data type is associated to the value:

```php
$s = "gopher";        // $s is a string
$s = 123;             // $s is now an integer
```

Now that we know a little about a statically typed language, this will help us better understand how numbers in Go work.

