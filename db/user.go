package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/A-Victory/blog-API/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrDb   = fmt.Errorf("user not found")
	ErrConn = fmt.Errorf("error communicating with db server")
)

// UpdateUser interacts directly with the database, updating the user information. It returns nil on success
func (db DbConn) UpdateUser(w http.ResponseWriter, username string, user *models.User) error {
	filter := bson.M{"username": username}

	if user.Firstname != "" {
		coll := db.Db.Collection("users")
		updateFilter := bson.D{{Key: "$set", Value: bson.D{{Key: "firstname", Value: user.Firstname}}}}
		update, err := coll.UpdateOne(ctx, filter, updateFilter)
		if err != nil {
			return ErrConn
		}
		json.NewEncoder(w).Encode(update.MatchedCount)
	}
	if user.Lastname != "" {
		coll := db.Db.Collection("users")
		updateFilter := bson.D{{Key: "$set", Value: bson.D{{Key: "lastnamw", Value: user.Lastname}}}}
		update, err := coll.UpdateOne(ctx, filter, updateFilter)
		if err != nil {
			return ErrConn
		}
		json.NewEncoder(w).Encode(update.MatchedCount)
	}

	// create a seperate function for updating password using bycrpt package.
	/*
		if user.Password != "" {
			coll := db.Db.Collection("users")
			updateFilter := bson.D{{Key: "$set", Value: bson.D{{Key: "firstname", Value: user.Firstname}}}}
			update, err := coll.UpdateOne(ctx, filter, updateFilter)
			if err != nil {
				return ErrConn
			}
			json.NewEncoder(w).Encode(update.MatchedCount)
		}

		// In case of email, validation code should be sent to the new email before updating in the database.
		if user.Email != "" {
			coll := db.Db.Collection("users")
			updateFilter := bson.D{{Key: "$set", Value: bson.D{{Key: "firstname", Value: user.Firstname}}}}
			update, err := coll.UpdateOne(ctx, filter, updateFilter)
			if err != nil {
				return ErrConn
			}
			json.NewEncoder(w).Encode(update.MatchedCount)
		}
	*/
	return nil
}

// GetUser returns user password from the database. It returns a nil error if no error is returned.
func (db DbConn) GetUser(user models.User) (username, password string, err error) {
	ctx := context.Background()
	filter := bson.D{primitive.E{Key: "email", Value: user.Email}}
	coll := db.Db.Collection("users")
	find := coll.FindOne(ctx, filter)

	info := models.User{}
	err = find.Decode(&info)
	if err != nil {
		return "", "", err
	}
	return info.Username, info.Password, nil
}

// CreateUser creates a new user document in the database.
func (db DbConn) CreateUser(user models.User) (interface{}, error) {
	coll := db.Db.Collection("users")
	insert, err := coll.InsertOne(ctx, user)
	if err != nil {
		return nil, ErrConn
	}
	return insert.InsertedID, nil
}

func (db DbConn) DeleteUser(user string) (*mongo.DeleteResult, error) {
	coll := db.Db.Collection("users")
	filter := bson.D{primitive.E{Key: "username", Value: user}}
	delete, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return nil, ErrConn
	}

	return delete, nil
}

func (db DbConn) SearchUser(user string) ([]models.Post, error) {
	posts := []models.Post{}

	filter := bson.M{"username": user}
	coll := db.Db.Collection("posts")
	int, err := coll.CountDocuments(ctx, bson.D{})
	if err != nil {
		if int <= 1 {
			return nil, fmt.Errorf("post collection not yet created")
		}
		return nil, err
	}

	cur, err := coll.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNilDocument {
			return nil, errors.New("no record found for user")
		}
		return nil, ErrConn
	}

	for cur.Next(ctx) {
		err := cur.Decode(posts)
		if err != nil {
			return nil, ErrConn
		}
	}

	return posts, nil
}
