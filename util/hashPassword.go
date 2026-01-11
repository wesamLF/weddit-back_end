package util

import "golang.org/x/crypto/bcrypt"

func HashPassword(passwrod string) (string, error) {

	b, err := bcrypt.GenerateFromPassword([]byte(passwrod), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
