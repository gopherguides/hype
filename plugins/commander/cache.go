package commander

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/gobuffalo/flect/name"
	"github.com/gopherguides/hype"
)

var cache = &cacher{}

func init() {
	u, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	fp := filepath.Join(u, ".hype", runtime.Version(), "commander")
	os.MkdirAll(fp, 0755)

	cache = &cacher{
		Dir: fp,
	}
}

type cacher struct {
	Dir string
	sync.Mutex
}

func (c *cacher) Retrieve(cmd *Cmd, data Data) error {

	cfp := c.Path(cmd, data)

	c.Lock()
	defer c.Unlock()

	var cf CacheFile

	f, err := os.Open(cfp)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", cfp, err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&cf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("could not decode %s: %w", cfp, err)
	}

	cmd.Children = hype.Tags{hype.QuickText(cf.HTML)}

	// fmt.Println("cache: hit:", cfp, cmd.String())
	return cmd.Validate()
}

func (c *cacher) Store(cmd *Cmd, data Data, res Result) error {
	cfp := c.Path(cmd, data)

	c.Lock()
	defer c.Unlock()

	s, err := res.Out(cmd.Attrs(), data)
	if err != nil {
		return err
	}

	cf := CacheFile{
		Result: res,
		HTML:   s,
	}

	f, err := os.Create(cfp)
	if err != nil {
		return err
	}
	defer f.Close()

	w := io.MultiWriter(f)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	err = enc.Encode(cf)
	if err != nil {
		return err
	}

	// fmt.Println("cache: store:", cfp, cmd.String())
	return nil
}

func (c *cacher) Path(cmd *Cmd, data Data) string {
	c.Lock()
	defer c.Unlock()
	ats := cmd.Attrs()
	src := ats["src"]

	runDir := filepath.Join(src)
	h, _ := hash(runDir)

	cargs := name.File(cmd.Node.StartTag(), ".json")

	cfp := filepath.Join(c.Dir, h, cargs)
	os.MkdirAll(filepath.Dir(cfp), 0755)
	return cfp
}
