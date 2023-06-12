package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/A-Victory/blog-API/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

// CreatePost creates a new post document in the database.
func (db DbConn) CreatePost(u string, p *models.Post) (*mongo.InsertOneResult, error) {
	user := models.User{}

	// first return user information from user collection
	coll := db.Db.Collection("users")
	filter := bson.M{"username": u}
	res := coll.FindOne(ctx, filter)
	if err := res.Decode(&user); err != nil {
		return nil, fmt.Errorf("could not retrieve user from user collection: %v", err)
	}

	p.Username = u
	p.Id = primitive.NewObjectID()
	p.Created_At = time.Now()
	// set user_id field to the user's _id
	p.User_id = user.Id.Hex()
	p.Reactions.UpVote = make([]models.UpVote, 0)
	p.Reactions.DownVote = make([]models.DownVote, 0)
	p.Comments = make([]models.Comment, 0)

	// insert the post into the post collection
	coll2 := db.Db.Collection("posts")
	insert, err := coll2.InsertOne(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("unable to insert post: %v", err)
	}

	// insert the post id back to the post_id flied in the user document.
	inserted := insert.InsertedID.(primitive.ObjectID)
	coll3 := db.Db.Collection("users")
	filter2 := bson.D{primitive.E{Key: "username", Value: u}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "post_id", Value: inserted.Hex()}}}}
	_, err = coll3.UpdateOne(ctx, filter2, update)
	if err != nil {
		return nil, fmt.Errorf("could not upate postID in user collection: %v", err)
	}
	return insert, nil
}

// DeletePost deletes a post document from the database
func (db DbConn) DeletePost(id string) (*mongo.DeleteResult, error) {
	newID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("postID %v is not a valid", id)
	}
	post := models.Post{}
	coll := db.Db.Collection("posts")

	// first get the post document and decode into post struct
	filter := bson.D{primitive.E{Key: "_id", Value: newID}}
	find := coll.FindOne(ctx, filter)
	if err := find.Decode(&post); err != nil {
		return nil, fmt.Errorf("error decoding post: %v", err)
	}

	coll2 := db.Db.Collection("users")

	// remove the post id from the user's collection.
	filter2 := bson.M{"username": post.Username}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "post_id", Value: newID.Hex()}}}}
	_, err = coll2.UpdateOne(ctx, filter2, update)
	if err != nil {
		return nil, fmt.Errorf("unable to update user's collection: %v", err)
	}

	// delete the post document from the post collection
	delete, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return nil, ErrConn
	}

	return delete, nil

}

// UpdatePost updates an existing post
func (db DbConn) UpdatePost(w http.ResponseWriter, id string, p models.Post) error {
	newID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("postID %v is not a valid", id)
	}
	coll := db.Db.Collection("posts")
	filter := bson.M{"_id": newID}
	//filter := bson.D{primitive.E{Key: "_id", Value: newID}}
	if p.Title != "" {
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "title", Value: p.Title}}}}
		UpdateResult, err := coll.UpdateOne(ctx, filter, update)
		if err != nil {
			return ErrConn
		}

		if UpdateResult.MatchedCount == 0 {
			json.NewEncoder(w).Encode("PostID does not match an existing post.")
		} else {
			json.NewEncoder(w).Encode(UpdateResult)
		}

	}

	if p.Body != "" {
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "body", Value: p.Body}}}}
		UpdateResult, err := coll.UpdateOne(ctx, filter, update)
		if err != nil {
			return ErrConn
		}

		if UpdateResult.MatchedCount == 0 {
			json.NewEncoder(w).Encode("PostID does not match an existing post.")
		} else {
			json.NewEncoder(w).Encode(UpdateResult)
		}
	}

	return nil
}

// Comment creates(adds) a comment field to a post document in the database
func (db DbConn) Comment(id string, com models.Comment) (*mongo.UpdateResult, error) {
	coll := db.Db.Collection("posts")
	PostId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("id is not a valid post id")
	}
	com.ID = primitive.NewObjectID()
	filter := bson.M{"_id": PostId}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "comments", Value: com}}}}
	upload, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return upload, nil
}

// GetPost returns a post document from the database.
func (db DbConn) GetPost(id string) *mongo.SingleResult {
	newID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil
	}
	coll := db.Db.Collection("posts")
	filter := bson.D{primitive.E{Key: "_id", Value: newID}}
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

	// Delete the posts from the users collection as well

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
		return nil, err
	}

	return find, nil
}

// DeleteComment deletes a comment from a post document.
func (db *DbConn) DeleteComment(id string, user string) (*mongo.UpdateResult, error) {
	coll := db.Db.Collection("posts")
	postID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("id is not valid")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: postID}}
	update := bson.M{"$pull": bson.M{"comments": bson.M{"username": user}}}

	delete, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return delete, nil
}
