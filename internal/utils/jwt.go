package utils

import (
	"time"

	"github.com/gin-gonic/gin" // Import ini WAJIB ada
	"github.com/golang-jwt/jwt/v5"

)

var SecretKey = []byte("mqfm_secret_key_123")

func GenerateToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
}

func GetUserID(c *gin.Context) uint {
	id, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	// JWT claims dari JSON seringkali dibaca sebagai float64
	if floatVal, ok := id.(float64); ok {
		return uint(floatVal)
	}
	// Jaga-jaga jika sudah uint
	if uintVal, ok := id.(uint); ok {
		return uintVal
	}
	return 0
}