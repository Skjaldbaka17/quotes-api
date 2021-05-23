//https://github.com/nicholasjackson/building-microservices-youtube/blob/episode_7/product-api/main.go
package main

import (
	"net/http"

	"github.com/Skjaldbaka17/quotes-api/handlers"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	gets := r.Methods(http.MethodGet).Subrouter()
	posts := r.Methods(http.MethodPost).Subrouter()
	posts.HandleFunc("/api/quotes", handlers.GetQuotesById)
	gets.HandleFunc("/api/authors/{id}", handlers.GetAuthorById)
	gets.HandleFunc("/api/quotes/{id}", handlers.GetQuoteById)
	gets.HandleFunc("/api/search/{searchString}", handlers.SearchByString)
	gets.HandleFunc("/api/search/authors/{searchString}", handlers.SearchAuthorsByString)
	gets.HandleFunc("/api/search/quotes/{searchString}", handlers.SearchQuotesByString)
	http.ListenAndServe(":8080", r)
}
