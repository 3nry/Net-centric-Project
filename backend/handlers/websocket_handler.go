package handlers

import (
	"encoding/json"
	"fmt"
	"mangahub/ws"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 👉 inject DB từ main.go
var DB *gorm.DB

type ChatMessage struct {
	Username string `json:"username"`
	Text     string `json:"text"`
	Time     string `json:"time"`
}

func HandleWS(c *gin.Context) {

	username := c.Query("username")
	if username == "" {
		username = "Anonymous"
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WS upgrade error:", err)
		return
	}
	defer conn.Close()

	// register client
	ws.HubInstance.Mu.Lock()
	ws.HubInstance.Clients[conn] = username
	ws.HubInstance.Mu.Unlock()

	fmt.Printf("[WS] %s joined\n", username)

	// system message
	ws.Broadcast(ChatMessage{
		Username: "System",
		Text:     fmt.Sprintf("%s joined the chat 👋", username),
		Time:     time.Now().Format("15:04"),
	})

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			break
		}

		fmt.Println("[WS RAW]", string(raw)) // 🔥 DEBUG

		var incoming struct {
			Text string `json:"text"`
		}

		if err := json.Unmarshal(raw, &incoming); err != nil {
			fmt.Println("JSON error:", err)
			continue
		}

		now := time.Now()

		// ======================
		// 💾 SAVE TO DATABASE
		// ======================
		if DB != nil {
			err := DB.Exec(`
				INSERT INTO messages (username, content, created_at)
				VALUES (?, ?, ?)
			`, username, incoming.Text, now).Error

			if err != nil {
				fmt.Println("DB insert error:", err)
			}
		} else {
			fmt.Println("DB is NIL ❌")
		}

		// ======================
		// 📡 BROADCAST
		// ======================
		ws.Broadcast(ChatMessage{
			Username: username,
			Text:     incoming.Text,
			Time:     now.Format("15:04"),
		})
	}

	// remove client
	ws.HubInstance.Mu.Lock()
	delete(ws.HubInstance.Clients, conn)
	ws.HubInstance.Mu.Unlock()

	fmt.Printf("[WS] %s left\n", username)
}
