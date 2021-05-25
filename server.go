//https://github.com/nicholasjackson/building-microservices-youtube/blob/episode_7/product-api/main.go
package main

import (
	"net/http"

	"github.com/Skjaldbaka17/quotes-api/handlers"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	posts := r.Methods(http.MethodPost).Subrouter()
	posts.HandleFunc("/api/quotes", handlers.GetQuotesById)
	posts.HandleFunc("/api/search", handlers.SearchByString)
	posts.HandleFunc("/api/search/authors", handlers.SearchAuthorsByString)
	posts.HandleFunc("/api/search/quotes", handlers.SearchQuotesByString)
	posts.HandleFunc("/api/authors", handlers.GetAuthorsById)

	http.ListenAndServe(":8080", r)
}
