package server

import (
	"bookstore/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	http "net/http"
	"strings"
	"gopkg.in/mgo.v2/bson"
	//"bookstore/helper"
)

func WriteHeader(writer http.ResponseWriter){
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	return
}

func getBooks(writer http.ResponseWriter, request *http.Request) {
	queryString := request.URL.Query().Get("q")
	//log.Println("URL : ", request.URL.RequestURI())
	data, err := repository.GetBooks(queryString)
	WriteHeader(writer)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(data)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}

func createBook(writer http.ResponseWriter, request *http.Request) {
	var book models.Book
	body, err := ioutil.ReadAll(io.LimitReader(request.Body,1048576))

	if err != nil {
		log.Println("Error book: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := request.Body.Close(); err != nil {
		log.Fatalln("Error createBook: ", err)
	}

	if err := json.Unmarshal(body, &book); err != nil {
		writer.WriteHeader(422)
		if err := json.NewEncoder(writer).Encode(err); err != nil {
			log.Fatalln("Error createBook unmarshalling data: ", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	success, _ := repository.CreateBook(book)
	if success == "0" {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write([]byte(success))
	writer.WriteHeader(http.StatusCreated)

}

func getBook(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]
	log.Println("URL : ", request.URL.RequestURI())
	book, err := repository.GetBook(id)
	WriteHeader(writer)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(book)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(book)
}

func updateBook(writer http.ResponseWriter, request *http.Request) {
	var book models.Book
	body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))

	if err != nil {
		log.Fatalln("Error Update Data: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := request.Body.Close(); err != nil {
		log.Fatalln("Error update data (Body.Close()): ", err)
	}

	if err := json.Unmarshal(body, &book); err != nil {
		WriteHeader(writer)
		writer.WriteHeader(422)
		if err := json.NewEncoder(writer).Encode(err); err != nil {
			log.Fatalln("Error updateBook unmarshaling data", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	vars := mux.Vars(request)
	book.ID = bson.ObjectIdHex(vars["id"])
	err = repository.UpdateBook(book)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteHeader(writer)
	writer.WriteHeader(http.StatusOK)
}

func deleteBook(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]
	if outString, err := repository.DeleteBook(id); err != nil {
		if strings.Contains(outString, "404") {
			writer.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(outString, "500"){
			writer.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	WriteHeader(writer)
	writer.WriteHeader(http.StatusOK)
}


