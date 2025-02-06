package helper

import (
	"net/mail"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func MatchesHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	reUppercase := regexp.MustCompile(`[A-Z]`)
	reDigit := regexp.MustCompile(`\d`)
	reSymbol := regexp.MustCompile(`[\W_]`)
	return reUppercase.MatchString(password) && reDigit.MatchString(password) && reSymbol.MatchString(password)
}
