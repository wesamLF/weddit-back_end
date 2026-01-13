package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURL = ""
var DBname = "Posts-DB"
var Database *mongo.Database
var UsersCollection *mongo.Collection
var PostsCollection *mongo.Collection
var CommentsCollection *mongo.Collection

var PostsCOLname = "My-Posts"
var usersCOLname = "users"
var CommentsCOLname = "comments"

func init() {

	mongoURL := os.Getenv("MONGODB_URL")
	clientOpts := options.Client().ApplyURI(mongoURL)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	Client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		fmt.Println("can not conntect to mongo")
		log.Fatal(err)
	}
	Database = Client.Database(DBname)
	UsersCollection = Database.Collection(usersCOLname)
	PostsCollection = Database.Collection(PostsCOLname)
	CommentsCollection = Database.Collection(CommentsCOLname)

}
