package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytepassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err.Error())
	}
	return string(bytepassword)
}

func VerifyPassword(candidatePassword string, password string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(candidatePassword))
	check := true
	msg := ""
	if err != nil {
		msg = "Incorrect Password"
		check = false
	}
	return check, msg
}
