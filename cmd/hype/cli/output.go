package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type OutPath struct {
	val *string
}

func (o *OutPath) String() string {
	if o.val == nil {
		return ""
	}
	return o.Value()
}

func (o *OutPath) Value() string {
	return *o.val
}

func (o *OutPath) Exists() bool {
	return o.val != nil
}

func dirExists(path string) error {
	pInfo, err := os.Stat(path)

	// path exists
	if err == nil {
		// path exists but it's a directory
		if pInfo.IsDir() {
			return fmt.Errorf("expected a file path but received path to a directory: %s. ", path)
		}
		return nil
	}

	// path doesn't exist, check if the dir exists
	dir := filepath.Dir(path)
	if _, err = os.Stat(dir); err != nil {
		return err
	}
	return nil
}

func cleanPath(p string) (string, error) {
	if filepath.IsAbs(p) {
		return filepath.Clean(p), nil
	}

	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to identify current working directory: %s", err.Error())
	}
	filepath.Clean(p)
	fullPath := filepath.Join(pwd, p)
	return fullPath, nil
}

func (o *OutPath) Set(path string) error {
	path = strings.TrimSpace(path)

	// path not set
	if len(path) == 0 {
		return nil
	}

	fullPath, err := cleanPath(path)
	if err != nil {
		return err
	}

	err = dirExists(fullPath)
	if err != nil {
		return err
	}

	o.val = &fullPath
	return nil
}
