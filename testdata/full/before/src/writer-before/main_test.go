package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// snippet: test
func Test_WriteData(t *testing.T) {
	t.Parallel()

	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		t.Fatal(err)
	}

	fn := filepath.Join(dir, "hello.txt")

	f, err := os.Create(fn)

	if err != nil {
		t.Fatal(err)
	}

	data := []byte("Hello, World!")
	WriteData(f, data)

	f.Close()

	f, err = os.Open(fn)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	act := string(b)
	exp := string(data)
	if act != exp {
		t.Fatalf("expected %q, got %q", exp, act)
	}

}

// snippet: test
