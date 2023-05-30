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
	ErrConn = fmt.Errorf("error connecting to db server")
)

// UpdateUser interacts directly with the database, updating the user information. It returns nil on success
func (db DbConn) UpdateUser(w http.ResponseWriter, user *models.User) error {
	ctx := context.Background()
	filter := bson.D{primitive.E{Key: "username", Value: user.Username}}

	if user.Firstname != "" {
		coll := db.Db.Collection("users")
		updateFilter := bson.D{{Key: "$set", Value: bson.D{{Key: "firstname", Value: user.Firstname}}}}
		update, err := coll.UpdateOne(ctx, filter, updateFilter)
		if err != nil {
			return ErrConn
		}
		json.NewEncoder(w).Encode(update)
	}
	if user.Lastname != "" {
		coll := db.Db.Collection("users")
		updateFilter := bson.D{{Key: "$set", Value: bson.D{{Key: "firstname", Value: user.Firstname}}}}
		update, err := coll.UpdateOne(ctx, filter, updateFilter)
		if err != nil {
			return ErrConn
		}
		json.NewEncoder(w).Encode(update)
	}
	if user.Password != "" {
		coll := db.Db.Collection("users")
		updateFilter := bson.D{{Key: "$set", Value: bson.D{{Key: "firstname", Value: user.Firstname}}}}
		update, err := coll.UpdateOne(ctx, filter, updateFilter)
		if err != nil {
			return ErrConn
		}
		json.NewEncoder(w).Encode(update)
	}

	// In case of email, validation code should be sent to the new email before updating in the database.
	if user.Email != "" {
		coll := db.Db.Collection("users")
		updateFilter := bson.D{{Key: "$set", Value: bson.D{{Key: "firstname", Value: user.Firstname}}}}
		update, err := coll.UpdateOne(ctx, filter, updateFilter)
		if err != nil {
			return ErrConn
		}
		json.NewEncoder(w).Encode(update)
	}
	return nil
}

// GetUser returns user password from the database. It returns a nil error if no error is returned.
func (db DbConn) GetUser(user models.User) (string, error) {
	ctx := context.Background()
	filter := bson.D{primitive.E{Key: "email", Value: user.Email}}
	coll := db.Db.Collection("users")
	find := coll.FindOne(ctx, filter)

	info := models.User{}
	err := find.Decode(&info)
	if err != nil {
		return "", ErrDb
	}
	return info.Password, nil
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

	filter := bson.D{primitive.E{Key: "username", Value: user}}
	coll := db.Db.Collection("post")
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
