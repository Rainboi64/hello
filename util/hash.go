package util

import (
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(password string, salt string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password + salt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(password string, salt string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password + salt))
	return err == nil
}