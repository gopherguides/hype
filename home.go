package hype

import (
	"log"
	"os"
)

// homeDirectory retursn the home directory of the current user
// it only runs once
var homeDir string

func homeDirectory() string {
	if homeDir != "" {
		return homeDir
	}

	hd, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	homeDir = hd
	return homeDir
}
