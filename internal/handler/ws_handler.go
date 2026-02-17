package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/arturo/autohost-cloud-api/internal/handler/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// En producción, deberías validar el origen apropiadamente
		return true
	},
}

type WSHandler struct {
	clients   map[string]*Client
	clientsMu sync.RWMutex
}

func NewWSHandler() *WSHandler {
	return &WSHandler{
		clients: make(map[string]*Client),
	}
}

type Client struct {
	NodeID string
	Conn   *websocket.Conn
	Send   chan []byte
	mu     sync.Mutex
}

type Message struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

func (h *WSHandler) Routes(nodeAuthMiddleware func(http.Handler) http.Handler) http.Handler {
	r := chi.NewRouter()

	r.Group(func(protected chi.Router) {
		protected.Use(nodeAuthMiddleware)
		protected.Get("/ws", h.HandleWebSocket)
	})

	return r
}

func (h *WSHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	nodeToken := middleware.GetNodeToken(r.Context())
	if nodeToken == nil || nodeToken.NodeID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Upgrade la conexión HTTP a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	client := &Client{
		NodeID: nodeToken.NodeID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	h.registerClient(client)
	defer h.unregisterClient(client)

	// Enviar mensaje de bienvenida
	welcomeMsg := Message{
		Type:      "connected",
		Timestamp: time.Now(),
	}
	if err := client.WriteJSON(welcomeMsg); err != nil {
		log.Printf("Error sending welcome message: %v", err)
		return
	}

	// Iniciar goroutines para lectura y escritura
	go client.writePump()
	client.readPump(h)
}

func (h *WSHandler) registerClient(client *Client) {
	h.clientsMu.Lock()
	defer h.clientsMu.Unlock()
	h.clients[client.NodeID] = client
	log.Printf("Node connected: %s (Total: %d)", client.NodeID, len(h.clients))
}

func (h *WSHandler) unregisterClient(client *Client) {
	h.clientsMu.Lock()
	defer h.clientsMu.Unlock()
	if _, ok := h.clients[client.NodeID]; ok {
		delete(h.clients, client.NodeID)
		close(client.Send)
		log.Printf("Node disconnected: %s (Total: %d)", client.NodeID, len(h.clients))
	}
}

func (h *WSHandler) BroadcastToNode(nodeID string, message []byte) error {
	h.clientsMu.RLock()
	client, ok := h.clients[nodeID]
	h.clientsMu.RUnlock()

	if !ok {
		return websocket.ErrCloseSent
	}

	select {
	case client.Send <- message:
		return nil
	default:
		return websocket.ErrCloseSent
	}
}

func (c *Client) readPump(handler *WSHandler) {
	defer func() {
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		msg.Timestamp = time.Now()
		log.Printf("Received message from node %s: %s", c.NodeID, msg.Type)

		// Procesar mensaje según su tipo
		switch msg.Type {
		case "ping":
			response := Message{
				Type:      "pong",
				Timestamp: time.Now(),
			}
			if err := c.WriteJSON(response); err != nil {
				log.Printf("Error sending pong: %v", err)
				return
			}
		case "echo":
			// Echo el mensaje de vuelta
			if err := c.WriteJSON(msg); err != nil {
				log.Printf("Error echoing message: %v", err)
				return
			}
		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Agregar mensajes en cola al actual
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) WriteJSON(v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Conn.WriteJSON(v)
}
