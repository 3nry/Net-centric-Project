package handlers

import (
	"net/http"

	"mangahub/config"
	"mangahub/middleware"
	"mangahub/models"

	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {

	userID, ok := middleware.GetUserIDFromToken(c)

	if !ok {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})

		return
	}

	var user models.User

	if err := config.DB.First(&user, userID).Error; err != nil {

		c.JSON(404, gin.H{
			"error": "user not found",
		})

		return
	}

	c.JSON(200, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"avatar":   user.Avatar,
		"role":     user.Role,
	})
}
