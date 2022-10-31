package hype

import "fmt"

type Page struct {
	Title string
	*Element
}

func (page *Page) Body() (*Body, error) {
	if page == nil {
		return nil, ErrIsNil("document")
	}

	bodies := ByType[*Body](page.Nodes)

	if len(bodies) == 0 {
		return nil, ErrIsNil("body")
	}

	body := bodies[0]

	return body, nil
}

func (page *Page) PostParse(p *Parser, d *Document, err error) error {
	if err != nil {
		return nil
	}

	if page == nil {
		return ErrIsNil("page")
	}

	if len(page.Title) == 0 {
		page.Title = FindTitle(page.Nodes)
	}

	mds := ByType[*Metadata](page.Children())
	if len(mds) > 1 {
		return fmt.Errorf("page has more than one metadata")
	}

	return nil
}

func NewPage(el *Element) (*Page, error) {
	if el == nil {
		return nil, ErrIsNil("el")
	}

	p := &Page{
		Element: el,
	}

	return p, nil
}

func NewPageNodes(p *Parser, el *Element) (Nodes, error) {
	page, err := NewPage(el)
	if err != nil {
		return nil, err
	}

	return Nodes{page}, nil
}
