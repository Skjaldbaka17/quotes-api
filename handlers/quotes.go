package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Quotes1 struct {
	gorm.Model
	Id          int    `json:"id" gorm:"primary_key"`
	Author_Id   int    `json:"author_id"`
	Quote       string `json:"quote"`
	Count       int    `json:"count"`
	IsIcelandic bool   `json:"isicelandic"`
}

type Authors1 struct {
	gorm.Model
	Name   string    `json:"name"`
	Count  int       `json:"count"`
	Quotes []Quotes1 `gorm:"foreignKey:Author_Id;references:Id" json:"quotes" `
}

type QuotesRequest struct {
	Ids []int `json:"ids"`
}

var db *gorm.DB

func GetAuthorById(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var author Authors1
	var quotes []Quotes1
	db.First(&author, params["id"])
	db.Where("author_id = ?", author.ID).Find(&quotes)
	author.Quotes = quotes
	json.NewEncoder(rw).Encode(&author)
}

func GetQuotesById(rw http.ResponseWriter, r *http.Request) {
	var requestBody QuotesRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	requestBody.Ids = append(requestBody.Ids, 10)
	out, _ := json.Marshal(requestBody)
	rw.Write(out)
}
