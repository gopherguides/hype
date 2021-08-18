# Strings

A string is a sequence of one or more characters (letters, numbers, symbols) that can be either a constant or a variable. Strings exist within either back quotes **`** or double quotes **"** in Go and have different characteristics depending on which quotes you use.

If you use the back quotes, you are creating a `raw` string literal. If you use the double quotes, you are creating an `interpreted` string literal.

## Raw String Literals

Raw string literals are character sequences between back quotes, often called back ticks. Within the quotes, any character may appear except back quote.

```text
a := `Say "hello" to Go!`
```

Backslashes have no special meaning inside of raw string literals.  Raw string literals may also be used to create multiline strings:

```text
a := `Go is expressive, concise, clean, and efficient.
Its concurrency mechanisms make it easy to write programs
that get the most out of multicore and networked machines,
while its novel type system enables flexible and modular
program construction. Go compiles quickly to machine code
yet has the convenience of garbage collection and the power
of run-time reflection. It's a fast, statically typed,
compiled language that feels like a dynamically typed,
interpreted language.`
```

## Interpreted String Literals

Interpreted string literals are character sequences between double quotes, as in `"bar"`. Within the quotes, any character may appear except newline and unescaped double quote.

```go
a := "Say \"hello\" to Go!"
```

You will almost always use interpreted string literals because they allow for escape characters within them.

> Suggested Reading: [https://blog.golang.org/strings](https://blog.golang.org/strings)
