package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Book struct {
	ID     primitive.ObjectID `json:"_id, omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	Author string             `json:"author,omitempty" bson:"author,omitempty"`
}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/book", CreateBookEndpoint).Methods("POST")
	router.HandleFunc("/books", GetBooksEndpoint).Methods("GET")
	router.HandleFunc("/book/{id}", GetBookEndpoint).Methods("GET")
	router.HandleFunc("/book/{id}", UpdateBookEndpoint).Methods("PUT")
	router.HandleFunc("/book/{id}", DeleteBookEndpoint).Methods("DELETE")
	http.ListenAndServe(":12345", router)
}

func CreateBookEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var book Book
	_ = json.NewDecoder(request.Body).Decode(&book)
	collection := client.Database("benfica").Collection("books")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, book)
	json.NewEncoder(response).Encode(result)
}

func GetBookEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var book Book
	collection := client.Database("benfica").Collection("books")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Book{ID: id}).Decode(&book)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(book)
}

func GetBooksEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var books []Book
	collection := client.Database("benfica").Collection("books")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var book Book
		cursor.Decode(&book)
		books = append(books, book)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	}
	json.NewEncoder(response).Encode(books)
}

func UpdateBookEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var book Book
	collection := client.Database("benfica").Collection("books")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	filter := bson.M{"_id": id}
	_ = json.NewDecoder(request.Body).Decode(&book)

	update := bson.D{
		{"$set", bson.D{
			{"name", book.Name},
			{"author", book.Author},
		}},
	}

	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&book)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	book.ID = id
	json.NewEncoder(response).Encode(book)

}

func DeleteBookEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database("benfica").Collection("books")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	filter := bson.M{"_id": id}
	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(deleteResult)
}
