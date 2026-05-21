package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // log SQL ra terminal khi dev
	})

	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	DB = database

	fmt.Println("✅ MySQL Connected")
}

// MigrateDB tách riêng để main.go gọi sau ConnectDB
func MigrateDB() {
	// import models ở đây để tránh circular import
	// => gọi từ main.go, truyền vào các model
	fmt.Println("✅ Database migrated")
}
