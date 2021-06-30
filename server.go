//https://github.com/nicholasjackson/building-microservices-youtube/blob/episode_7/product-api/main.go
package main

import (
	"net/http"

	"github.com/Skjaldbaka17/quotes-api/handlers"
	"github.com/Skjaldbaka17/quotes-api/routes"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	posts := r.Methods(http.MethodPost).Subrouter()
	posts.HandleFunc("/api/quotes", routes.GetQuotes)
	posts.HandleFunc("/api/quotes/list", routes.GetQuotesList)
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

	posts.HandleFunc("/api/users/signup", routes.CreateUser)
	posts.HandleFunc("/api/users/login", routes.Login)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	gets := r.Methods(http.MethodGet).Subrouter()
	gets.HandleFunc("/api/meta/languages", routes.ListLanguagesSupported)
	gets.Handle("/docs", sh)
	gets.Handle("/swagger/swagger.yaml", http.FileServer(http.Dir("./")))

	r.HandleFunc("/", routes.Home)
	s := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))
	r.PathPrefix("/assets/").Handler(s)

	err := http.ListenAndServe(":"+handlers.GetEnvVariable("PORT"), r)

	if err != nil {
		panic(err)
	}
}
