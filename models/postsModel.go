package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PostID        primitive.ObjectID `bson:"postId" json:"postId"`
	Content       string             `bson:"content" json:"content"`
	OwnerId       primitive.ObjectID `bson:"ownerId" json:"ownerId"`
	OwnerUsername string             `bson:"ownerUsername" json:"ownerUsername"`
}

type Post struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title         string             `json:"title" bson:"title"`
	Desc          string             `json:"desc" bson:"desc"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	OwnerUsername string             `bson:"ownerUsername" json:"ownerUsername"`
	Owner         primitive.ObjectID `bson:"ownerid,omitempty" json:"ownerid"`
}

type RespnsePostData struct {
	Message  string `bson:"message" json:"message"`
	Username string `bson:"username" json:"username"`
	PostId   string `bson:"ownerid,omitempty" json:"ownerid"`
}
