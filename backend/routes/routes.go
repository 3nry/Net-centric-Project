package routes

import (
	"mangahub/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// auth routes
	AuthRoutes(r)

	// manga routes
	r.GET("/manga", handlers.GetManga)
	r.GET("/manga/:id", handlers.GetMangaByID)
	r.POST("/manga", handlers.CreateManga)

	// chapter routes
	r.GET("/manga/:id/chapters", handlers.GetChapters)
	r.GET("/chapter/:id/pages", handlers.GetPages)

	// user routes
	r.GET("/users/me", handlers.GetMe)

	// websocket
	r.GET("/ws", handlers.HandleWS)
	r.GET("/messages", handlers.GetMessages)
}
