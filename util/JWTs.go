package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(userID string, ownerusername string) string {
	fmt.Println("token user id" + userID)
	var Secret_key = "posodu33333"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":           userID,
		"ownerusername": ownerusername,
		"exp":           time.Now().Add(time.Hour * 24).Unix(),
	})
	stringToken, err := token.SignedString([]byte(Secret_key))
	if err != nil {
		fmt.Println("error JWT", err)
	}
	fmt.Println("token is ", stringToken)
	return stringToken
}

func ParseJWT(token string) string {

	var Secret_key = "posodu33333"
	return Secret_key
}
