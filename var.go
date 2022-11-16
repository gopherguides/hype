package hype

import (
	"context"
	"fmt"
	"strings"

	"github.com/markbates/syncx"
)

type Var struct {
	*Element
	Data *syncx.Map[string, any]

	value any
}

func (v *Var) String() string {
	const blank = "<var></var>"
	if v == nil {
		return blank
	}

	if v.value != nil {
		return fmt.Sprintf("%v", v.value)
	}

	if v.Element == nil {
		return blank
	}

	s := v.Children().String()
	s = strings.TrimSpace(s)

	return fmt.Sprintf("<var>%v</var>", s)
}

func (v *Var) Execute(ctx context.Context, d *Document) error {
	if v == nil {
		return ErrIsNil("var")
	}

	s := v.Children().String()
	s = strings.TrimSpace(s)

	if IsEmptyNode(v) {
		return v.WrapErr(fmt.Errorf("variable name is empty"))
	}

	val, ok := v.Data.Get(s)
	if !ok {
		return v.WrapErr(fmt.Errorf("variable %q not found", s))
	}

	v.value = val

	return nil
}

func NewVarParserFn(data map[string]any) (ParseElementFn, error) {
	mm := syncx.NewMap(data)
	return func(p *Parser, el *Element) (Nodes, error) {
		if el == nil {
			return nil, ErrIsNil("element")
		}

		v := &Var{
			Element: el,
			Data:    mm,
		}

		s := v.Children().String()
		s = strings.TrimSpace(s)

		v.Nodes = Nodes{Text(s)}

		if IsEmptyNode(v) {
			return nil, v.WrapErr(fmt.Errorf("variable name is empty"))
		}

		if _, ok := v.Data.Get(s); !ok {
			return nil, v.WrapErr(fmt.Errorf("variable %q not found", s))
		}

		return Nodes{v}, nil
	}, nil
}
