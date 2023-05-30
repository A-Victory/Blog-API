package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Firstname string             `json:"firstname" bson:"firstname" validate:"required"`
	Lastname  string             `json:"lastname" bson:"lastname" validate:"omitempty"`
	Username  string             `json:"username" bson:"username" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Password  string             `json:"password" bson:"password" validate:"required,min=8,alphanum"`
	Post_Id   []string           `json:"post_id" bson:"post_id"`
	Id        primitive.ObjectID `json:"id" bson:"_id"`
}
