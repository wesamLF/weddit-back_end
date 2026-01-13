package util

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(userID string, ownerusername string) string {
	var JWTsecret = os.Getenv("JWT_SECRET")

	if JWTsecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":           userID,
		"ownerusername": ownerusername,
		"exp":           time.Now().Add(time.Hour * 24).Unix(),
	})
	stringToken, err := token.SignedString([]byte(JWTsecret))
	if err != nil {
		fmt.Println("error JWT", err)
	}
	return stringToken
}

func ParseJWT(token string) string {
	var JWTsecret = os.Getenv("JWT_SECRET")

	if JWTsecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	return JWTsecret
}
