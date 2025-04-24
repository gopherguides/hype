package hype_test

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildBin(binPath string) {
	buildCmd := exec.Command("go", "build", "-o", binPath, "./cmd/hype")
	if err := buildCmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func getProjectWd() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}

func Test_Cli_Error(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	binPath := filepath.Join(t.TempDir(), "hype")
	buildBin(binPath)

	root := filepath.Join(getProjectWd(), "testdata", "to_md", "source_code", "broken")
	err := os.Chdir(root)
	r.NoError(err)

	args := []string{"export", "-format=markdown", "-f", "hype.md"}

	t.Logf("Running command: %s %s, root: %s", binPath, strings.Join(args, " "), root)

	hypeCmd := exec.Command(binPath, args...)
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	hypeCmd.Stderr = &stderr
	hypeCmd.Stdout = &stdout

	err = hypeCmd.Run()
	r.Error(err)

	r.Equal(0, stdout.Len())

	var exitErr *exec.ExitError
	r.True(errors.As(err, &exitErr), "expected exec.ExitError, got %T", err)

	exitCode := hypeCmd.ProcessState.ExitCode()
	r.Equal(1, exitCode)
}

func Test_Cli_Success(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	binPath := filepath.Join(t.TempDir(), "hype")
	buildBin(binPath)

	root := filepath.Join(getProjectWd(), "testdata", "to_md", "source_code", "full")
	err := os.Chdir(root)
	r.NoError(err)

	args := []string{"export", "-format=markdown", "-f", "hype.md"}
	t.Logf("Running command: %s %s", binPath, strings.Join(args, " "))
	hypeCmd := exec.Command(binPath, args...)
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	hypeCmd.Stderr = &stderr
	hypeCmd.Stdout = &stdout

	err = hypeCmd.Run()
	if err != nil {
		fmt.Printf("%s\n", stderr.String())
	}
	r.NoError(err)
	r.Equal(0, stderr.Len())
	r.True(stdout.Len() > 0)

	exitCode := hypeCmd.ProcessState.ExitCode()
	r.Equal(0, exitCode)

	b, err := os.ReadFile(filepath.Join(root, "hype.gold"))
	r.NoError(err)
	exp := strings.TrimSpace(string(b))
	act := strings.TrimSpace(stdout.String())

	r.Equal(exp, act)
}
