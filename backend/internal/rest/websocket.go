package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"cargomax-api/internal/auth"
	"cargomax-api/internal/config"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// newUpgrader creates a WebSocket upgrader that validates the Origin header
// against the configured frontend URL to prevent cross-site WebSocket hijacking.
func newUpgrader(cfg *config.Config) websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				// Allow connections with no Origin header (non-browser clients).
				return true
			}
			return origin == cfg.FrontendURL
		},
	}
}

type wsClient struct {
	conn     *websocket.Conn
	tenantID uuid.UUID
	send     chan []byte
}

type Hub struct {
	clients    map[*wsClient]bool
	broadcast  chan broadcastMsg
	register   chan *wsClient
	unregister chan *wsClient
	mu         sync.RWMutex
	config     *config.Config
	upgrader   websocket.Upgrader
}

type broadcastMsg struct {
	tenantID uuid.UUID
	msgType  string // "tracking" or "alert"
	data     []byte
}

func NewHub(cfg *config.Config) *Hub {
	return &Hub{
		clients:    make(map[*wsClient]bool),
		broadcast:  make(chan broadcastMsg, 256),
		register:   make(chan *wsClient),
		unregister: make(chan *wsClient),
		config:     cfg,
		upgrader:   newUpgrader(cfg),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				if client.tenantID == msg.tenantID {
					select {
					case client.send <- msg.data:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastTracking(tenantID uuid.UUID, data interface{}) {
	msg, err := json.Marshal(map[string]interface{}{
		"type": "tracking",
		"data": data,
	})
	if err != nil {
		return
	}
	h.broadcast <- broadcastMsg{tenantID: tenantID, msgType: "tracking", data: msg}
}

func (h *Hub) BroadcastAlert(tenantID uuid.UUID, data interface{}) {
	msg, err := json.Marshal(map[string]interface{}{
		"type": "alert",
		"data": data,
	})
	if err != nil {
		return
	}
	h.broadcast <- broadcastMsg{tenantID: tenantID, msgType: "alert", data: msg}
}

// HandleTrackingWS handles WS /ws/tracking/live
func (h *Hub) HandleTrackingWS(w http.ResponseWriter, r *http.Request) {
	tenantID := h.authenticateWS(w, r)
	if tenantID == uuid.Nil {
		return
	}
	h.serveWS(w, r, tenantID)
}

// HandleAlertsWS handles WS /ws/alerts
func (h *Hub) HandleAlertsWS(w http.ResponseWriter, r *http.Request) {
	tenantID := h.authenticateWS(w, r)
	if tenantID == uuid.Nil {
		return
	}
	h.serveWS(w, r, tenantID)
}

func (h *Hub) authenticateWS(w http.ResponseWriter, r *http.Request) uuid.UUID {
	// Get token from query parameter
	token := r.URL.Query().Get("token")
	if token == "" {
		// Try Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 {
				token = parts[1]
			}
		}
	}
	if token == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return uuid.Nil
	}

	claims, err := auth.ValidateToken(h.config.JWTPublicKey, token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return uuid.Nil
	}
	return claims.TenantID
}

func (h *Hub) serveWS(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}

	client := &wsClient{
		conn:     conn,
		tenantID: tenantID,
		send:     make(chan []byte, 256),
	}

	h.register <- client

	// Writer goroutine
	go func() {
		defer func() {
			conn.Close()
		}()
		for msg := range client.send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				break
			}
		}
	}()

	// Reader goroutine (just to detect disconnect)
	go func() {
		defer func() {
			h.unregister <- client
			conn.Close()
		}()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()
}
