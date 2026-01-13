package main

import (
	"Weddit_back-end/middleware"
	"Weddit_back-end/routes"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if os.Getenv("RAILWAY_ENV") == "" {
		_ = godotenv.Load()
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local testing
	}

	mux := routes.SetUpRoutes()
	log.Printf("Server listening on port %s", port)
	err := http.ListenAndServe(":"+port, middleware.CorsMiddleware(mux))
	if err != nil {
		log.Fatal(err)
	}

}
