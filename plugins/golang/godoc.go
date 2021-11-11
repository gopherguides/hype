package golang

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gopherguides/hype"
	"golang.org/x/net/html"
)

var _ hype.Tag = &Godoc{}
var _ hype.Sourceable = &Godoc{}

type Godoc struct {
	*hype.Node
	cache string
	flags []string
}

func (gd Godoc) Source() (hype.Source, bool) {
	return hype.SrcAttr(gd.Attrs())
}

func (gd Godoc) StartTag() string {
	bb := &bytes.Buffer{}
	fmt.Fprint(bb, gd.Node.StartTag())
	fmt.Fprint(bb, `<pre><code language="godoc" class="language-godoc">`)
	return bb.String()
}

func (gd Godoc) EndTag() string {
	bb := &bytes.Buffer{}
	fmt.Fprint(bb, "</code></pre>")
	fmt.Fprint(bb, gd.Node.EndTag())
	return bb.String()
}

func (gd Godoc) String() string {
	bb := &bytes.Buffer{}
	fmt.Fprint(bb, gd.StartTag())
	fmt.Fprint(bb, gd.Children.String())
	fmt.Fprint(bb, gd.EndTag())
	return bb.String()
}

func (gd *Godoc) Validate(checks ...hype.ValidatorFn) error {
	if gd == nil {
		return fmt.Errorf("godoc is nil")
	}

	_, ok := hype.TagSource(gd)
	if !ok {
		return fmt.Errorf("%s is not a tag source %v", gd.Atom(), gd)
	}

	checks = append(checks, hype.AtomValidator(GODOC))
	return gd.Node.Validate(html.ElementNode, checks...)
}

// github.com/gobuffalo/buffalo.App/Name.short.all.godoc
func (gd Godoc) key(pkg string, flags ...string) string {
	var sep string = string(filepath.Separator)

	fp := strings.ReplaceAll(pkg, "#", sep)
	fp = filepath.Clean(fp)

	for _, flag := range flags {
		flag = strings.ReplaceAll(flag, "-", ".")
		fp += flag
	}
	fp += ".godoc"

	return fp
}

func NewGodoc(n *hype.Node) (*Godoc, error) {
	root, err := CachePath()
	if err != nil {
		return nil, err
	}

	gd := &Godoc{
		Node:  n,
		cache: root,
	}

	if err := gd.Validate(); err != nil {
		return nil, err
	}

	source, ok := gd.Source()
	if !ok {
		return nil, fmt.Errorf("%s is not a sourceable %v", gd.Atom(), gd)
	}

	if f, err := gd.Get("flags"); err == nil {
		gd.flags = CleanFlags(strings.Split(f, ",")...)
	}

	key := gd.key(source.String(), gd.flags...)

	fp := filepath.Join(root, key)

	if b, err := os.ReadFile(fp); err == nil {

		qt := hype.QuickText(string(b))

		gd.Children = append(gd.Children, qt)

		return gd, gd.Validate()
	}

	if err := os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	defer gd.tidy(ctx)
	if err := gd.goGet(ctx, string(source)); err != nil {
		return nil, err
	}

	s, err := gd.Doc(ctx, string(source), gd.flags...)
	if err != nil {
		return nil, err
	}

	gd.Children = append(gd.Children, hype.QuickText(s))
	gd.Children = hype.Tags{
		hype.QuickText(s),
	}

	f, err := os.Create(fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fmt.Fprint(f, s)

	return gd, gd.Validate()
}

func (gd *Godoc) Doc(ctx context.Context, src string, flags ...string) (string, error) {
	if gd == nil {
		return "", fmt.Errorf("doctor is nil")
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	defer gd.tidy(ctx)

	if err := gd.goGet(ctx, src); err != nil {
		return "", err
	}

	args := CleanFlags(flags...)

	pk := strings.ReplaceAll(src, "#", ".")
	pk = strings.TrimSpace(pk)
	args = append(args, pk)

	bb := &bytes.Buffer{}

	r := NewCommander("doc", args...)

	// fmt.Fprintf(bb, "$ %s\n\n", r)

	// v := strings.TrimPrefix(runtime.Version(), "go")
	// fmt.Fprintf(bb, "// Go Version:\t\t%s\n", v)

	// u := fmt.Sprintf("https://pkg.go.dev/%s", src)
	// fmt.Fprintf(bb, "// Documentation:\t<a href=%[1]q target=\"_blank\">%[1]s</a>\n\n", u)

	r.IO = WithOut(r.IO, bb)

	if err := r.Run(ctx, cacheDir); err != nil {
		return "", err
	}

	val := bb.String()
	e := Envelope{
		Body: val,
		Cmd:  r.Cmd(ctx),
		Doc:  src,
	}

	return e.String(), nil
}

// /*
// Usage of [go] doc:
//         go doc
//         go doc <pkg>
//         go doc <sym>[.<methodOrField>]
//         go doc [<pkg>.]<sym>[.<methodOrField>]
//         go doc [<pkg>.][<sym>.]<methodOrField>
//         go doc <pkg> <sym>[.<methodOrField>]
// For more information run
//         go help doc

// Flags:
//   -all
//         show all documentation for package
//   -c    symbol matching honors case (paths not affected)
//   -cmd
//         show symbols with package docs even if package is a command
//   -short
//         one-line representation for each symbol
//   -src
//         show source code for symbol
//   -u    show unexported symbols as well as exported
// */
