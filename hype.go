package hype

import "github.com/gopherguides/hype/atomx"

// DefaultElements returns a map of all the default
// element parser functions.
// For example, `include`, `body`, `code`, etc.
func DefaultElements() map[Atom]ParseElementFn {
	m := map[Atom]ParseElementFn{
		"godoc":          NewGoDocLinkNodes,
		"godoc#a":        NewGoDocLinkNodes,
		"now":            NewNowNodes,
		"toc":            NewToCNodes,
		"youtube":        NewYouTubeNodes,
		atomx.A:          NewLinkNodes,
		atomx.Body:       NewBodyNodes,
		atomx.Cmd:        NewCmdNodes,
		atomx.Code:       NewCodeNodes,
		atomx.Figcaption: NewFigcaptionNodes,
		atomx.Figure:     NewFigureNodes,
		atomx.Go:         NewGolangNodes,
		atomx.Image:      NewImageNodes,
		atomx.Img:        NewImageNodes,
		atomx.Include:    NewIncludeNodes,
		atomx.Li:         NewLINodes,
		atomx.Link:       NewLinkNodes,
		atomx.Metadata:   NewMetadataNodes,
		atomx.Ol:         NewOLNodes,
		atomx.P:          NewParagraphNodes,
		atomx.Page:       NewPageNodes,
		atomx.Ref:        NewRefNodes,
		atomx.Table:      NewTableNodes,
		atomx.Td:         NewTDNodes,
		atomx.Th:         NewTHNodes,
		atomx.Thead:      NewTHeadNodes,
		atomx.Tr:         NewTRNodes,
		atomx.Ul:         NewULNodes,
		atomx.Var:        NewVarNodes,
	}

	for _, h := range atomx.Headings() {
		m[h] = NewHeadingNodes
	}

	return m
}
