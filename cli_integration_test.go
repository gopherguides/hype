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

	hypeCmd, stdout, _ := setupCmd(binPath, args)
	err = hypeCmd.Run()
	r.Error(err)

	r.Equal(0, stdout.Len())

	var exitErr *exec.ExitError
	r.True(errors.As(err, &exitErr), "expected exec.ExitError, got %T", err)

	exitCode := hypeCmd.ProcessState.ExitCode()
	r.Equal(1, exitCode)
}

func setupCmd(binPath string, args []string) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer) {
	cmd := exec.Command(binPath, args...)
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	return cmd, &stdout, &stderr
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

	hypeCmd, stdout, stderr := setupCmd(binPath, args)
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

func Test_Output_Flag_Error(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	binPath := filepath.Join(t.TempDir(), "hype")
	buildBin(binPath)

	root := filepath.Join(getProjectWd(), "testdata", "to_md", "source_code", "broken")
	err := os.Chdir(root)
	r.NoError(err)

	inputFile := "hype.md"
	outputFile := inputFile

	args := []string{"export", "-format=markdown", "-f", inputFile, "-o", outputFile}
	hypeCmd, _, _ := setupCmd(binPath, args)
	err = hypeCmd.Run()
	r.Error(err)

	info, err := os.Stat(filepath.Join(root, inputFile))
	r.NoError(err)

	r.True(info.Size() > 0)
}

func Test_Output_Flag_Success(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	binPath := filepath.Join(t.TempDir(), "hype")
	buildBin(binPath)

	root := filepath.Join(getProjectWd(), "testdata", "to_md", "source_code", "full")
	err := os.Chdir(root)
	r.NoError(err)

	inputFile := "hype.md"
	outputFile := filepath.Join(t.TempDir(), "output.md")

	args := []string{"export", "-format=markdown", "-f", inputFile, "-o", outputFile}
	hypeCmd, _, _ := setupCmd(binPath, args)
	err = hypeCmd.Run()
	r.NoError(err)

	info, err := os.Stat(outputFile)
	r.NoError(err)
	r.True(info.Size() > 0)

	bExp, err := os.ReadFile(filepath.Join(root, "hype.gold"))
	r.NoError(err)
	bAct, err := os.ReadFile(outputFile)
	r.NoError(err)

	exp := strings.TrimSpace(string(bExp))
	act := strings.TrimSpace(string(bAct))

	r.Equal(exp, act)
}
