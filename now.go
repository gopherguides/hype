package hype

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

var _ ExecutableNode = &Now{}

type Now struct {
	*Element
}

func (now *Now) MarshalJSON() ([]byte, error) {
	if now == nil {
		return nil, ErrIsNil("now")
	}

	m, err := now.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = fmt.Sprintf("%T", now)

	return json.MarshalIndent(m, "", "  ")
}

func (now *Now) Execute(ctx context.Context, doc *Document) error {
	if now == nil {
		return now.WrapErr(ErrIsNil("now"))
	}

	now.Lock()
	defer now.Unlock()

	if doc == nil {
		return now.WrapErr(ErrIsNil("doc"))
	}

	p := doc.Parser

	tm := p.Now()

	fmt, _ := now.Get("gofmt")

	s := now.format(tm, fmt)
	now.Nodes = Nodes{Text(s)}

	return nil
}

func (now *Now) format(tm time.Time, fmt string) string {

	switch fmt {
	case "Layout":
		fmt = time.Layout
	case "ANSIC":
		fmt = time.ANSIC
	case "UnixDate":
		fmt = time.UnixDate
	case "RubyDate":
		fmt = time.RubyDate
	case "RFC822":
		fmt = time.RFC822
	case "RFC822Z":
		fmt = time.RFC822Z
	case "RFC850":
		fmt = time.RFC850
	case "RFC1123":
		fmt = time.RFC1123
	case "RFC1123Z":
		fmt = time.RFC1123Z
	case "RFC3339":
		fmt = time.RFC3339
	case "RFC3339Nano":
		fmt = time.RFC3339Nano
	case "Kitchen":
		fmt = time.Kitchen
	case "Stamp":
		fmt = time.Stamp
	case "StampMilli":
		fmt = time.StampMilli
	case "StampMicro":
		fmt = time.StampMicro
	case "StampNano":
		fmt = time.StampNano
	}

	if len(fmt) == 0 {
		fmt = TIME_FORMAT
	}

	return tm.Format(fmt)
}

func NewNowNodes(p *Parser, el *Element) (Nodes, error) {
	if p == nil {
		return nil, el.WrapErr(ErrIsNil("parser"))
	}

	if el == nil {
		return nil, el.WrapErr(ErrIsNil("element"))
	}

	node := &Now{
		Element: el,
	}

	return Nodes{node}, nil
}
