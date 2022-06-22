package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gobuffalo/flect"
	"github.com/gopherguides/hype"
)

type Part struct {
	flect.Ident // chapter

	Name   flect.Ident // "Arrays and Slices"
	Key    string      // arrays-and-slices
	Number int         // 9
}

func (p Part) Children() hype.Nodes {
	return nil
}

func (part Part) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"ident":  part.Ident,
		"name":   part.Name.Titleize(),
		"number": part.Number,
	}

	return json.MarshalIndent(mm, "", "  ")
}

func (s Part) String() string {
	return fmt.Sprintf("\"%s %d: %s\"", s.Ident.Titleize(), s.Number, s.Name.Titleize())
}

type Parts map[string]*Part

func (parts Parts) UpdateIdent(ident flect.Ident) {
	if parts == nil {
		return
	}

	for _, part := range parts {
		if part == nil {
			continue
		}
		part.Ident = ident
	}
}

func PartFromPath(mp string) (*Part, error) {
	if len(mp) == 0 {
		return nil, fmt.Errorf("dir is empty")
	}

	base := filepath.Base(mp)
	ext := filepath.Ext(base)

	if len(ext) > 0 {
		base = filepath.Base(filepath.Dir(mp))
	}

	rx, err := regexp.Compile(`^(\d+)-(.+)`)

	if err != nil {
		return nil, err
	}

	matches := rx.FindAllStringSubmatch(base, -1)
	if len(matches) < 1 {
		return nil, sectionPathError(mp)
	}

	match := matches[0]
	if len(match) < 3 {
		return nil, sectionPathError(mp)
	}

	id, err := strconv.Atoi(match[1])
	if err != nil {
		return nil, sectionPathError(mp)
	}

	part := &Part{
		Number: id,
	}

	name := match[2]
	name = strings.TrimSuffix(name, filepath.Ext(name))

	part.Key = name

	part.Name = flect.New(name)

	fp := filepath.Join(mp, "module.md")

	if _, err := os.Stat(fp); err != nil {
		return part, nil
	}

	dir := filepath.Dir(fp)
	p := hype.NewParser(os.DirFS(dir))

	doc, err := p.ParseFile("module.md")
	if err != nil {
		return nil, err
	}

	part.Name = flect.New(doc.Title)

	return part, nil
}

type sectionPathError string

func (e sectionPathError) Error() string {
	return fmt.Sprintf("could not parse section from: %q", string(e))
}
