package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func getSession() *mongo.Client {
	mongoString := os.Getenv("MONGO")
	cltOptns := options.Client().ApplyURI(mongoString)
	s, err := mongo.Connect(context.TODO(), cltOptns)

	if err != nil {
		panic(err)
	}
	if err := s.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	return s
}

func UserDb() *DbConn {
	db := getSession().Database("blog-API")
	return &DbConn{db}
}

type DbConn struct {
	Db *mongo.Database
}
