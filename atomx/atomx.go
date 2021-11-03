package atomx

import (
	adam "github.com/gopherguides/hype/atomx/internal/atom"
	"golang.org/x/net/html/atom"
)

const (
	File      = atom.Atom(adam.File)
	FileGroup = atom.Atom(adam.Filegroup)
	Include   = atom.Atom(adam.Include)
	Page      = atom.Atom(adam.Page)
)
