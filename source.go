package hype

import (
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"testing/fstest"
	"time"
)

// Sourceable is a tag that has a src attribute.
type Sourceable interface {
	Tag
	Source() (Source, bool)
}

// SetSourceable is a tag that can have its src attribute set.
type SetSourceable interface {
	Sourceable
	SetSource(s string)
}

type Source string

// Ext returns the file extension of the source.
func (s Source) Ext() string {
	return filepath.Ext(string(s))
}

// Base returns the base name of the source.
func (s Source) Base() string {
	return filepath.Base(string(s))
}

// Dir returns the directory of the source.
func (s Source) Dir() string {
	return filepath.Dir(string(s))
}

func (s Source) String() string {
	return string(s)
}

// Lang returns the language of the source
// based on the file extension.
// Defaults to "plain".
func (s Source) Lang() string {
	ext := s.Ext()
	ext = strings.TrimPrefix(ext, ".")

	if len(ext) == 0 {
		return "plain"
	}

	return ext
}

// MimeType of the source. Defaults to "text/plain".
func (s Source) MimeType() string {
	m := mime.TypeByExtension(s.Ext())
	if len(m) == 0 {
		return "text/plain"
	}

	return m
}

// SchemeLess returns the source without the scheme.
func (s Source) Schemeless() string {
	sc, _ := s.Scheme()
	return strings.TrimPrefix(s.String(), fmt.Sprintf("%s://", sc))
}

// Scheme returns the scheme of the source.
func (s Source) Scheme() (string, error) {
	u, err := url.Parse(string(s))

	if err != nil {
		return "", err
	}

	sc := u.Scheme
	if len(sc) == 0 {
		sc = "file"
	}

	return sc, nil
}

// IsFile returns true if the source is a file.
func (s Source) IsFile() bool {
	// its ok to ignore the error
	// before the zero value of
	// source can be used for the
	// check.
	// if an error is returned
	// then the scheme will not
	// match "file"
	sc, _ := s.Scheme()
	return sc == "file"
}

// IsHTTP returns true if the source is a http(s) source.
func (s Source) IsHTTP() bool {
	// its ok to ignore the error
	// before the zero value of
	// source can be used for the
	// check.
	// if an error is returned
	// then the scheme will not
	// match "http" or "https"
	sc, _ := s.Scheme()
	return sc == "http" || sc == "https"
}

// StatFile returns the file info of the source.
func (s Source) StatFile(cab fs.FS) (fs.FileInfo, error) {
	if !s.IsFile() {
		return nil, fmt.Errorf("source is not a file: %q", s)
	}

	return fs.Stat(cab, s.Schemeless())
}

// StatHTTP returns the file info of an http(s) source.
func (s Source) StatHTTP(client *http.Client) (fs.FileInfo, error) {
	if !s.IsHTTP() {
		return nil, fmt.Errorf("source is not a http(s) source: %q", s)
	}

	res, err := client.Get(string(s))
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("http(s) source returned status code: %q (%d)", s, res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	mod, err := time.Parse(http.TimeFormat, res.Header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}

	cab := fstest.MapFS{
		s.Schemeless(): {
			Data:    b,
			Mode:    fs.ModeIrregular,
			ModTime: mod,
		},
	}

	return cab.Stat(s.Schemeless())
}

// SetTag will set the "src" attribute of the tag.
// If the tag implements the SetSourceable
// that will be used.
// Otherwise the tags attributes will be set directly.
func (s Source) SetTag(tag Tag) {
	if sc, ok := tag.(SetSourceable); ok {
		sc.SetSource(string(s))
		return
	}

	node := tag.DaNode()
	node.Set("src", string(s))
}
