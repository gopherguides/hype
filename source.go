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

type Sourceable interface {
	Tag
	Source() (Source, bool)
}

type SetSourceable interface {
	Sourceable
	SetSource(s string)
}

type Source string

func (s Source) Ext() string {
	return filepath.Ext(string(s))
}

func (s Source) Base() string {
	return filepath.Base(string(s))
}

func (s Source) Dir() string {
	return filepath.Dir(string(s))
}

func (s Source) String() string {
	return string(s)
}

func (s Source) Lang() string {
	ext := s.Ext()
	ext = strings.TrimPrefix(ext, ".")

	if len(ext) == 0 {
		return "plain"
	}

	return ext
}

func (s Source) MimeType() string {
	m := mime.TypeByExtension(s.Ext())
	if len(m) == 0 {
		return "text/plain"
	}

	return m
}

func (s Source) Schemeless() string {
	sc, _ := s.Scheme()
	return strings.TrimPrefix(s.String(), fmt.Sprintf("%s://", sc))
}

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

func (s Source) StatFile(cab fs.FS) (fs.FileInfo, error) {
	if !s.IsFile() {
		return nil, fmt.Errorf("source is not a file: %q", s)
	}

	return fs.Stat(cab, s.Schemeless())
}

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

	// Last-Modified
	// Last-Modified: <day-name>, <day> <month> <year> <hour>:<minute>:<second> GMT
	// Last-Modified: Wed, 21 Oct 2015 07:28:00 GMT
	// RFC1123

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

// SetTag will set the "src" attribute
// of the tag.
// If the tag implements the SetSourceable
// that will be used.
// Otherwise the tags attributes will be set
// directly.
func (s Source) SetTag(tag Tag) {
	if sc, ok := tag.(SetSourceable); ok {
		sc.SetSource(string(s))
		return
	}

	node := tag.DaNode()
	node.Set("src", string(s))
}
