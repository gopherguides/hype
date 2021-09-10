# UTF-8

Go supports [UTF-8](https://en.wikipedia.org/wiki/UTF-8) characters out of the box, without any special setup, libaries, or packages.

```go
a := "Hello, 世界"
```

## Runes

A rune is a special type in go that represents special characters.

You can define a rune using the single quote (`'`) character:

<code src="src/utf8/utf8-rune/main.go" snippet="main"></code>

If you run the program, it prints out the value of `65`.

The reason for this is because `runes` in Go are a special type.

A `rune` is an alias for `int32`. A rune can be made up of `1` to `3` int32 values.

## Care Must Be Taken

In many languages, the correct way to iterate over a string would look very much like the following code sample...

<code src="src/utf8/utf8-loop/main.go" snippet="main"></code>

But, this will not have the output you would expect:

```plain
0: H
1: e
2: l
3: l
4: o
5: ,
6:
7: ä
8: ¸
9: 
10: ç
11: 
12: 
```

Notice the unexpected characters that were printed out for index 7-12? This is because we were taking part of the rune as an int32, not the entire set of int32's that make up the rune.

## The Right Way

If you intend to walk through each character in a string, the proper way is to use the `range` keyword in the loop.

<code src="src/utf8/utf8-range/main.go" snippet="main"></code>

```text
0: H
1: e
2: l
3: l
4: o
5: ,
6:
7: 世
10: 界
```

Range ensures that we use the proper index and length of int32's to capture the proper rune value.

> Suggested Reading/Viewing:
> [GopherCon 2016: Marcel van Lohuizen - Handling Text from Around the World in Go](https://www.youtube.com/watch?v=K7rMS9Y7_x0)
> [Strings, bytes, runes and characters in Go](https://blog.golang.org/strings)
> [The Rune Type](https://golang.org/doc/go1#rune)
> [Stack Overflow: What is a rune?](https://stackoverflow.com/questions/19310700/what-is-a-rune)

