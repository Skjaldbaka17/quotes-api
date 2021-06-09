//https://github.com/nicholasjackson/building-microservices-youtube/blob/episode_7/product-api/main.go
package main

import (
	"net/http"

	"github.com/Skjaldbaka17/quotes-api/routes"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	posts := r.Methods(http.MethodPost).Subrouter()
	posts.HandleFunc("/api/quotes", routes.GetQuotes)
	posts.HandleFunc("/api/quotes/random", routes.GetRandomQuote)
	posts.HandleFunc("/api/quotes/qod/new", routes.SetQuoteOfTheDay)
	posts.HandleFunc("/api/quotes/qod", routes.GetQuoteOfTheDay)
	posts.HandleFunc("/api/quotes/qod/history", routes.GetQODHistory)

	posts.HandleFunc("/api/search", routes.SearchByString)
	posts.HandleFunc("/api/search/authors", routes.SearchAuthorsByString)
	posts.HandleFunc("/api/search/quotes", routes.SearchQuotesByString)

	posts.HandleFunc("/api/authors", routes.GetAuthorsById)
	posts.HandleFunc("/api/authors/list", routes.GetAuthorsList)
	posts.HandleFunc("/api/authors/random", routes.GetRandomAuthor)
	posts.HandleFunc("/api/authors/aod/new", routes.SetAuthorOfTheDay)
	posts.HandleFunc("/api/authors/aod", routes.GetAuthorOfTheDay)
	posts.HandleFunc("/api/authors/aod/history", routes.GetAODHistory)

	posts.HandleFunc("/api/topics", routes.GetTopics)
	posts.HandleFunc("/api/topic", routes.GetTopic)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	gets := r.Methods(http.MethodGet).Subrouter()
	gets.HandleFunc("/api/languages", routes.ListLanguagesSupported)
	gets.Handle("/docs", sh)
	gets.Handle("/swagger/swagger.yaml", http.FileServer(http.Dir("./")))

	http.ListenAndServe(":8080", r)
}
