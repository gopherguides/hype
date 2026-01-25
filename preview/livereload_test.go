package preview

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/websocket"
)

func TestNewLiveReload(t *testing.T) {
	r := require.New(t)

	lr := NewLiveReload()

	r.NotNil(lr)
	r.NotNil(lr.clients)
	r.Equal(0, lr.ClientCount())
}

func TestLiveReload_ClientCount(t *testing.T) {
	r := require.New(t)

	lr := NewLiveReload()
	r.Equal(0, lr.ClientCount())

	lr.mu.Lock()
	lr.clients[nil] = true
	lr.mu.Unlock()

	r.Equal(1, lr.ClientCount())
}

func TestLiveReload_HandleWebSocket(t *testing.T) {
	r := require.New(t)

	lr := NewLiveReload()

	server := httptest.NewServer(liveReloadHandler{lr})
	defer server.Close()

	wsURL := "ws" + server.URL[4:] + "/"

	ws, err := websocket.Dial(wsURL, "", server.URL)
	r.NoError(err)
	defer func() { _ = ws.Close() }()

	time.Sleep(50 * time.Millisecond)
	r.Equal(1, lr.ClientCount())
}

func TestLiveReload_Reload(t *testing.T) {
	r := require.New(t)

	lr := NewLiveReload()

	server := httptest.NewServer(liveReloadHandler{lr})
	defer server.Close()

	wsURL := "ws" + server.URL[4:] + "/"

	ws, err := websocket.Dial(wsURL, "", server.URL)
	r.NoError(err)
	defer func() { _ = ws.Close() }()

	time.Sleep(50 * time.Millisecond)

	msgChan := make(chan string, 1)
	go func() {
		var msg string
		_ = websocket.Message.Receive(ws, &msg)
		msgChan <- msg
	}()

	lr.Reload()

	select {
	case msg := <-msgChan:
		r.Equal("reload", msg)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for reload message")
	}
}

type liveReloadHandler struct {
	lr *LiveReload
}

func (h liveReloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.lr.HandleWebSocket(w, r)
}
