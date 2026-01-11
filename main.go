package main

import (
	"Weddit_back-end/middleware"
	"Weddit_back-end/routes"
	"net/http"
)

func main() {

	mux := routes.SetUpRoutes()
	http.ListenAndServe(":8080", middleware.CorsMiddleware(mux))

}
