package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/A-Victory/blog-API/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
		updateFilter := bson.D{{Key: "$set", Value: bson.D{{Key: "lastname", Value: user.Lastname}}}}
		update, err := coll.UpdateOne(ctx, filter, updateFilter)
		if err != nil {
			return ErrConn
		}
		json.NewEncoder(w).Encode(update.MatchedCount)
	}

	return nil
}

// GetUser returns user password from the database. It returns a nil error if no error is returned.
func (db DbConn) GetUser(user models.User) (username, password string, err error) {
	ctx := context.Background()
	filter := bson.M{"email": user.Email}
	coll := db.Db.Collection("users")
	find := coll.FindOne(ctx, filter)

	info := models.User{}
	err = find.Decode(&info)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", "", err
		}
		return "", "", fmt.Errorf("an error occurred, please try again")
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

func (db DbConn) SearchUser(user string) (*mongo.Cursor, error) {

	filter := bson.M{"username": user}
	userData := db.Db.Collection("users").FindOne(ctx, filter)
	if userData.Err() == mongo.ErrNoDocuments {
		return nil, errors.New("user does not exist, enter a valid username")
	}

	coll := db.Db.Collection("posts")

	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, errors.New("no record found for user")
	}

	return cur, nil
}

func (db DbConn) ChangePassword(user string, password models.Password) (*mongo.UpdateResult, error) {

	u := models.User{}
	p := password
	coll := db.Db.Collection("users")
	filter := bson.D{primitive.E{Key: "username", Value: user}}
	if err := coll.FindOne(ctx, filter).Decode(&u); err != nil {
		return nil, errors.New("unable to decode into user: " + err.Error())
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p.Old))
	if err != nil {
		return nil, fmt.Errorf("password mismatch: %v", err)
	}

	newHashed, err := bcrypt.GenerateFromPassword([]byte(p.New), 8)
	if err != nil {
		log.Print("Error generating password: ", err)
		return nil, err
	}

	hashpswd := string(newHashed)

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: hashpswd}}}}
	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}
