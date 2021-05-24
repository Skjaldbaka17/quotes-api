package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Quotes struct {
	gorm.Model
	Id          int    `json:"id"`
	Authorid    int    `json:"authorid"`
	Quote       string `json:"quote"`
	Count       int    `json:"count"`
	IsIcelandic bool   `json:"isicelandic"`
}

type Authors struct {
	gorm.Model
	Name   string   `json:"name"`
	Count  int      `json:"count"`
	Quotes []Quotes `json:"quotes" gorm:"foreignKey:authorid"`
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
	start := time.Now()
	params := mux.Vars(r)
	searchString := params["searchString"]

	var results []SearchView
	m1 := regexp.MustCompile(` `)
	phrasesearch := m1.ReplaceAllString(searchString, " <-> ")
	generalsearch := m1.ReplaceAllString(searchString, " | ")

	err := db.Table("searchview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq,to_tsquery(?) as generalq ",
		searchString, phrasesearch, generalsearch).
		Where("tsv @@ plainq").
		Or("tsv @@ phraseq").
		Select("*, ts_rank(tsv, plainq || generalq || phraseq) as plainrank").
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "plainrank DESC, similarity(name, ?) DESC", Vars: []interface{}{searchString}, WithoutParentheses: true},
		}).
		Or("tsv @@ generalq").
		Limit(25).
		Find(&results).Error

	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&results)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Time: %d", elapsed.Milliseconds())
}

func SearchAuthorsByString(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()
	params := mux.Vars(r)
	searchString := params["searchString"]
	fmt.Print(searchString)

	var results []SearchView

	//% is same as SIMILARITY but with default threshold 0.3
	err := db.Table("searchview").
		Where("nametsv @@ plainto_tsquery(?)", searchString).
		Or("? % ANY(STRING_TO_ARRAY(name,' '))", searchString).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "similarity(name, ?) DESC", Vars: []interface{}{searchString}, WithoutParentheses: true},
		}).
		Limit(25).
		Find(&results).Error

	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&results)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Time: %d", elapsed.Milliseconds())
}

func SearchQuotesByString(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()
	params := mux.Vars(r)
	searchString := params["searchString"]
	fmt.Print(searchString)

	var results []SearchView
	m1 := regexp.MustCompile(` `)
	phrasesearch := m1.ReplaceAllString(searchString, " <-> ")
	generalsearch := m1.ReplaceAllString(searchString, " | ")

	err := db.Table("searchview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq,to_tsquery(?) as generalq ",
		searchString, phrasesearch, generalsearch).
		Where("quotetsv @@ plainq").
		Or("quotetsv @@ phraseq").
		Select("*, ts_rank(quotetsv, plainq) as plainrank, ts_rank(quotetsv, phraseq) as phraserank, ts_rank(quotetsv, generalq) as generalrank").
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "plainrank DESC, phraserank DESC, generalrank DESC", Vars: []interface{}{}, WithoutParentheses: true},
		}).
		Or("quotetsv @@ generalq").
		Limit(25).
		Find(&results).Error

	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&results)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Time: %d", elapsed.Milliseconds())

}
