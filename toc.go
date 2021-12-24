package hype

import (
	"bytes"
	"fmt"

	"github.com/gopherguides/hype/atomx"
)

type Tocs []*Toc

// Toc represents a table of contents.
type Toc struct {
	Children Tocs   // Children of this Toc.
	Indent   int    // Indentation level.
	Title    string // Title of this Toc.
}

func (t Toc) String() string {
	bb := &bytes.Buffer{}
	for i := 0; i < t.Indent; i++ {
		bb.WriteString("  ")
	}

	fmt.Fprintln(bb, t.Title)
	for _, child := range t.Children {
		fmt.Fprint(bb, child.String())
	}
	return bb.String()
}

// Tocs returns a table of contents for the given documents.
func TocsFromDocs(title string, docs ...*Document) (*Toc, error) {
	toc := &Toc{
		Title: title,
	}

	for _, doc := range docs {

		dt := &Toc{
			Title:  doc.Title(),
			Indent: 1,
		}

		tocKids(dt, doc.Children.ByAtom(atomx.Page))
		toc.Children = append(toc.Children, dt)
	}

	return toc, nil
}

func tocKids(parent *Toc, tags Tags) {

	tos := func(tag Tag) string {
		return tag.GetChildren().String()
	}

	for _, tag := range tags {
		toc := &Toc{
			Title:  tos(tag),
			Indent: parent.Indent + 1,
		}

		switch tag.Atom() {
		case atomx.H1:
			// toc.Indent += 1
			tocKids(toc, tag.GetChildren())
			parent.Children = append(parent.Children, toc)
		case atomx.H2:
			toc.Indent += 1
			tocKids(toc, tag.GetChildren())
			parent.Children = append(parent.Children, toc)
		case atomx.H3:
			toc.Indent += 2
			tocKids(toc, tag.GetChildren())
			parent.Children = append(parent.Children, toc)
		case atomx.H4:
			toc.Indent += 3
			tocKids(toc, tag.GetChildren())
			parent.Children = append(parent.Children, toc)
		case atomx.H5:
			toc.Indent += 4
			tocKids(toc, tag.GetChildren())
			parent.Children = append(parent.Children, toc)
		case atomx.Page:
			p, ok := tag.(*Page)
			if ok {
				toc.Title = p.Title()
			}

			tocKids(parent, tag.GetChildren())
		}
	}
}
