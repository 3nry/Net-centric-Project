package routes

import (
	"mangahub/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {

	auth := r.Group("/auth")

	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}
}
