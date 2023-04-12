package hype

import (
	"fmt"
	"strings"

	"github.com/gopherguides/hype/atomx"
	"github.com/markbates/table"
)

type Table struct {
	*Element
}

func (tab *Table) IsEmptyNode() bool {
	if tab == nil {
		return true
	}

	data, err := tab.Data()

	if err != nil {
		return false
	}

	rows, err := data.Rows()
	if err != nil {
		return false
	}

	return len(rows) == 0
}

func (tab *Table) Data() (*table.Table, error) {
	if tab == nil {
		return nil, ErrIsNil("table")
	}

	if tab.Element == nil {
		return nil, ErrIsNil("table.Element")
	}

	res := &table.Table{}

	if err := tab.setColumns(res); err != nil {
		return nil, err
	}

	if err := tab.setData(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (tab *Table) setColumns(res *table.Table) error {
	var cols []string

	heads := ByAtom(tab.Children(), atomx.Thead)
	if len(heads) == 0 {
		return nil
	}

	head := heads[0]
	for _, th := range ByAtom(head.Children(), atomx.Th) {
		s := fmt.Sprintf("%s", th.Children())
		s = strings.TrimSpace(s)
		cols = append(cols, s)
	}

	if err := res.SetColumns(cols...); err != nil {
		return err
	}

	return nil
}

func (tab *Table) setData(res *table.Table) error {
	bodies := ByAtom(tab.Children(), atomx.Tbody)
	if len(bodies) == 0 {
		return nil
	}

	body := bodies[0]
	for _, tr := range ByAtom(body.Children(), atomx.Tr) {
		var row []any
		for _, td := range ByAtom(tr.Children(), atomx.Td) {
			s := fmt.Sprintf("%s", td.Children())
			s = strings.TrimSpace(s)
			row = append(row, s)
		}

		if err := res.QuickRow(row...); err != nil {
			return err
		}
	}

	return nil
}

func NewTable(p *Parser, el *Element) (*Table, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	t := &Table{
		Element: el,
	}

	return t, nil
}

func NewTableNodes(p *Parser, el *Element) (Nodes, error) {
	h, err := NewTable(p, el)
	if err != nil {
		return nil, err
	}

	return Nodes{h}, nil
}

func (tab *Table) MD() string {
	if tab == nil {
		return ""
	}

	bb := &strings.Builder{}

	heads := ByAtom(tab.Children(), atomx.Th)
	cols := make([]string, 0, len(heads))
	if len(heads) > 0 {
		for _, h := range heads {
			c := h.Children().MD()
			cols = append(cols, c)
			fmt.Fprintf(bb, "| %s ", c)
		}
		bb.WriteString("|\n")
	}

	if len(cols) > 0 {
		for _, h := range cols {
			fmt.Fprintf(bb, "| %s ", strings.Repeat("-", len(h)))
		}
		bb.WriteString("|\n")
	}

	rows := ByAtom(tab.Children(), atomx.Tr)
	for _, r := range rows {
		tds := ByAtom(r.Children(), atomx.Td)
		if len(tds) == 0 {
			continue
		}

		for _, td := range tds {
			fmt.Fprintf(bb, "| %s ", td.Children().MD())
		}
		bb.WriteString("|\n")
	}

	return bb.String()
}
