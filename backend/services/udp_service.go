package services

import (
	"encoding/json"
	"fmt"
	"mangahub/ws"
	"net"
	"time"
)

type Notification struct {
	Type      string `json:"type"`
	MangaID   string `json:"manga_id"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func StartUDP() {

	addr, _ := net.ResolveUDPAddr("udp", ":9001")

	conn, err := net.ListenUDP("udp", addr)

	if err != nil {
		return
	}

	defer conn.Close()

	fmt.Println("📡 UDP Server running on port 9001")

	buf := make([]byte, 1024)

	for {

		n, remoteAddr, _ := conn.ReadFromUDP(buf)

		message := string(buf[:n])

		fmt.Printf("[UDP] From %v: %s\n", remoteAddr, message)

		// TẠO NOTIFICATION
		noti := Notification{
			Type:      "chapter_release",
			MangaID:   "one-piece",
			Message:   message,
			Timestamp: time.Now().Unix(),
		}

		// JSON
		data, _ := json.Marshal(noti)

		// BROADCAST TO WS
		ws.BroadcastRaw(data)
	}
}
