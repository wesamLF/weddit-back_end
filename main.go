package main

import (
	"Weddit_back-end/middleware"
	"Weddit_back-end/routes"
	"net/http"
	"os"
	"log"
)

func main() {
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
