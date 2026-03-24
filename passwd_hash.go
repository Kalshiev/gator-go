package main

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(passwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), 14)
	return string(bytes), err
}

func VerifyPassword(passwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))
	return err == nil
}
