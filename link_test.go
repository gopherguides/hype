package hype

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func Test_Link_MarshalJSON(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	link := &Link{
		Element: NewEl("a", nil),
	}
	link.Nodes = append(link.Nodes, Text("This is a link"))

	err := link.Set("href", "https://example.com")
	r.NoError(err)

	testJSON(t, "link", link)
}

func Test_Link_Execute_Disabled(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	link := &Link{Element: NewEl("a", nil)}
	r.NoError(link.Set("href", "https://example.com"))

	doc := &Document{
		FS:     fstest.MapFS{},
		Parser: &Parser{},
	}

	err := link.Execute(context.Background(), doc)
	r.NoError(err)
}

func Test_Link_Execute_HTTP_OK(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cfg := DefaultLinkCheckConfig()
	cfg.Enabled = true

	link := &Link{Element: NewEl("a", nil)}
	r.NoError(link.Set("href", srv.URL+"/page"))

	doc := &Document{
		FS: fstest.MapFS{},
		Parser: &Parser{
			LinkCheck:     cfg,
			LinkValidator: NewLinkValidator(cfg),
		},
	}

	err := link.Execute(context.Background(), doc)
	r.NoError(err)
}

func Test_Link_Execute_HTTP_NotFound(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	cfg := DefaultLinkCheckConfig()
	cfg.Enabled = true

	link := &Link{Element: NewEl("a", nil)}
	r.NoError(link.Set("href", srv.URL+"/missing"))

	doc := &Document{
		FS: fstest.MapFS{},
		Parser: &Parser{
			LinkCheck:     cfg,
			LinkValidator: NewLinkValidator(cfg),
		},
	}

	err := link.Execute(context.Background(), doc)
	r.Error(err)

	var lce LinkCheckError
	r.ErrorAs(err, &lce)
	r.Equal(404, lce.StatusCode)
}

func Test_Link_Execute_NonHTTP(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cfg := DefaultLinkCheckConfig()
	cfg.Enabled = true

	link := &Link{Element: NewEl("a", nil)}
	r.NoError(link.Set("href", "mailto:test@example.com"))

	doc := &Document{
		FS: fstest.MapFS{},
		Parser: &Parser{
			LinkCheck:     cfg,
			LinkValidator: NewLinkValidator(cfg),
		},
	}

	err := link.Execute(context.Background(), doc)
	r.NoError(err)
}

func Test_Link_Execute_NoHref(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cfg := DefaultLinkCheckConfig()
	cfg.Enabled = true

	link := &Link{Element: NewEl("a", nil)}

	doc := &Document{
		FS: fstest.MapFS{},
		Parser: &Parser{
			LinkCheck:     cfg,
			LinkValidator: NewLinkValidator(cfg),
		},
	}

	err := link.Execute(context.Background(), doc)
	r.NoError(err)
}
