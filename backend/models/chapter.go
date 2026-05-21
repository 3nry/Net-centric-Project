package models

import "gorm.io/gorm"

type Chapter struct {
	gorm.Model

	MangaID uint   `json:"mangaId"`
	Title   string `json:"title"`
	Number  int    `json:"number"`
}
