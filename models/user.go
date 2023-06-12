package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Firstname string             `json:"firstname" bson:"firstname" validate:"required"`
	Lastname  string             `json:"lastname" bson:"lastname" validate:"omitempty"`
	Username  string             `json:"username" bson:"username" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Password  string             `json:"password" bson:"password" validate:"required,min=8"`
	Post_Id   []Post_ID          `json:"posts" bson:"posts"`
	Id        primitive.ObjectID `json:"id" bson:"_id"`
}

type Post_ID struct {
	Post_id string `json:"post_id" bson:"post_id"`
}

type Password struct {
	Old     string `json:"old password"`
	New     string `json:"new password"`
	Confirm string `json:"confirm password"`
}
