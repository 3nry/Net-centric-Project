package models

import "gorm.io/gorm"

// UserProgress lưu tiến độ đọc của từng user với từng manga
// Primary key là cặp (UserID, MangaID) — theo đúng schema đề bài
type UserProgress struct {
	gorm.Model

	UserID         uint `json:"userId" gorm:"not null;index"`
	MangaID        uint `json:"mangaId" gorm:"not null;index"`
	CurrentChapter int  `json:"currentChapter" gorm:"default:0"`

	// "reading" | "completed" | "plan_to_read"
	Status string `json:"status" gorm:"default:plan_to_read"`

	// Quan hệ (optional — để GORM biết join)
	User  User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Manga Manga `json:"manga,omitempty" gorm:"foreignKey:MangaID"`
}
