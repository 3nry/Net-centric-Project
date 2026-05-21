package services

import (
	"errors"
	"os"
	"time"

	"mangahub/config"
	"mangahub/models"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user models.User) (string, error) {

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		secret = "mangahub-secret-key"
	}

	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(secret))
}

func RegisterUser(
	username string,
	email string,
	password string,
) (models.User, string, error) {

	var existing models.User

	if err := config.DB.
		Where("username = ? OR email = ?", username, email).
		First(&existing).Error; err == nil {

		return models.User{}, "", errors.New(
			"username or email already exists",
		)
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	if err := config.DB.Create(&user).Error; err != nil {

		return models.User{}, "", err
	}

	token, err := GenerateToken(user)

	if err != nil {
		return models.User{}, "", err
	}

	return user, token, nil
}

func LoginUser(
	username string,
	password string,
) (models.User, string, error) {

	var user models.User

	if err := config.DB.
		Where("username = ?", username).
		First(&user).Error; err != nil {

		return models.User{}, "", errors.New(
			"invalid username or password",
		)
	}

	if !user.CheckPassword(password) {

		return models.User{}, "", errors.New(
			"invalid username or password",
		)
	}

	token, err := GenerateToken(user)

	if err != nil {
		return models.User{}, "", err
	}

	return user, token, nil
}
