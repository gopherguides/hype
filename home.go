package hype

import (
	"os"
)

// homeDirectory retursn the home directory of the current user
// it only runs once
var homeDir string

func homeDirectory() (string, error) {
	if homeDir != "" {
		return homeDir, nil
	}

	hd, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	homeDir = hd
	return homeDir, nil
}
