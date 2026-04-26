package utils

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateStateToken() string {
	byteArray := make([]byte, 16)
	rand.Read(byteArray)
	return base64.URLEncoding.EncodeToString(byteArray)
}

func GenerateSessionID() string {
	id := uuid.New().String()
	return id
}

func HashPassword(password string) string {
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPassword := string(hashedBytes)
	return hashedPassword
}

func CompareHashPasswords(password_hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
