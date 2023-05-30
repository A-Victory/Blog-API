package db

import (
	"context"

	"github.com/A-Victory/blog-API/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

// CreatePost creates a new post document in the database.
func (db DbConn) CreatePost(p models.Post) (*mongo.InsertOneResult, error) {
	coll := db.Db.Collection("posts")
	insert, err := coll.InsertOne(ctx, p)
	if err != nil {
		return nil, ErrConn
	}
	return insert, nil
}

// DeletePost deletes a post document from the database
func (db DbConn) DeletePost(id int) (*mongo.DeleteResult, error) {
	coll := db.Db.Collection("posts")
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	delete, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return nil, ErrConn
	}

	return delete, nil

}

// UpdatePost updates an existing post
func (db DbConn) UpdatePost(id int, p models.Post) (*mongo.UpdateResult, error) {
	coll := db.Db.Collection("posts")
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: p}}
	UpdateResult, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, ErrConn
	}
	return UpdateResult, nil
}

// Comment creates(adds) a comment field to a post document in the database
func (db DbConn) Comment(id int, com models.Comment) (*mongo.UpdateResult, error) {
	coll := db.Db.Collection("posts")
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "comment", Value: bson.D{{Key: "$each", Value: com}}}}}}
	upload, err := coll.UpdateByID(ctx, id, update)
	if err != nil {
		return nil, ErrConn
	}

	return upload, nil
}

// GetPost returns a post document from the database.
func (db DbConn) GetPost(id int) *mongo.SingleResult {
	coll := db.Db.Collection("posts")
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	post := coll.FindOne(ctx, filter)

	return post

}

// GetPosts returns all post in the database.
func (db DbConn) GetPosts() (*mongo.Cursor, error) {
	coll := db.Db.Collection("posts")
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	find, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, ErrConn
	}

	return find, nil
}

// DeletePosts deletes all posts associated to a user.
func (db DbConn) DeletePosts(user string) (*mongo.DeleteResult, error) {
	coll := db.Db.Collection("posts")
	filter := bson.D{primitive.E{Key: "username", Value: user}}

	delete, err := coll.DeleteMany(ctx, filter)
	if err != nil {
		return nil, ErrConn
	}

	return delete, nil
}

// GetUserPosts returns all posts associated with a user		“
func (db DbConn) GetUserPosts(user string) (*mongo.Cursor, error) {
	coll := db.Db.Collection("posts")
	filter := bson.D{primitive.E{Key: "username", Value: user}}

	find, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, ErrConn
	}

	return find, nil
}

// DeleteComment deletes a comment from a post document.
func (db *DbConn) DeleteComment(id int, user string) (*mongo.UpdateResult, error) {
	coll := db.Db.Collection("posts")
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"comment": bson.M{"username": user}}}

	delete, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, ErrConn
	}

	return delete, nil
}
