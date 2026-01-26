package preview

import (
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type LiveReload struct {
	clients map[*websocket.Conn]bool
	mu      sync.RWMutex
}

func NewLiveReload() *LiveReload {
	return &LiveReload{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (lr *LiveReload) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	handler := websocket.Handler(func(ws *websocket.Conn) {
		lr.mu.Lock()
		lr.clients[ws] = true
		lr.mu.Unlock()

		defer func() {
			lr.mu.Lock()
			delete(lr.clients, ws)
			lr.mu.Unlock()
			_ = ws.Close()
		}()

		for {
			var msg string
			if err := websocket.Message.Receive(ws, &msg); err != nil {
				break
			}
		}
	})
	handler.ServeHTTP(w, r)
}

func (lr *LiveReload) Reload() {
	lr.mu.RLock()
	clients := make([]*websocket.Conn, 0, len(lr.clients))
	for client := range lr.clients {
		clients = append(clients, client)
	}
	lr.mu.RUnlock()

	for _, client := range clients {
		_ = websocket.Message.Send(client, "reload")
	}
}

func (lr *LiveReload) ClientCount() int {
	lr.mu.RLock()
	defer lr.mu.RUnlock()
	return len(lr.clients)
}
