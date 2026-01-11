package db

import (
	"Weddit_back-end/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func LoginToDB(myUser models.User) (models.User, error) {
	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"username": myUser.Username}
	err := UsersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil || err == mongo.ErrNoDocuments {
		return user, fmt.Errorf("wrong password or username")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(myUser.Password))
	if err != nil {
		return user, fmt.Errorf("wrong password or username")
	} else {
		return user, nil
	}

}

func SinginToDB(myUser models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"username": myUser.Username}
	err := UsersCollection.FindOne(ctx, filter).Decode(&struct{}{})
	if err == nil {
		return fmt.Errorf("this username is already used")
	}

	if err == mongo.ErrNoDocuments {
		_, err = UsersCollection.InsertOne(ctx, myUser)
		if err != nil {
			return fmt.Errorf("errort")
		}
		fmt.Println("account created")

		return nil
	}

	return fmt.Errorf("somthing went wrong")

}
