package handlers

import (
	"net/http"

	"mangahub/config"
	"mangahub/models"

	"github.com/gin-gonic/gin"
)

func GetManga(c *gin.Context) {

	var mangas []models.Manga

	config.DB.Find(&mangas)

	c.JSON(http.StatusOK, mangas)
}

func GetMangaByID(c *gin.Context) {

	id := c.Param("id")

	var manga models.Manga

	if err := config.DB.First(&manga, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "manga not found",
		})

		return
	}

	c.JSON(http.StatusOK, manga)
}

func CreateManga(c *gin.Context) {

	var manga models.Manga

	if err := c.BindJSON(&manga); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})

		return
	}

	config.DB.Create(&manga)

	c.JSON(http.StatusOK, manga)
}
