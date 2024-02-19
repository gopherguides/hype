package hype

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func saveErrorJSON(t testing.TB, name string, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("error is nil")
		return
	}

	fp := filepath.Join("testdata", "errors", "json")
	// os.RemoveAll(fp)
	if err := os.MkdirAll(fp, 0755); err != nil {
		t.Fatal(err)
	}
	fp = filepath.Join(fp, name+".json")

	f, ex := os.Create(fp)
	if ex != nil {
		t.Fatal(ex)
	}
	defer f.Close()

	val, ok := err.(json.Marshaler)
	if !ok {
		fmt.Fprintf(f, "%+v", err)
		return
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(val); err != nil {
		t.Fatal(err)
	}
}
