package binding

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gobuffalo/flect"
)

type Part struct {
	Ident flect.Ident // chapter

	Name   flect.Ident // "Arrays and Slices"
	Key    string      // arrays-and-slices
	Number int         // 9
	Path   string      // path to the part
}

func (part Part) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"ident":  part.Ident,
		"name":   part.Name.Titleize(),
		"number": part.Number,
		"path":   part.Path,
	}

	return json.MarshalIndent(mm, "", "  ")
}

// func (s Part) String() string {
// 	return fmt.Sprintf("\"%s %d: %s\"", s.Ident.Titleize(), s.Number, s.Name.Titleize())
// }

type Parts map[string]Part

func (parts Parts) UpdateIdent(ident flect.Ident) {
	if parts == nil {
		return
	}

	for _, part := range parts {
		part.Ident = ident
	}
}

func PartFromPath(cab fs.FS, mp string) (Part, error) {
	part := Part{
		Path: mp,
	}

	if len(mp) == 0 {
		return part, fmt.Errorf("dir is empty")
	}

	base := filepath.Base(mp)
	ext := filepath.Ext(base)

	if len(ext) > 0 {
		base = filepath.Base(filepath.Dir(mp))
	}

	rx, err := regexp.Compile(`^(\d+)-(.+)`)

	if err != nil {
		return part, err
	}

	matches := rx.FindAllStringSubmatch(base, -1)
	if len(matches) < 1 {
		return part, ErrPath(mp)
	}

	match := matches[0]
	if len(match) < 3 {
		return part, ErrPath(mp)
	}

	id, err := strconv.Atoi(match[1])
	if err != nil {
		return part, ErrPath(mp)
	}

	part.Number = id

	name := match[2]
	name = strings.TrimSuffix(name, filepath.Ext(name))

	part.Key = name

	part.Name = flect.New(name)

	fp := filepath.Join(mp, "module.md")

	if _, err := fs.Stat(cab, fp); err != nil {
		return part, nil
	}

	return part, nil
}
