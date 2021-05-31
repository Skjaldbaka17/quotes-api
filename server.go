//https://github.com/nicholasjackson/building-microservices-youtube/blob/episode_7/product-api/main.go
package main

import (
	"net/http"

	"github.com/Skjaldbaka17/quotes-api/handlers"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func main() {
	requestHandler := handlers.Request{}
	validationChain := alice.New(requestHandler.BodyValidationHandler)

	r := mux.NewRouter()
	posts := r.Methods(http.MethodPost).Subrouter()
	posts.Handle("/api/quotes", validationChain.ThenFunc(requestHandler.GetQuotesById))
	posts.Handle("/api/quotes/random", validationChain.ThenFunc(requestHandler.GetRandomQuote))
	posts.Handle("/api/search", validationChain.ThenFunc(requestHandler.SearchByString))
	posts.Handle("/api/search/authors", validationChain.ThenFunc(requestHandler.SearchAuthorsByString))
	posts.Handle("/api/search/quotes", validationChain.ThenFunc(requestHandler.SearchQuotesByString))
	posts.Handle("/api/authors", validationChain.ThenFunc(requestHandler.GetAuthorsById))
	posts.Handle("/api/topics", validationChain.ThenFunc(requestHandler.GetTopics))
	posts.Handle("/api/topic", validationChain.ThenFunc(requestHandler.GetTopic))

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	gets := r.Methods(http.MethodGet).Subrouter()
	gets.HandleFunc("/api/languages", handlers.ListLanguagesSupported)
	gets.Handle("/docs", sh)
	gets.Handle("/swagger/swagger.yaml", http.FileServer(http.Dir("./")))

	http.ListenAndServe(":8080", r)
}
