package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Book struct {
	ID     bson.ObjectId `json:"_id, omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	Author string             `json:"author,omitempty" bson:"author,omitempty"`
}

type Config struct {
	MongoDBhost 	string 		`json:"MongoDBhost"`
	MongoDBname 	string 		`json:"MongoDBname"`
	MongoDBCollectionName 	string 		`json:"MongoDBCollectionName"`
}
