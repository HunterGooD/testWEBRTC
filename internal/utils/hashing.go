package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateHashFromString(s string) {
	sha := sha256.New()
	hex.EncodeToString(sha.Sum([]byte(s)))
}

func CreateHashFromBytes(b []byte) {
	sha := sha256.New()
	hex.EncodeToString(sha.Sum(b))
}

func RandomBytes(length int) ([]byte, error) {
	rBytes := make([]byte, length)
	_, err := rand.Read(rBytes)
	if err != nil {
		return nil, err
	}
	return rBytes, nil
}
