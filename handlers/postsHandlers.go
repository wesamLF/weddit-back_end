package handlers

import (
	"Weddit_back-end/db"
	"Weddit_back-end/middleware"
	"Weddit_back-end/models"
	"encoding/json"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	Message string `json:"message"`
}
type PostRespnse struct {
	Message string      `json:"message"`
	Data    models.Post `json:"data"`
}
type PostsRespnse struct {
	Message string        `json:"message"`
	Data    []models.Post `json:"data"`
}
type CommentResponse struct {
	Message string         `json:"message"`
	Data    models.Comment `json:"data"`
}
type CommentsResponse struct {
	Message string           `json:"message"`
	Data    []models.Comment `json:"data"`
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {

	result, err := db.GetAllPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
	resp := PostsRespnse{
		Message: "success",
		Data:    *result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GetPostsByUsernameHandler(w http.ResponseWriter, r *http.Request) {

	username := r.PathValue("username")
	result, err := db.GetPostsByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(result)

}

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	json.NewDecoder(r.Body).Decode(&post)
	myId, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	userID, err := primitive.ObjectIDFromHex(myId)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}
	myUsername, ok := r.Context().Value(middleware.UsernameKey).(string)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	if post.Title == "" || post.Desc == "" {
		http.Error(w, "missing data", http.StatusBadRequest)
		return
	}
	newPostId, err := db.InsertOnePost(post.Title, post.Desc, myUsername, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := PostRespnse{
		Message: "post has been created successfully",
		Data:    *newPostId,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {

	var comment models.Comment
	postID := r.PathValue("id")
	json.NewDecoder(r.Body).Decode(&comment)

	if postID == "undefined" || postID == "null" {
		http.Error(w, "invalid postid", http.StatusBadRequest)
		return
	}
	postIDobj, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		http.Error(w, "invalid postid", http.StatusBadRequest)

		return
	}

	myId, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	myUsername, ok := r.Context().Value(middleware.UsernameKey).(string)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	onwerIDobj, err := primitive.ObjectIDFromHex(myId)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)

		return
	}
	if comment.Content == "" {
		http.Error(w, "missing data", http.StatusBadRequest)
		return
	}

	_, err = db.GetOnePostByID(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := db.InsertOneComment(comment.Content, postIDobj, onwerIDobj, myUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CommentResponse{
		Message: "comment has been created successfully",
		Data:    *result,
	})
}

func GetCommentsByPostIdHandler(w http.ResponseWriter, r *http.Request) {

	var comment models.Comment
	postID := r.PathValue("id")
	json.NewDecoder(r.Body).Decode(&comment)

	if postID == "undefined" || postID == "null" {
		http.Error(w, "invalid postid", http.StatusBadRequest)
		return
	}
	postIDobj, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		http.Error(w, "invalid postid", http.StatusBadRequest)

		return
	}

	result, err := db.GetCommentsByPostId(postIDobj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CommentsResponse{
		Message: "success",
		Data:    *result,
	})
}
func GetCommentsByUsernameHandler(w http.ResponseWriter, r *http.Request) {

	var comment models.Comment
	username := r.PathValue("username")
	json.NewDecoder(r.Body).Decode(&comment)

	if username == "undefined" || username == "null" {
		http.Error(w, "invalid postid", http.StatusBadRequest)
		return
	}

	result, err := db.GetCommentsByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CommentsResponse{
		Message: "success",
		Data:    *result,
	})
}
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	myIdString, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	if postID == "undefined" || postID == "null" {
		http.Error(w, "invalid postid", http.StatusBadRequest)
		return
	}
	userID, err := primitive.ObjectIDFromHex(myIdString)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	err = db.DeleteOnePost(postID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(Response{
		Message: "post has been deleted",
	})
}
func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentId := r.PathValue("id")
	myIdString, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	if commentId == "undefined" || commentId == "null" {
		http.Error(w, "invalid postid", http.StatusBadRequest)
		return
	}
	myId, err := primitive.ObjectIDFromHex(myIdString)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	err = db.DeleteComment(commentId, myId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(Response{
		Message: "comment has been deleted",
	})
}

func GetOnePostByIDHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	newPostId, err := db.GetOnePostByID(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := PostRespnse{
		Message: "success",
		Data:    *newPostId,
	}

	json.NewEncoder(w).Encode(resp)

}

func UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var myComment models.Comment
	commentId := r.PathValue("id")
	myIdString, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	myId, err := primitive.ObjectIDFromHex(myIdString)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&myComment)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if myComment.Content == "" {
		http.Error(w, "Invalid request body, some data are missing", http.StatusBadRequest)
		return
	}
	err = db.UpdateComment(commentId, myComment.Content, myId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(Response{
		Message: "comment has been updated",
	})

}

//
///
////
/////
//////
///////
////////
///////
//////
/////
////
///
//

func GetPostsPaginationHandler(w http.ResponseWriter, r *http.Request) {
	pageNumber := r.PathValue("pagenumber")
	number, err := strconv.Atoi(pageNumber)
	if err != nil {
		http.Error(w, "bad url", http.StatusBadRequest)
		return
	}
	result, pagesCount, postsCount, err := db.GetPostsPagination(number)
	if err != nil {
		http.Error(w, "did not found the posts", http.StatusBadRequest)

	}
	response := map[string]any{
		"page":       pageNumber,
		"postsCount": postsCount,
		"pagesCount": pagesCount,
		"posts":      result,
	}
	json.NewEncoder(w).Encode(response)
}

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	var myPost models.Post
	postID := r.PathValue("id")
	myIdString, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	userID, err := primitive.ObjectIDFromHex(myIdString)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&myPost)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if myPost.Title == "" || myPost.Desc == "" {
		http.Error(w, "Invalid request body, some data are missing", http.StatusBadRequest)
		return
	}
	err = db.UpdateOnePost(postID, myPost.Title, myPost.Desc, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(Response{
		Message: "post has been updated",
	})

}
