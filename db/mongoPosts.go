package db

import (
	"Weddit_back-end/models"
	"context"
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllPosts() (*[]models.Post, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var posts []models.Post
	filter := bson.D{}
	cur, err := PostsCollection.Find(ctx, filter)
	if err != nil {
		return &posts, err
	}
	defer cur.Close(ctx)
	if err := cur.All(ctx, &posts); err != nil {
		return &posts, err
	}

	return &posts, nil

}
func GetPostsByUsername(username string) (*[]models.Post, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var posts []models.Post

	filter := bson.M{"ownerUsername": username}
	cur, err := PostsCollection.Find(ctx, filter)
	if err != nil {
		return &posts, err
	}
	defer cur.Close(ctx)
	if err := cur.All(ctx, &posts); err != nil {
		return &posts, err
	}
	if len(posts) == 0 {
		return &[]models.Post{}, nil
	}
	return &posts, nil
}

func InsertOnePost(title string, desc string, ownerUsername string, ownerid primitive.ObjectID) (*models.Post, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	myPost := models.Post{
		ID:            primitive.NewObjectID(),
		Title:         title,
		Desc:          desc,
		CreatedAt:     time.Now(),
		OwnerUsername: ownerUsername,
		Owner:         ownerid,
	}

	_, err := PostsCollection.InsertOne(ctx, myPost)
	if err != nil {
		return nil, fmt.Errorf("could not store the post")
	}
	fmt.Println("resualt :", myPost)

	return &myPost, nil
}

func InsertOneComment(content string, postID primitive.ObjectID, ownerId primitive.ObjectID, myUsername string) (*models.Comment, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	myComment := models.Comment{
		ID:            primitive.NewObjectID(),
		PostID:        postID,
		Content:       content,
		OwnerId:       ownerId,
		OwnerUsername: myUsername,
	}

	result, err := CommentsCollection.InsertOne(ctx, myComment)
	if err != nil {
		return nil, fmt.Errorf("could not create a comment")
	}
	fmt.Println("resualt :", result)
	return &myComment, nil
}

func GetCommentsByPostId(postID primitive.ObjectID) (*[]models.Comment, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"postId": postID}

	cursor, err := CommentsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	comments := []models.Comment{}
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	if comments == nil {
		return &[]models.Comment{}, nil
	}

	return &comments, nil
}
func GetCommentsByUsername(ownerUsername string) (*[]models.Comment, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"ownerUsername": ownerUsername}

	cursor, err := CommentsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	comments := []models.Comment{}
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	if comments == nil {
		return &[]models.Comment{}, nil
	}

	return &comments, nil
}

func GetOnePostByID(postID string) (*models.Post, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var result models.Post
	id, _ := primitive.ObjectIDFromHex(postID)
	filter := bson.M{"_id": id}
	err := PostsCollection.FindOne(ctx, filter).Decode(&result)

	if err == mongo.ErrNoDocuments || err != nil {
		return nil, err

	}

	return &result, nil

}
func GetPostsPagination(pageNumber int) (*[]models.Post, int, int, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var posts []models.Post
	limit := 10
	skip := (pageNumber - 1) * limit
	filter := bson.D{}
	ops := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetSkip(int64(skip)).SetLimit(int64(limit))

	postsCount, err := PostsCollection.CountDocuments(ctx, bson.D{})
	if err != nil || postsCount == 0 {
		return &posts, 0, 0, err

	}
	totalPages := int(math.Ceil(float64(postsCount) / float64(limit)))

	cursor, err := PostsCollection.Find(ctx, filter, ops)
	if err != nil {
		return &posts, 0, 0, err

	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &posts)
	if err != nil {
		return &posts, 0, 0, err
	}
	return &posts, totalPages, int(postsCount), nil

}

func UpdateOnePost(postID string, newTitle string, newDesc string, myId primitive.ObjectID) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var myPost models.Post
	id, _ := primitive.ObjectIDFromHex(postID)
	filter := bson.M{"_id": id}

	err := PostsCollection.FindOne(ctx, filter).Decode(&myPost)
	if err != nil || err == mongo.ErrNoDocuments {
		return err
	}
	if myPost.Owner != myId {
		return err
	}
	update := bson.M{"$set": bson.M{"title": newTitle, "desc": newDesc}}
	result, err := PostsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("No document found to update")
	} else if result.ModifiedCount == 0 {
		return fmt.Errorf("Document found but nothing was changed")
	}
	return nil
}

func DeleteOnePost(postID string, myId primitive.ObjectID) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var myPost models.Post
	id, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err

	}
	filter := bson.M{"_id": id}
	err = PostsCollection.FindOne(ctx, filter).Decode(&myPost)
	if err != nil {
		return err

	}
	if myPost.Owner != myId {
		return fmt.Errorf("error, you can only delete your own posts!")
	}

	_, err = CommentsCollection.DeleteMany(ctx, bson.M{
		"postId": id,
	})
	if err != nil {
		return err
	}

	result, err := PostsCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no document found")
	}
	return nil
}

func DeleteComment(commentId string, myId primitive.ObjectID) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var myComment models.Comment
	id, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		return err

	}
	filter := bson.M{"_id": id}
	err = CommentsCollection.FindOne(ctx, filter).Decode(&myComment)
	if err != nil {
		return err

	}
	if myComment.OwnerId != myId {
		return fmt.Errorf("error, you can only delete your own comments!")
	}
	result, err := CommentsCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no document found")
	}
	return nil
}

func UpdateComment(commentId string, newContent string, myId primitive.ObjectID) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var myComment models.Comment
	id, _ := primitive.ObjectIDFromHex(commentId)
	filter := bson.M{"_id": id}

	err := CommentsCollection.FindOne(ctx, filter).Decode(&myComment)
	if err != nil || err == mongo.ErrNoDocuments {
		return err
	}
	if myComment.OwnerId != myId {
		return err
	}
	update := bson.M{"$set": bson.M{"content": newContent}}
	result, err := CommentsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("No document found to update")
	} else if result.ModifiedCount == 0 {
		return fmt.Errorf("Document found but nothing was changed")
	}
	return nil
}
