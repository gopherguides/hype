package hype

import (
	"context"
	"fmt"
	"io/fs"
	"strings"
	"sync"

	"github.com/gopherguides/hype/atomx"
	"golang.org/x/sync/errgroup"
)

var _ Node = &Document{}

type Document struct {
	fs.FS
	sync.RWMutex

	Nodes     Nodes
	Parser    *Parser // Parser used to create the document
	Root      string
	SectionID int
	Snippets  *Snippets
	Title     string
}

func (doc *Document) Pages() ([]*Page, error) {
	if doc == nil {
		return nil, ErrIsNil("document")
	}

	pages := ByType[*Page](doc.Nodes)

	if len(pages) == 0 {
		body, err := doc.Body()
		if err != nil {
			return nil, err
		}

		pages = append(pages, body.AsPage())
	}

	return pages, nil
}

func (doc *Document) Body() (*Body, error) {
	if doc == nil {
		return nil, ErrIsNil("document")
	}

	bodies := ByType[*Body](doc.Nodes)

	if len(bodies) == 0 {
		return nil, ErrIsNil("body")
	}

	body := bodies[0]

	return body, nil
}

func (doc *Document) Children() Nodes {
	return doc.Nodes
}

func (doc *Document) Format(f fmt.State, verb rune) {
	if doc == nil {
		return
	}

	switch verb {
	case 'v':
		fmt.Fprintf(f, "%v", doc.Children())
	default:
		fmt.Fprintf(f, "%s", doc.String())
	}
}

func (doc *Document) String() string {
	return doc.Children().String()
}

func (doc *Document) Execute(ctx context.Context) error {
	if doc == nil {
		return ErrIsNil("document")
	}

	err := doc.Children().PreExecute(ctx, doc)
	if err != nil {
		return err
	}

	wg := &errgroup.Group{}

	// execute
	// error gets passed to post executers
	err = doc.Nodes.Execute(wg, ctx, doc)
	if err != nil {
		return err
	}

	err = wg.Wait()

	if err == nil {
		err = doc.processRefs()
	}

	if perr := doc.Children().PostExecute(ctx, doc, err); perr != nil {
		return perr
	}

	return err
}

func (doc *Document) processRefs() error {

	figs := ByType[*Figure](doc.Nodes)

	for i, fig := range figs {

		fig.Lock()
		fig.SectionID = doc.SectionID
		fig.Pos = i + 1
		fig.Unlock()

		caps := ByType[*Figcaption](fig.Nodes)

		if len(caps) > 1 {
			return fmt.Errorf("more than one figcaption")
		}

		fc := &Figcaption{
			Element: NewEl(atomx.Figcaption, fig),
		}

		if len(caps) == 0 {
			return fmt.Errorf("no figcaption: %s", fig.StartTag())
		}

		if len(caps) == 1 {
			fc = caps[0]
		}

		fcb := fc.Nodes.String()
		fcb = strings.TrimSpace(fcb)

		if len(fcb) == 0 {
			return fmt.Errorf("empty figcaption: %s", fig.StartTag())
		}

		em := NewEl(atomx.Em, fc)

		if err := em.Set("class", "figure-name"); err != nil {
			return err
		}

		em.Nodes = append(em.Nodes, Text(fmt.Sprintf("%s:", fig.Name())))

		fcns := fc.Nodes
		fc.Nodes = Nodes{em, Text(" ")}
		fc.Nodes = append(fc.Nodes, fcns...)

	}

	fn := func(i int, fig *Figure) (string, error) {
		return fmt.Sprintf("fig-%d-%d", fig.SectionID, fig.Pos), nil
	}

	return RestripeFigureIDs(doc.Nodes, fn)
}
