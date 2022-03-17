package hype

import (
	"fmt"
	"io/fs"
)

func (p *Parser) ParseFragment(b []byte) (Tags, error) {
	doc, err := p.ParseMD(b)
	if err != nil {
		return nil, err
	}

	pages := doc.Pages()
	if len(pages) == 0 {
		return nil, fmt.Errorf("no pages found")
	}

	page := pages[0]

	kids := page.GetChildren()
	if len(kids) == 0 {
		return nil, fmt.Errorf("no children found")
	}

	res := make(Tags, 0, len(kids))

	for _, kid := range kids {
		if _, ok := kid.(*Text); ok {
			continue
		}

		res = append(res, kid)
	}

	return res, nil
}

func QuickFragment(b []byte, cab fs.FS) (Tags, error) {
	p, err := NewParser(cab)
	if err != nil {
		return nil, err
	}

	return p.ParseFragment(b)
}
