package hype

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gopherguides/hype/atomx"
)

type RefProcessor struct {
	IDGenerator IDGenerator

	mu sync.RWMutex

	indexes map[string]int
}

func (rp *RefProcessor) Process(doc *Document) error {
	if err := rp.validate(); err != nil {
		return err
	}

	figs := ByType[*Figure](doc.Nodes)

	for _, fig := range figs {
		if err := rp.ProcessFigure(doc.SectionID, fig); err != nil {
			return err
		}
	}

	return RestripeFigureIDs(doc.Nodes, rp.IDGenerator)
}

func (rp *RefProcessor) ProcessFigure(sectionID int, fig *Figure) error {
	if err := rp.validate(); err != nil {
		return err
	}

	i := rp.NextIndex(fig.Style())

	fig.Lock()
	fig.SectionID = sectionID
	fig.Pos = i
	fig.Unlock()

	caps := ByType[*Figcaption](fig.Nodes)

	if len(caps) == 0 {
		return fmt.Errorf("no figcaption: %s", fig.StartTag())
	}

	if len(caps) > 1 {
		return fmt.Errorf("more than one figcaption")
	}

	fc := caps[0]

	return rp.processCaption(fig, fc)
}

func (rp *RefProcessor) processCaption(fig *Figure, fc *Figcaption) error {
	fcb := fc.Nodes.String()
	fcb = strings.TrimSpace(fcb)

	if len(fcb) == 0 {
		return fmt.Errorf("empty figcaption: %s", fig.StartTag())
	}

	klass := map[string]string{
		"class": "figure-name",
	}

	if ems := ByAttrs(fc.Nodes, klass); len(ems) > 0 {
		return nil
	}

	em := NewEl(atomx.Em, fc)

	for k, v := range klass {
		if err := em.Set(k, v); err != nil {
			return err
		}
	}

	em.Nodes = append(em.Nodes, Text(fmt.Sprintf("%s:", fig.Name())))

	fcns := fc.Nodes
	fc.Nodes = Nodes{em, Text(" ")}
	fc.Nodes = append(fc.Nodes, fcns...)
	return nil
}

// CurIndex will return the current index for the given key.
func (rp *RefProcessor) CurIndex(key string) int {
	if err := rp.validate(); err != nil {
		return 0
	}

	rp.mu.RLock()
	defer rp.mu.RUnlock()

	if rp.indexes == nil {
		return 0
	}

	return rp.indexes[key]
}

// NextIndex will increment the index for the given key,
// and return the new index.
func (rp *RefProcessor) NextIndex(key string) int {
	if err := rp.validate(); err != nil {
		return 0
	}

	rp.mu.Lock()
	defer rp.mu.Unlock()

	rp.indexes[key]++

	return rp.indexes[key]
}

func (rp *RefProcessor) validate() error {
	if rp == nil {
		return ErrIsNil("RefProcessor")
	}

	rp.mu.Lock()
	defer rp.mu.Unlock()

	if rp.indexes == nil {
		rp.indexes = map[string]int{}
	}

	if rp.IDGenerator == nil {
		rp.IDGenerator = func(i int, fig *Figure) (string, error) {
			return fmt.Sprintf("%s-%d-%d", fig.Style(), fig.SectionID, fig.Pos), nil
		}
	}

	return nil
}
