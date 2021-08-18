package hype

import "golang.org/x/net/html"

type Constructor func(p *Parser, node *html.Node) (Tag, error)

type Constructors map[*html.Node]Constructor
