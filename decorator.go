package hype

import "golang.org/x/net/html"

type Decorator func(p *Parser, tag Tag) (Tag, error)

type Decorators map[*html.Node]Decorator
