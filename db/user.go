package db

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/A-Victory/blog-API/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrDb   = fmt.Errorf("user not found")
	ErrConn = fmt.Errorf("error connecting to db server")
)

// UpdateUser interacts directly with the database, updating the user information. It returns nil on success
func UpdateUser(w http.ResponseWriter, user *models.User, ac *mongo.Database) error {
	ctx := context.Background()
	filter := bson.D{{Key: "user", Value: user.Username}}

	if user.Firstname != "" {
		coll := ac.Collection("users")
		update, err := coll.UpdateOne(ctx, filter, user.Firstname)
		if err != nil {
			return ErrConn
		}
		json.NewEncoder(w).Encode(update)
	}
	if user.Lastname != "" {
		coll := ac.Collection("users")
		update, err := coll.UpdateOne(ctx, filter, user.Lastname)
		if err != nil {
			return ErrConn
		}
		json.NewEncoder(w).Encode(update)
	}
	if user.Password != "" {
		coll := ac.Collection("users")
		update, err := coll.UpdateOne(ctx, filter, user.Password)
		if err != nil {
			return ErrConn
		}
		json.NewEncoder(w).Encode(update)
	}

	// In case of email, validation code should be sent to the new email before updating in the database.
	if user.Email != "" {
		coll := ac.Collection("users")
		update, err := coll.UpdateOne(ctx, filter, user.Firstname)
		if err != nil {
			return ErrConn
		}
		json.NewEncoder(w).Encode(update)
	}
	return nil
}

// GetUser returns user password from the database. It returns a nil error if no error is returned.
func GetUser(user models.User, ac *mongo.Database) (string, error) {
	ctx := context.Background()
	filter := bson.D{{Key: "email", Value: user.Email}}
	coll := ac.Collection("users")
	find := coll.FindOne(ctx, filter)

	info := models.User{}
	err := find.Decode(&info)
	if err != nil {
		return "", ErrDb
	}
	return info.Password, nil
}

// CreateUser creates a new user document in the database.
func CreateUser(user models.User, ac *mongo.Database) (interface{}, error) {
	coll := ac.Collection("users")
	insert, err := coll.InsertOne(context.Background(), user)
	if err != nil {
		return nil, ErrConn
	}
	return insert.InsertedID, nil
}
