//https://github.com/nicholasjackson/building-microservices-youtube/blob/episode_7/product-api/main.go
package main

import (
	"net/http"

	"github.com/Skjaldbaka17/quotes-api/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.GetQuotesById)
	// http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
