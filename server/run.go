package server

import (
	"bookstore/helper"
	"log"
	"net/http"
)

var repository = *helper.NewRepository()

func Run() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8081", router))
}
