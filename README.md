# Don't Believe the Hype! **[Requires Go 1.18+]**

A Package for parsing Markdown and HTML files that allows for custom HTML tags and ships with a lot of fun stuff right out of the box!

## Installation

```bash
$ github.com/gopherguides/hype
```

## Parsing

Before you can parse, you need a parser. Use the `NewParser` function to create a new parser giving it an implementation of `fs.FS`.

```go
p, err := hype.NewParser(os.DirFS("."))
```

## Parsing a File

```go
doc, err := p.ParseFile("README.md")
```

## Parse a Reader

```go
doc, err := p.ParseReader(strings.NewReader(`# Hello World`))
```

## Parse Markdown as Bytes

```go
doc, err := p.ParseMD([]byte(`# Hello World`))
```

## Parse an `html.Node`

If you have a `golang.org/x/net/html.Node` you can parse it with `ParseNode`.

```go
node, err := html.Parse(strings.NewReader(s))
doc, err := p.ParseNode(node)
```

## Custom Tags

You can create your own tag implements and register a factory function with the parser.

```go
type MyTag struct {
    *hype.Node
}

func (t *MyTag) String() string {
    bb := &bytes.Buffer{}
    fmt.Fprintln(bb, t.StartTag())
    fmt.Fprintln(bb, t.GetChildren().String()
    fmt.Fprintln(bb, t.EndTag())
    return bb.String()
}

p , err := hype.NewParser(os.DirFS("."))
p.RegisterTag("mytag", func(n *hype.Node) (hype.Tag, error) {
    return &MyTag{n}, nil
})

in := `<mytag>Hello World</mytag>`

doc, err := p.ParseMD([]byte(in))
```

## Finding Tags

### Finding a Tag by Atom

```go
imgs := ByAtom(doc.Children, "img")
for _, img := range imgs {
    fmt.Println(img.String())
}
```

### Finding a Tag Go Type

```go
imgs := ByType(doc.Children, &hype.Image{})
for _, img := range imgs {
    fmt.Println(img.String())
}
```

## Printing Tags

By default printing a `Document` or a `Tag` will print the complete HTML for that element.

```go
fmt.Println(doc.String()) // prints the entire document
fmt.Println(tag.String()) // prints the tag and its children
```

### Customizing Printing

Printing can be customized by using the `Printer`. By default the `Printer` will print the same HTML as the `String` method.

```go
p, err := NewParser(os.DirFS("testdata"))
if err != nil {
    log.Fatal(err)
}

doc, err := p.ParseFile("files.md")
if err != nil {
    log.Fatal(err)
}

printer := NewPrinter(os.Stdout)
printer.Print(doc)
```

### Setting a Transformer

A transformer function can be set that will be called for each tag before printing. This can be used to modify the tag before printing.

```go
p, err := NewParser(os.DirFS("testdata"))
if err != nil {
    log.Fatal(err)
}

doc, err := p.ParseFile("files.md")
if err != nil {
    log.Fatal(err)
}

printer := NewPrinter(os.Stdout)
printer.SetTransformer(func(t Tag) (Tag, error) {
    if head, ok := t.(*Heading); ok {
        head.Level++
        return head, nil
    }

    // return other tags unchanged
    return t, nil
})

printer.Print(doc)
```
