package commander

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gopherguides/hype"
)

func fromCache(cmd *Cmd, cfp string, data Data) (*Cmd, error) {
	var cf CacheFile
	f, err := os.Open(cfp)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %w", cfp, err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&cf)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("could not decode %s: %w", cfp, err)
	}

	cmd.Children = append(cmd.Children, hype.QuickText(string(cf.HTML)))
	return cmd, cmd.Validate()
}
