package routes

import (
	handlers "Weddit_back-end/handlers"
	"Weddit_back-end/middleware"
	"net/http"
)

func SetUpRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /post/{id}", handlers.GetOnePostByIDHandler)
	mux.HandleFunc("GET /posts", handlers.GetPostsHandler)
	mux.HandleFunc("GET /user/{username}/posts", handlers.GetPostsByUsernameHandler)
	mux.Handle("POST /createpost", middleware.ValidateToken(http.HandlerFunc(handlers.CreatePostsHandler)))
	mux.Handle("DELETE /delete/{id}", middleware.ValidateToken(http.HandlerFunc(handlers.DeletePostHandler)))
	mux.Handle("PUT /update/{id}", middleware.ValidateToken(http.HandlerFunc(handlers.UpdatePostHandler)))
	//////////////////
	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("POST /logout", handlers.LogoutHandler)
	mux.HandleFunc("POST /singin", handlers.SigninHandler)
	///////////////////////
	mux.Handle("POST /comment/{id}", middleware.ValidateToken(http.HandlerFunc(handlers.CreateCommentHandler)))
	mux.HandleFunc("GET /comments/{id}", handlers.GetCommentsByPostIdHandler)
	mux.HandleFunc("GET /comments/user/{username}", handlers.GetCommentsByUsernameHandler)
	mux.Handle("DELETE /comment/{id}", middleware.ValidateToken(http.HandlerFunc(handlers.DeleteCommentHandler)))
	mux.Handle("PUT /comment/{id}", middleware.ValidateToken(http.HandlerFunc(handlers.UpdateCommentHandler)))
	mux.HandleFunc("GET /posts/{pagenumber}", handlers.GetPostsPaginationHandler)

	return mux
}
