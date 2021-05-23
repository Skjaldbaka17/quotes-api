package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Quotes struct {
	gorm.Model
	Id          int    `json:"id" gorm:"primary_key"`
	Authorid    int    `json:"authorid"`
	Quote       string `json:"quote"`
	Count       int    `json:"count"`
	IsIcelandic bool   `json:"isicelandic"`
}

type Authors struct {
	gorm.Model
	Name   string   `json:"name"`
	Count  int      `json:"count"`
	Quotes []Quotes `gorm:"foreignKey:authorid;references:Id" json:"quotes" `
}

type SearchView struct {
	//if AuthorId then gorm cant map the values correctly, but works with Authorid and Quoteid etc. Why? TODO
	Authorid    int    `json:"authorid"`
	Name        string `json:"name"`
	Quoteid     int    `json:"quoteid" `
	Quote       string `json:"quote"`
	Isicelandic bool   `json:"isicelandic"`
}

type QuotesRequest struct {
	Ids []int `json:"ids"`
}

var db *gorm.DB

func GetAuthorById(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var author []SearchView
	err := db.Table("searchview").
		Select("*, authorid").
		Where("authorid = ?", params["id"]).
		Find(&author).
		Error

	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&author)
}

func GetQuoteById(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var quote SearchView

	err := db.
		Table("searchview").
		Where("quoteid = ?", params["id"]).
		First(&quote).
		Error

	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&quote)
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

func SearchByString(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	searchString := params["searchString"]
	fmt.Print(searchString)

	var results []SearchView

	err := db.Table("searchview").
		Where("quote @@ plainto_tsquery(?)", searchString).
		Select("*, ts_rank(to_tsvector(name), plainto_tsquery(?)) as rank", searchString).
		Or("name @@ plainto_tsquery(?)", searchString).
		Or("? % ANY(STRING_TO_ARRAY(name,' '))", searchString).
		Order(gorm.Expr("Similarity(name, ?) DESC", searchString)).
		Limit(25).
		Find(&results).Error

	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&results)
}

func SearchAuthorsByString(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	searchString := params["searchString"]
	fmt.Print(searchString)

	var results []SearchView

	m1 := regexp.MustCompile(` `)
	search := m1.ReplaceAllString(searchString, " <-> ")

	//% is same as SIMILARITY but with default threshold 0.3
	err := db.Table("searchview").
		Where("name @@ to_tsquery(?)", search).
		Or("? % ANY(STRING_TO_ARRAY(name,' '))", search).
		Order(gorm.Expr("Similarity(name, ?) DESC", searchString)).
		Limit(25).
		Find(&results).Error

	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&results)
}

func SearchQuotesByString(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	searchString := params["searchString"]
	fmt.Print(searchString)

	var results []SearchView

	err := db.Table("searchview").
		Select("*, ts_rank(to_tsvector(quote), plainto_tsquery(?)) as rank", searchString).
		Where("quote @@ plainto_tsquery(?)", searchString).
		Order("rank DESC").
		Limit(25).
		Find(&results).Error

	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&results)
}
