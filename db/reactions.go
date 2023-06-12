package db

import (
	"github.com/A-Victory/blog-API/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Takes the post ID and updates the post documents
func (db DbConn) UpVote(postID, user string) error {

	post_ID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}

	userStruct := &models.User{}
	coll := db.Db.Collection("users")
	filter := bson.D{primitive.E{Key: "username", Value: user}}
	if err := coll.FindOne(ctx, filter).Decode(userStruct); err != nil {
		return err
	}

	vote := models.UpVote{}
	vote.Id = userStruct.Id.Hex()
	coll1 := db.Db.Collection("posts")
	filter1 := bson.M{"_id": post_ID}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "reactions.upvote", Value: vote}}}}
	_, err = coll1.UpdateOne(ctx, filter1, update)
	if err != nil {
		return err
	}

	return nil
}

// Takes the post ID and updates the reaction field in the post collection
func (db DbConn) DownVote(postID, user string) error {

	Post_ID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}

	userStruct := &models.User{}
	coll := db.Db.Collection("users")
	filter := bson.D{primitive.E{Key: "username", Value: user}}
	if err := coll.FindOne(ctx, filter).Decode(userStruct); err != nil {
		return err
	}

	vote := models.UpVote{}
	vote.Id = userStruct.Id.Hex()
	coll1 := db.Db.Collection("posts")
	filter1 := bson.D{primitive.E{Key: "_id", Value: Post_ID}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "reactions.downvote", Value: vote}}}}
	_, err = coll1.UpdateOne(ctx, filter1, update)
	if err != nil {
		return err
	}

	return nil
}
