package hype

import (
	"context"
	"fmt"
	"io/fs"
	"sync"

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
	if doc == nil {
		return nil
	}

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

// Execute the Document with the given context.
// Any child nodes that implement the PreExecuter,
// ExecutableNode, or PostExecuter interfaces will be executed.
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

	if err := doc.processRefs(err); err != nil {
		return err
	}

	if perr := doc.Children().PostExecute(ctx, doc, err); perr != nil {
		return perr
	}

	return err
}

func (doc *Document) processRefs(err error) error {
	if err != nil {
		return nil
	}

	if doc == nil {
		return ErrIsNil("document")
	}

	rp := &RefProcessor{}

	err = rp.Process(doc)
	if err != nil {
		return err
	}

	return nil
}

type Documents []*Document

func (docs Documents) Execute(ctx context.Context) error {
	if docs == nil {
		return ErrIsNil("documents")
	}

	wg := &errgroup.Group{}

	for _, doc := range docs {
		doc := doc
		wg.Go(func() error {
			return doc.Execute(ctx)
		})
	}

	return wg.Wait()
}
