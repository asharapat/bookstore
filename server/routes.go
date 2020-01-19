package server

import (
	"net/http"
)


type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes {
	Route{
		"Get all books",
		"GET",
		"/books",
		getBooks,
	},
	Route{
		"Get book by id",
		"GET",
		"/books/{id}",
		getBook,
	},
	Route{
		"Add book",
		"POST",
		"/books",
		createBook,
	},
	Route{
		"Update book",
		"PUT",
		"/books/{id}",
		updateBook,
	},
	Route{
		"Delete book",
		"DELETE",
		"/books/{id}",
		deleteBook,
	},
}
