package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Firstname string             `json:"firstname" bson:"firstname"`
	Lastname  string             `json:"lastname" bson:"lastname"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Post_Id   int                `json:"post_id" bson:"post_id"`
	Id        primitive.ObjectID `json:"id" bson:"_id"`
}
