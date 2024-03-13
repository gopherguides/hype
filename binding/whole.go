package binding

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/gobuffalo/flect"
)

type Whole struct {
	Ident     flect.Ident // book
	PartIdent flect.Ident // chapter

	Name  flect.Ident // "My Big Book"
	Parts Parts       // chapters of the book
	Path  string      // path to the whole
}

func (w Whole) String() string {

	mm := map[string]any{
		"ident": w.Ident,
		"name":  w.Name.Titleize(),
		"parts": w.Parts,
		"path":  w.Path,
	}

	b, _ := json.MarshalIndent(mm, "", "  ")
	return string(b)
}

func (w *Whole) UpdatePartIdent(ident flect.Ident) {
	if w == nil {
		return
	}
	w.Parts.UpdateIdent(ident)
}

func WholeFromPath(cab fs.FS, root string, wholeName string, partName string) (*Whole, error) {
	if len(root) == 0 {
		return nil, fmt.Errorf("dir is empty")
	}

	w := &Whole{
		Ident:     flect.New(wholeName),
		Name:      flect.New(filepath.Base(root)),
		PartIdent: flect.New(partName),
		Parts:     Parts{},
		Path:      root,
	}

	err := fs.WalkDir(cab, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		dir := filepath.Dir(path)
		if len(dir) == 0 || dir == "." {
			return nil
		}

		base := filepath.Base(path)
		if base != "hype.md" {
			return nil
		}

		cab, err := fs.Sub(cab, dir)
		if err != nil {
			return fmt.Errorf("error getting sub fs: %q: %w", dir, err)
		}

		part, err := PartFromPath(cab, dir)
		if err != nil {
			return err
		}

		part.Ident = w.PartIdent

		w.Parts[part.Key] = part

		return nil
	})

	if err != nil && !errors.Is(err, ErrPath("")) {
		return w, err
	}

	return w, nil
}
