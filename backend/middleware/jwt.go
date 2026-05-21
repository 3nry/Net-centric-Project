package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserIDFromToken(c *gin.Context) (uint, bool) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, false
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		secret = "mangahub-secret-key"
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return 0, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		return 0, false
	}

	return uint(idFloat), true
}
