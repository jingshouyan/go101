package util

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"

var cl = big.NewInt(int64(len(charset)))

func GenerateRandomPassword(length int) (string, error) {

	var password []byte
	for i := 0; i < length; i++ {
		charIndex, err := rand.Int(rand.Reader, cl)
		if err != nil {
			return "", err
		}
		password = append(password, charset[charIndex.Int64()])
	}

	return string(password), nil
}
