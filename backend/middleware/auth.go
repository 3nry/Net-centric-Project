package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		userID, ok := GetUserIDFromToken(c)

		if !ok {

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})

			c.Abort()

			return
		}

		// lưu user id vào context
		c.Set("user_id", userID)

		c.Next()
	}
}
