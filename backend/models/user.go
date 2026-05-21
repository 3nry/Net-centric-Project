package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username string `json:"username" gorm:"uniqueIndex;not null"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Password string `json:"-"` // json:"-" => KHÔNG bao giờ trả password về client
	Avatar   string `json:"avatar"`
	Role     string `json:"role" gorm:"default:user"`
}

// BeforeSave — GORM tự gọi trước mỗi lần Create/Save
// Hash password nếu chưa hash
func (u *User) BeforeSave(tx *gorm.DB) error {
	if len(u.Password) > 0 && len(u.Password) < 60 {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}
	return nil
}

// CheckPassword so sánh plain password với hash
func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}
