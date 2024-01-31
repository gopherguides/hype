package hype

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
)

func VarProcessor() PreParseFn {
	fn := func(p *Parser, r io.Reader) (io.Reader, error) {
		if p == nil {
			return nil, ErrIsNil("parser")
		}

		if r == nil {
			return nil, ErrIsNil("reader")
		}

		b, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		lines := strings.Split(string(b), "\n")

		var inside bool

		for _, line := range lines {
			line = strings.TrimSpace(line)

			ll := strings.ToLower(line)

			switch ll {
			case "<details>", "<metadata>", "<p><details>", "<p><metadata>":
				inside = true
			case "</details>", "</metadata>", "</p></details>", "</p></metadata>":
				inside = false
			default:
				if !inside {
					continue
				}

				x := strings.Split(line, ": ")
				if len(x) != 2 {
					return nil, fmt.Errorf("expected key:value got %q", line)
				}

				k := strings.TrimSpace(x[0])
				v := strings.TrimSpace(x[1])

				if _, ok := p.Vars.Get(k); ok {
					continue
				}

				if err := p.Vars.Set(k, v); err != nil {
					return nil, err
				}

			}
		}
		return bytes.NewReader(b), nil
	}

	return fn
}

type Var struct {
	*Element

	value any
}

func (v *Var) String() string {
	if v == nil {
		return ""
	}

	if v.value != nil {
		return fmt.Sprintf("%v", v.value)
	}

	return v.Element.String()
}

func (v *Var) Execute(ctx context.Context, doc *Document) error {
	if v == nil {
		return v.WrapErr(ErrIsNil("var"))
	}

	if v.Element == nil {
		return v.WrapErr(ErrIsNil("element"))
	}

	if doc == nil {
		return v.WrapErr(ErrIsNil("document"))
	}

	key := v.Nodes.String()
	key = strings.TrimSpace(key)

	v.Lock()
	defer v.Unlock()

	var ok bool
	v.value, ok = doc.Parser.Vars.Get(key)
	if !ok {
		return v.WrapErr(fmt.Errorf("unknown var key %q", key))
	}

	return nil
}

func NewVarNode(el *Element) (*Var, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	v := &Var{
		Element: el,
	}

	return v, nil
}

func NewVarNodes(p *Parser, el *Element) (Nodes, error) {
	v, err := NewVarNode(el)
	if err != nil {
		return nil, err
	}

	return Nodes{v}, nil
}
