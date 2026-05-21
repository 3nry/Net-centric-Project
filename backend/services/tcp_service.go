package services

import (
	"fmt"
	"net"

	"mangahub/config"
	"mangahub/models"
	"mangahub/ws"
)

func StartTCP() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Printf("❌ TCP error: %v\n", err)
		return
	}
	defer l.Close()

	fmt.Println("📡 TCP Server running on :9000")

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go handleTCP(conn)
	}
}

func handleTCP(c net.Conn) {
	defer c.Close()

	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		return
	}

	data := string(buf[:n])

	msg := models.Message{
		Username: "ADMIN NOTE !!!",
		Content:  data,
	}

	config.DB.Create(&msg)

	ws.Broadcast(map[string]interface{}{
		"username": msg.Username,
		"text":     msg.Content,
		"time":     msg.CreatedAt.Format("15:04"),
	})
}
