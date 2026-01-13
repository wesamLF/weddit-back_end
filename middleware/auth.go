package middleware

import (
	"Weddit_back-end/util"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKeyUserId string
type ContextKeyUserName string

const UserIDKey ContextKeyUserId = "userID"
const UsernameKey ContextKeyUserName = "username"

func ValidateToken(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.EnableCors(w, r)

		var JWTsecret = os.Getenv("JWT_SECRET")

		if JWTsecret == "" {
			log.Fatal("JWT_SECRET is not set")
		}
		cookie, err := r.Cookie("token")
		if err != nil {
			fmt.Println("in Missing cookie")

			http.Error(w, "Missing cookie", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value

		// Parse token with MapClaims
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Println("in unexpected signing")

				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(JWTsecret), nil
		})

		if err != nil || !token.Valid {
			fmt.Println("in Invalid token")

			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract user ID (username) from claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				fmt.Println("exp token")
				return
			}
			userID := claims["sub"]
			userName := claims["ownerusername"]

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UsernameKey, userName)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid claims", http.StatusUnauthorized)
			return
		}
	})
}
