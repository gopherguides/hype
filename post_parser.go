package hype

type PostParser interface {
	PostParse(p *Parser, d *Document, err error) error
}

type PostParseFn func(p *Parser, d *Document, err error) error

func (fn PostParseFn) PostParse(p *Parser, d *Document, err error) error {
	return fn(p, d, err)
}

func (list Nodes) PostParse(p *Parser, d *Document, err error) error {

	var err2 error

	for _, n := range list {
		if nodes, ok := n.(Nodes); ok {
			err2 = nodes.PostParse(p, d, err)
			if err2 != nil {
				if _, ok := err2.(PostParseError); ok {
					return err2
				}
				return PostParseError{
					Err:      err2,
					Filename: p.Filename,
					OrigErr:  err,
					Root:     p.Root,
				}
			}
			continue
		}

		pp, ok := n.(PostParser)

		if ok {
			err2 = pp.PostParse(p, d, err)
			if err2 != nil {
				return PostParseError{
					Err:        err2,
					Filename:   p.Filename,
					OrigErr:    err,
					PostParser: pp,
					Root:       p.Root,
				}
			}
		}

		err2 = n.Children().PostParse(p, d, err)
		if err2 != nil {
			if _, ok := err2.(PostParseError); ok {
				return err2
			}
			return PostParseError{
				Err:      err2,
				Filename: p.Filename,
				OrigErr:  err,
				Root:     p.Root,
			}
		}
	}

	return err
}
