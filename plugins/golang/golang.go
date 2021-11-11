package golang

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/gopherguides/hype"
)

const cacheDir = ".hype/golang"

func CachePath() (string, error) {
	root, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fp := filepath.Join(root, cacheDir, runtime.Version())
	return fp, nil
}

func Register(p *hype.Parser, root string) {
	p.SetCustomTag(GODOC, func(node *hype.Node) (hype.Tag, error) {
		return NewGodoc(node)
	})

	p.SetCustomTag(LINK, func(node *hype.Node) (hype.Tag, error) {
		return NewLink(node)
	})

	p.SetCustomTag(GORUN, func(node *hype.Node) (hype.Tag, error) {
		return NewGoRun(node, root)
	})
}
