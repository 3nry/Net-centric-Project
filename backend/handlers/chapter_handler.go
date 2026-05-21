package handlers

import (
	"mangahub/config"
	"mangahub/models"
	"mangahub/services"

	"github.com/gin-gonic/gin"
)

func GetChapters(c *gin.Context) {

	id := c.Param("id")

	var manga models.Manga

	if err := config.DB.First(&manga, id).Error; err != nil {

		c.JSON(404, gin.H{
			"error": "manga not found",
		})

		return
	}

	chapters, err := services.FetchChapters(manga.MangaDexID)

	if err != nil {

		c.JSON(500, gin.H{
			"error": "failed to fetch chapters",
		})

		return
	}

	c.JSON(200, chapters)
}

func GetPages(c *gin.Context) {

	chapterID := c.Param("id")

	pages, err := services.FetchPages(chapterID)

	if err != nil {

		c.JSON(500, gin.H{
			"error": "failed to fetch pages",
		})

		return
	}

	c.JSON(200, gin.H{
		"pages": pages,
	})
}
