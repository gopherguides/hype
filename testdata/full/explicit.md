# Explicit Interface Implementation

In a lot of object oriented languages, such as C# and Java, you have to _explicitly_ declare that your type is implementing a very specific interface.

In C#, <ref id="csharp"></ref>, you declare that you are using the `Performer` interface by using the `:` operator after the class name and listing the interfaces you want to use.

<figure id="csharp">

```c#
interface Performer {
	void Perform();
}

// explicitly implements Performer
class Musician : Performer {
	public void Perform() {}
}

```

<figcaption>Example C# implementation of the `Performer` interface.</figcaption>
</figure>

In Java, <ref id="java"></ref>, you use the `implements` keyword after the class name to tell the compiler that your type wants to implement the `Performer` interface.

<div>

<figure id="java">

```java
interface Performer {
	void Perform();
}

// explicitly implements Performer
class Musician implements Performer {
	void Perform() {}
}
```

<figcaption>Example Java implementation of `Performer` interface.</figcaption>
</figure>

</div>
