package hype

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_LinkValidator_Check_Success(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cfg := DefaultLinkCheckConfig()
	cfg.Enabled = true
	v := NewLinkValidator(cfg)

	err := v.Check(context.Background(), srv.URL+"/page")
	r.NoError(err)
}

func Test_LinkValidator_Check_HeadFallbackToGet(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cfg := DefaultLinkCheckConfig()
	cfg.Enabled = true
	v := NewLinkValidator(cfg)

	err := v.Check(context.Background(), srv.URL+"/page")
	r.NoError(err)
}

func Test_LinkValidator_Check_NotFound(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	cfg := DefaultLinkCheckConfig()
	cfg.Enabled = true
	v := NewLinkValidator(cfg)

	err := v.Check(context.Background(), srv.URL+"/missing")
	r.Error(err)

	var lce LinkCheckError
	r.ErrorAs(err, &lce)
	r.Equal(404, lce.StatusCode)
}

func Test_LinkValidator_Check_SkipNonHTTP(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cfg := DefaultLinkCheckConfig()
	cfg.Enabled = true
	v := NewLinkValidator(cfg)

	r.NoError(v.Check(context.Background(), "mailto:test@example.com"))
	r.NoError(v.Check(context.Background(), "tel:+1234567890"))
	r.NoError(v.Check(context.Background(), "ftp://example.com/file"))
	r.NoError(v.Check(context.Background(), "#section-anchor"))
}

func Test_LinkValidator_Check_SkipExcluded(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	cfg := DefaultLinkCheckConfig()
	cfg.Enabled = true
	cfg.ExcludePatterns = []string{srv.URL + "/*"}
	v := NewLinkValidator(cfg)

	err := v.Check(context.Background(), srv.URL+"/should-skip")
	r.NoError(err)
}

func Test_LinkValidator_Check_Cache(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var hitCount atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		hitCount.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cfg := DefaultLinkCheckConfig()
	cfg.Enabled = true
	v := NewLinkValidator(cfg)

	u := srv.URL + "/cached"
	r.NoError(v.Check(context.Background(), u))
	r.NoError(v.Check(context.Background(), u))

	r.Equal(int32(1), hitCount.Load())
}

func Test_LinkValidator_Check_RetryAfter(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		n := calls.Add(1)
		if n <= 2 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cfg := DefaultLinkCheckConfig()
	cfg.Enabled = true
	v := NewLinkValidator(cfg)

	err := v.Check(context.Background(), srv.URL+"/rate-limited")
	r.NoError(err)
}

func Test_LinkValidator_Check_TooManyRedirects(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, req.URL.Path, http.StatusMovedPermanently)
	}))
	defer srv.Close()

	cfg := DefaultLinkCheckConfig()
	cfg.Enabled = true
	cfg.MaxRedirects = 3
	v := NewLinkValidator(cfg)

	err := v.Check(context.Background(), srv.URL+"/loop")
	r.Error(err)
}

func Test_LinkValidator_Check_Timeout(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cfg := DefaultLinkCheckConfig()
	cfg.Enabled = true
	cfg.Timeout = 100 * time.Millisecond
	v := NewLinkValidator(cfg)

	err := v.Check(context.Background(), srv.URL+"/slow")
	r.Error(err)
}

func Test_LinkValidator_Check_FollowRedirects(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		n := calls.Add(1)
		if n == 1 {
			http.Redirect(w, req, "/final", http.StatusMovedPermanently)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cfg := DefaultLinkCheckConfig()
	cfg.Enabled = true
	v := NewLinkValidator(cfg)

	err := v.Check(context.Background(), srv.URL+"/redirect")
	r.NoError(err)
}
