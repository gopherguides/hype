# Numbers

Go has two types of numeric types.  The first type is an `architecture independent` type.  This means that regardless of the architecture you compile for, the type will have the correct size (bytes).  The second type is a `implementation` specific type, and the byte size of that numeric type can vary based on the architecture the program is built for.

Go has the following architecture-independent numeric types:

```text
uint8       unsigned  8-bit integers (0 to 255)
uint16      unsigned 16-bit integers (0 to 65535)
uint32      unsigned 32-bit integers (0 to 4294967295)
uint64      unsigned 64-bit integers (0 to 18446744073709551615)
int8        signed  8-bit integers (-128 to 127)
int16       signed 16-bit integers (-32768 to 32767)
int32       signed 32-bit integers (-2147483648 to 2147483647)
int64       signed 64-bit integers (-9223372036854775808 to 9223372036854775807)
float32     IEEE-754 32-bit floating-point numbers (+- 1O-45 -> +- 3.4 * 1038 )
float64     IEEE-754 64-bit floating-point numbers (+- 5 * 10-324 -> 1.7 * 10308 )
complex64   complex numbers with float32 real and imaginary parts
complex128  complex numbers with float64 real and imaginary parts
byte        alias for uint8
rune        alias for int32
```

In addition, Go has the following implementation specific types:

```text
uint     either 32 or 64 bits
int      same size as uint
uintptr  an unsigned integer large enough to store the uninterpreted bits of a pointer value
```

Implementation specific types will have their size defined by the architecture the program is compiled for.

## Picking The Correct Numeric Type

In Go, picking the correct [type](https://golang.org/ref/spec#Types) usually has more to do with performance for the target architecture you are programming for than the size of the data you are working with.  However, without needing to know the specific ramifications of performance for your program, you can follow some of these basic guidelines when first starting out.

For integer data, it's common in Go to use the implementation types like `int` or `uint`.  This will typically result in the fastest processing speed for your target architecture.

If you know you won't exceed a specific size range, then picking an architecture-independent type can both increase speed decrease memory usage.  To understand integer ranges, we can look at the following examples:

```plain
int8 (-128 -> 127)
int16 (-32768 -> 32767)
int32 (− 2,147,483,648 -> 2,147,483,647)
int64 (− 9,223,372,036,854,775,808 -> 9,223,372,036,854,775,807)
```

And for unsigned integers, we have the following ranges:

```plain
uint8 (with alias byte, 0 -> 255)
uint16 (0 -> 65,535)
uint32 (0 -> 4,294,967,295)
uint64 (0 -> 18,446,744,073,709,551,615)
```

For floats:

```plain
float32 (+- 1O-45 -> +- 3.4 * 1038 )
(IEEE-754) float64 (+- 5 * 10-324 -> 1.7 * 10308 )
```

Now that we have looked at some of the possible ranges for numeric data types, we can see how they will be affected if we exceed those ranges in our program.

## Overflow Vs. Wraparound

Go has the potential to both `overflow` a number as well as `wraparound` a number. An `overflow` happens when you try to store a value larger than the data type was designed to store.  When one happens vs. the other is dependent on if the value can be calculated at compile time or at runtime.

At compile time, if the compiler can determine a value will be too large to hold in the data type specified, it will throw an `overflow` error.  This means that the value you calculated to store is too large for the data type you specified.

If we take the following example:

<code src="src/numbers/maxuint8/main.go" snippet="main"></code>

It will compile and run with the following result:

```go
255
```

If we add `1` to the value at runtime, it will wraparound to `0`:

<code src="src/numbers/maxuint8/main.go" snippet="plus"></code>

Output:

```go
0
```

If we change the program to add `1` to the variable when we assign it, it will not compile:

<code src="src/numbers/maxuint8-overflow/main.go" snippet="main">></code>

Because the compiler can determine it will overflow the value it will now throw an error:

```plain
constant 256 overflows uint8
```

Understanding the boundaries of your data will help you avoid potential bugs in your program in the future.

## Saturation

Go does not [saturate](https://en.wikipedia.org/wiki/Saturation_arithmetic) variables during mathematical operations such as addition or multiplication. In languages that saturate, if you had a `uint8` with a max value of `255`, and added `1`, the value would still be the max (saturated) value of `255`.

In go, however, it will always wrap around.  There is no saturation in Go.

<code src="src/numbers/maxuint8-saturation/main.go" snippet="main"></code>

Output:

<code src="src/numbers/maxuint8-saturation/main.go" snippet="output"></code>
