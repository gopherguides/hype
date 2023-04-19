package cli

import (
	"os"
	"path/filepath"
)

func Getwd() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if mp := os.Getenv("MARKED_PATH"); len(mp) > 0 {
		pwd = filepath.Dir(mp)
	}

	return pwd, nil
}
