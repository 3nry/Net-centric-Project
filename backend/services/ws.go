package services

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	Username string `json:"username"`
	Text     string `json:"text"`
	Time     string `json:"time"`
}

type Hub struct {
	Clients map[*websocket.Conn]string
	Mu      sync.RWMutex
}

var WS = &Hub{
	Clients: make(map[*websocket.Conn]string),
}

func (h *Hub) Add(conn *websocket.Conn, username string) {
	h.Mu.Lock()
	h.Clients[conn] = username
	h.Mu.Unlock()

	fmt.Println("[WS] client added:", username)
}

func (h *Hub) Remove(conn *websocket.Conn) {
	h.Mu.Lock()
	delete(h.Clients, conn)
	h.Mu.Unlock()

	fmt.Println("[WS] client removed")
}

func (h *Hub) Broadcast(msg ChatMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("❌ JSON marshal error:", err)
		return
	}

	h.Mu.RLock()
	defer h.Mu.RUnlock()

	for conn := range h.Clients {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			fmt.Println("❌ WS write error, removing client")

			// safe remove outside loop lock
			go func(c *websocket.Conn) {
				h.Remove(c)
				c.Close()
			}(conn)
		}
	}
}
