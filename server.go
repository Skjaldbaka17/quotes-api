//https://github.com/nicholasjackson/building-microservices-youtube/blob/episode_7/product-api/main.go
package main

import (
	"net/http"

	"github.com/Skjaldbaka17/quotes-api/handlers"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	posts := r.Methods(http.MethodPost).Subrouter()
	posts.HandleFunc("/api/quotes", handlers.GetQuotesById)
	posts.HandleFunc("/api/quotes/random", handlers.GetRandomQuote)
	posts.HandleFunc("/api/search", handlers.SearchByString)
	posts.HandleFunc("/api/search/authors", handlers.SearchAuthorsByString)
	posts.HandleFunc("/api/search/quotes", handlers.SearchQuotesByString)
	posts.HandleFunc("/api/authors", handlers.GetAuthorsById)
	posts.HandleFunc("/api/authors/list", handlers.GetAuthorsList)
	posts.HandleFunc("/api/authors/random", handlers.GetRandomAuthor)
	posts.HandleFunc("/api/topics", handlers.GetTopics)
	posts.HandleFunc("/api/topic", handlers.GetTopic)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	gets := r.Methods(http.MethodGet).Subrouter()
	gets.HandleFunc("/api/languages", handlers.ListLanguagesSupported)
	gets.Handle("/docs", sh)
	gets.Handle("/swagger/swagger.yaml", http.FileServer(http.Dir("./")))

	http.ListenAndServe(":8080", r)
}
