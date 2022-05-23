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
				return err2
			}
			continue
		}

		pp, ok := n.(PostParser)

		if ok {
			err2 = pp.PostParse(p, d, err)
			if err2 != nil {
				return PostParseError{
					OrigErr:    err,
					Err:        err2,
					PostParser: pp,
				}
			}
		}

		err2 = n.Children().PostParse(p, d, err)
		if err2 != nil {
			// the error should already be wrapped
			return err2
		}
	}

	return err
}
