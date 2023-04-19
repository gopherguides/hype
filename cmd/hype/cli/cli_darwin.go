package cli

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	usr, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	xp := os.Getenv("PATH")
	paths := []string{
		xp,
		"/opt/homebrew/bin",
		"/usr/local/bin",
		"/usr/bin",
		"/bin",
		"/usr/sbin",
		"/sbin",
		filepath.Join(usr, "go", "bin"),
	}

	xp = strings.Join(paths, ":")
	os.Setenv("PATH", xp)
}
