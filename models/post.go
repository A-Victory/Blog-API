package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	Username  string `json:"username" bson:"username"`
	Header    string `json:"header" bson:"header"`
	Body      string `json:"post_body" bson:"post_body"`
	Reactions []Reaction
	Comments  []Comment
	User_id   int                `json:"user_id" bson:"user_id"`
	Id        primitive.ObjectID `json:"id" bson:"_id"`
}

type Comment struct {
	Username string `json:"username" bson:"username"`
	Body     string `json:"comment_body" bson:"comment_body"`
}

type Reaction struct {
	Upvote   string `json:"upvote" bson:"upvote"`
	Downvote string `json:"downvote" bson:"downvote"`
}
