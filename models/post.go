package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	Username   string             `json:"username" bson:"username"`
	Title      string             `json:"title" bson:"title"`
	Body       string             `json:"body" bson:"body"`
	Reactions  Reaction           `json:"reactions" bson:"reactions"`
	Comments   []Comment          `json:"comments" bson:"comments"`
	Created_At time.Time          `json:"created_at" bson:"created_at"`
	User_id    string             `json:"user_id" bson:"user_id"`
	Id         primitive.ObjectID `json:"id" bson:"_id"`
}

type Comment struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Body     string             `json:"body" bson:"body"`
}

type Reaction struct {
	UpVote   []UpVote   `json:"upvote" bson:"upvote"`
	DownVote []DownVote `json:"downvote" bson:"downvote"`
}

type UpVote struct {
	Id string `json:"user_id" bson:"user_id"`
}

type DownVote struct {
	Id string `json:"user_id" bson:"user_id"`
}
