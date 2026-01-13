package handlers

import (
	"Weddit_back-end/db"
	"Weddit_back-end/models"
	"Weddit_back-end/util"
	"encoding/json"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	if user.Username == "" || user.Password == "" {
		http.Error(w, "please enter valid username and password", http.StatusUnauthorized)
		return
	}
	user, err := db.LoginToDB(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	jwtToken := util.CreateJWT(user.ID.Hex(), user.Username)
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    jwtToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	resp := models.RespnseUserData{
		Message:  "logged in successfully",
		Username: user.Username,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	if user.Username == "" || user.Password == "" {
		http.Error(w, "please enter valid username and password", http.StatusUnauthorized)
		return
	}
	var err error
	user.Password, err = util.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "try another password", http.StatusUnauthorized)
	}
	err = db.SinginToDB(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	resp := models.RespnseUserData{
		Message:  "user has been created successfully",
		Username: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// remove cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, // use true in production HTTPS
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "logged out"}`))
}
