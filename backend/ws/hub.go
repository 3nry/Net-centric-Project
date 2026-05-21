package ws

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients map[*websocket.Conn]string
	Mu      sync.Mutex
}

var HubInstance = &Hub{
	Clients: make(map[*websocket.Conn]string),
}

func AddClient(conn *websocket.Conn, username string) {

	HubInstance.Mu.Lock()
	defer HubInstance.Mu.Unlock()

	HubInstance.Clients[conn] = username
}

func RemoveClient(conn *websocket.Conn) {

	HubInstance.Mu.Lock()
	defer HubInstance.Mu.Unlock()

	delete(HubInstance.Clients, conn)
}

func Broadcast(v interface{}) {

	data, _ := json.Marshal(v)

	HubInstance.Mu.Lock()
	defer HubInstance.Mu.Unlock()

	for conn := range HubInstance.Clients {
		conn.WriteMessage(websocket.TextMessage, data)
	}
}

func BroadcastRaw(data []byte) {

	HubInstance.Mu.Lock()
	defer HubInstance.Mu.Unlock()

	for conn := range HubInstance.Clients {
		conn.WriteMessage(websocket.TextMessage, data)
	}
}
