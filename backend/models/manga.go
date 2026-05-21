package models

import "gorm.io/gorm"

type Manga struct {
	gorm.Model

	Title       string `json:"title"`
	Description string `json:"description"`
	CoverImage  string `json:"coverImage"`
	Author      string `json:"author"`
	Genre       string `json:"genre"`
	Status      string `json:"status"`
	MangaDexID  string `json:"mangaDexId" gorm:"uniqueIndex"`
}
