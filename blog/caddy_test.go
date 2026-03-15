package blog

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func freePort(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer l.Close()
	_, port, err := net.SplitHostPort(l.Addr().String())
	require.NoError(t, err)
	return port
}

func setupPublicDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	largeHTML := "<html><body>" + string(make([]byte, 1024)) + "</body></html>"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.html"), []byte(largeHTML), 0644))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "about"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "about", "index.html"), []byte("<html><body>About</body></html>"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "style.css"), []byte("body { margin: 0; }"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "app.js"), []byte("console.log('hello');"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "404.html"), []byte("<html><body>Not Found</body></html>"), 0644))

	return dir
}

func startTestServer(t *testing.T, publicDir string) string {
	t.Helper()
	port := freePort(t)
	addr := ":" + port

	err := StartProductionServer(ProductionConfig{
		PublicDir: publicDir,
		Addr:      addr,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		StopProductionServer()
	})

	baseURL := "http://127.0.0.1:" + port
	require.Eventually(t, func() bool {
		resp, err := http.Get(baseURL + "/")
		if err != nil {
			return false
		}
		resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, 5*time.Second, 100*time.Millisecond, "server did not become ready")

	return baseURL
}

func TestProductionServerSecurityHeaders(t *testing.T) {
	publicDir := setupPublicDir(t)
	baseURL := startTestServer(t, publicDir)

	resp, err := http.Get(baseURL + "/")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	require.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
	require.Equal(t, "strict-origin-when-cross-origin", resp.Header.Get("Referrer-Policy"))
	require.Contains(t, resp.Header.Get("Permissions-Policy"), "camera=()")
}

func TestProductionServerCompression(t *testing.T) {
	publicDir := setupPublicDir(t)
	baseURL := startTestServer(t, publicDir)

	req, err := http.NewRequest("GET", baseURL+"/", nil)
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")

	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))
}

func TestProductionServerStaticAssetCaching(t *testing.T) {
	publicDir := setupPublicDir(t)
	baseURL := startTestServer(t, publicDir)

	resp, err := http.Get(baseURL + "/style.css")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Contains(t, resp.Header.Get("Cache-Control"), "max-age=31536000")
	require.Contains(t, resp.Header.Get("Cache-Control"), "immutable")
}

func TestProductionServerHTMLCaching(t *testing.T) {
	publicDir := setupPublicDir(t)
	baseURL := startTestServer(t, publicDir)

	resp, err := http.Get(baseURL + "/")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Contains(t, resp.Header.Get("Cache-Control"), "max-age=3600")
	require.Contains(t, resp.Header.Get("Cache-Control"), "must-revalidate")
}

func TestProductionServerCleanURLs(t *testing.T) {
	publicDir := setupPublicDir(t)
	baseURL := startTestServer(t, publicDir)

	resp, err := http.Get(baseURL + "/about/")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestProductionServer404(t *testing.T) {
	publicDir := setupPublicDir(t)
	baseURL := startTestServer(t, publicDir)

	resp, err := http.Get(baseURL + "/nonexistent")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestProductionServerRequiresPublicDir(t *testing.T) {
	err := StartProductionServer(ProductionConfig{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "public directory is required")
}

func TestProductionWatchConflict(t *testing.T) {
	err := fmt.Errorf("cannot use -production and -watch together")
	require.Contains(t, err.Error(), "cannot use -production and -watch together")
}
