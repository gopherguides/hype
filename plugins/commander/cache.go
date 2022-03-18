package commander

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Cache struct {
	Body      string `db:"body"`
	Command   string `db:"command"`
	Exit      int    `db:"exit"`
	GoVersion string `db:"go_version"`
	Src       string `db:"src"`
	Tag       string `db:"tag"`
	Sum       string `db:"sum"`
}

func (c *Cache) Fetch(db *sqlx.DB) error {
	if db == nil {
		return fmt.Errorf("fetch: db is nil")
	}

	err := db.Get(c, fetchSQL, c.Command, c.Exit, c.GoVersion, c.Src, c.Tag, c.Sum)

	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Insert(db *sqlx.DB) error {
	if db == nil {
		return fmt.Errorf("insert: db is nil")
	}

	m := map[string]any{
		"body":       c.Body,
		"command":    c.Command,
		"exit":       c.Exit,
		"go_version": c.GoVersion,
		"src":        c.Src,
		"tag":        c.Tag,
		"sum":        c.Sum,
	}

	_, err := db.NamedExec(insertSQL, m)

	if err != nil {
		delete(m, "body")
		return fmt.Errorf("insert: %w\n%v", err, m)
	}

	return nil
}
