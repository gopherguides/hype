package hype

import (
	"bytes"
	"context"
	"encoding/json"
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

	Key   string
	Value any
}

func (v *Var) MarshalJSON() ([]byte, error) {
	if v == nil {
		return nil, ErrIsNil("th")
	}

	v.RLock()
	defer v.RUnlock()

	m, err := v.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = toType(v)
	m["key"] = v.Key
	m["value"] = v.Value

	return json.MarshalIndent(m, "", "  ")
}

func (v *Var) String() string {
	if v == nil || v.Element == nil {
		return "<var></var>"
	}

	if v.Value != nil {
		return fmt.Sprintf("%v", v.Value)
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
	v.Value, ok = doc.Parser.Vars.Get(key)
	if !ok {
		return v.WrapErr(fmt.Errorf("unknown var key %q", key))
	}

	return nil
}

func NewVarNode(p *Parser, el *Element) (*Var, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	if el == nil {
		return nil, ErrIsNil("element")
	}

	key := el.Nodes.String()
	key = strings.TrimSpace(key)

	if len(key) == 0 {
		return nil, fmt.Errorf("missing var key")
	}

	var ok bool
	val, ok := p.Vars.Get(key)
	if !ok {
		return nil, el.WrapErr(fmt.Errorf("unknown var key %q", key))
	}

	v := &Var{
		Element: el,
		Key:     key,
		Value:   val,
	}

	return v, nil
}

func NewVarNodes(p *Parser, el *Element) (Nodes, error) {
	v, err := NewVarNode(p, el)
	if err != nil {
		return nil, err
	}

	return Nodes{v}, nil
}
