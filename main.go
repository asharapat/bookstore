package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"bookstore/helper"
	"bookstore/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	fmt.Println("Starting the application...")

	router := mux.NewRouter()
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/book", createBook).Methods("POST")
	router.HandleFunc("/book/{id}", getBook).Methods("GET")
	router.HandleFunc("/book/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/book/{id}", deleteBook).Methods("DELETE")
	http.ListenAndServe(":12345", router)
}

func getBooks(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("content-type", "application/json")

	var books []models.Book

	collection := helper.ConnectDB()

	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, response)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var book models.Book
		err := cur.Decode(&book)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}

	json.NewEncoder(response).Encode(books)

}

func createBook(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "application/json")

	var book models.Book

	_ = json.NewDecoder(request.Body).Decode(&book)
	collection := helper.ConnectDB()
	result, err := collection.InsertOne(context.TODO(), book)
	if err != nil {
		helper.GetError(err, response)
		return
	}

	json.NewEncoder(response).Encode(result)
}

func getBook(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var book models.Book
	var params = mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := helper.ConnectDB()
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&book)
	if err != nil {
		helper.GetError(err, response)
		return
	}
	json.NewEncoder(response).Encode(book)
}

func updateBook(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var book models.Book
	collection := helper.ConnectDB()
	filter := bson.M{"_id": id}
	_ = json.NewDecoder(request.Body).Decode(&book)

	update := bson.D{
		{"$set", bson.D{
			{"name", book.Name},
			{"author", book.Author},
		}},
	}

	fmt.Println(book)

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&book)

	fmt.Println(book)

	if err != nil {
		helper.GetError(err, response)
		return
	}

	json.NewEncoder(response).Encode(update)

}

func deleteBook(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(request)
	id, err := primitive.ObjectIDFromHex(params["id"])
	collection := helper.ConnectDB()
	filter := bson.M{"_id": id}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		helper.GetError(err, response)
		return
	}
	json.NewEncoder(response).Encode(deleteResult)
}
