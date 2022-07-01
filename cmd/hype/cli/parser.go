package cli

import (
	"io/fs"

	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/binding"
)

func NewParser(cab fs.FS, ctxPath string, pwd string) (*hype.Parser, error) {
	p := hype.NewParser(cab)

	p.Section = 1

	if sec, err := binding.PartFromPath(cab, pwd); err == nil {
		p.Section = sec.Number
	}

	if len(ctxPath) > 0 {
		w, err := binding.WholeFromPath(cab, ctxPath, "book", "chapter")
		if err != nil {
			return nil, err
		}

		p.NodeParsers[hype.Atom("binding")] = NewBindingNodes(w)

	}

	return p, nil
}
