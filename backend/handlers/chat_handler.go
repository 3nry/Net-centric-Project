package handlers

import (
	"net/http"
	"time"

	"mangahub/config"

	"github.com/gin-gonic/gin"
)

func GetMessages(c *gin.Context) {

	var messages []struct {
		Username  string    `json:"username"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}

	config.DB.Raw(`
		SELECT username, content, created_at
		FROM messages
		WHERE created_at >= NOW() - INTERVAL 1 DAY
		ORDER BY created_at ASC
	`).Scan(&messages)

	c.JSON(http.StatusOK, messages)
}
