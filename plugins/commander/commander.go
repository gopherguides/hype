package commander

import (
	"fmt"

	_ "embed"

	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/atomx"
	"github.com/jmoiron/sqlx"
)

const (
	CMD     atomx.Atom = "cmd"
	COMMAND atomx.Atom = "command"
)

type Attributes = hype.Attributes
type Data map[string]string

func CustomTag(p *hype.Parser) (hype.CustomTagFn, error) {
	fn := func(node *hype.Node) (hype.Tag, error) {
		return NewCmd(node)
	}

	return fn, nil
}

// Register registers all of the atoms and tags
// that this plugin provides.
func Register(p *hype.Parser) error {
	if p == nil {
		return fmt.Errorf("parser is nil")
	}

	if err := migrate(p.DB); err != nil {
		return err
	}

	fn, err := CustomTag(p)
	if err != nil {
		return err
	}

	p.SetCustomTag(CMD, fn)
	p.SetCustomTag(COMMAND, fn)

	return nil
}

func migrate(db *sqlx.DB) error {
	if db == nil {
		return nil
	}

	_, err := db.Exec(migrationSQL)
	if err != nil {
		return err
	}

	return nil
}

// cmds is a map of special data attributes for different
// exacutable commands.
var cmds = map[string]hype.Attributes{
	"tree": hype.Attributes{
		"hide-data": "true",
	},
	"cat": hype.Attributes{
		"hide-data": "true",
	},
}
