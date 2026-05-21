package main

import (
	"fmt"

	"mangahub/config"
	"mangahub/handlers"
	"mangahub/models"
	"mangahub/routes"
	"mangahub/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	config.ConnectDB()
	handlers.DB = config.DB

	config.DB.AutoMigrate(
		&models.User{},
		&models.Manga{},
		&models.Chapter{},
		&models.UserProgress{},
	)

	fmt.Println("✅ Database migrated")

	// ===== SEED MANGA =====
	var count int64

	config.DB.Model(&models.Manga{}).Count(&count)

	if count == 0 {
		services.SeedMangaFromMangaDex()
	}

	// ===== SERVICES =====
	go services.StartGRPC()
	go services.StartTCP()
	go services.StartUDP()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	routes.SetupRoutes(r)

	fmt.Println("🚀 Server running at :8080")

	r.Run(":8080")
}
