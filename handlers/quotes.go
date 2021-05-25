package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

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

type Request struct {
	Ids          []int  `json:"ids,omitempty"`
	Id           int    `json:"id,omitempty"`
	Page         int    `json:"page,omitempty"`
	SearchString string `json:"searchString,omitempty"`
	PageSize     int    `json:"pageSize,omitempty"`
}

const defaultPageSize = 25
const maxPageSize = 200

func GetAuthorsById(rw http.ResponseWriter, r *http.Request) {
	requestBody, err := validateRequestBody(r)

	if err != nil {
		http.Error(rw, "Could not finish", 404)
		return
	}
	var authors []SearchView
	err = db.Table("searchview").
		Select("*").
		Where("authorid in (?)", requestBody.Ids).
		First(&authors).
		Error
	log.Println(authors)
	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&authors)
}

// func GetAuthorsById(rw http.ResponseWriter, r *http.Request) {
// 	var requestBody Request
// 	err := json.NewDecoder(r.Body).Decode(&requestBody)

// 	if err != nil {
// 		log.Printf("Got error when decoding: %s", err)
// 		http.Error(rw, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	var quotes []SearchView
// 	err = db.Table("searchview").
// 		Select("*").
// 		Where("authorid in ?", requestBody.Ids).
// 		Order("quoteid ASC").
// 		Find(&quotes).
// 		Error

// 	if err != nil {
// 		log.Printf("Got error when decoding: %s", err)
// 		http.Error(rw, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	json.NewEncoder(rw).Encode(&quotes)
// }

func GetQuotesById(rw http.ResponseWriter, r *http.Request) {
	var requestBody Request
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	log.Println(requestBody)
	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var quotes []SearchView
	err = db.Table("searchview").
		Select("*").
		Where("quoteid in ?", requestBody.Ids).
		Order("quoteid ASC").
		Find(&quotes).
		Error

	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&quotes)
}

func validateRequestBody(r *http.Request) (Request, error) {
	var requestBody Request
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		return Request{}, err
	}

	//TODO: add validation for searchString and page etc.

	if requestBody.PageSize == 0 || requestBody.PageSize > maxPageSize {
		requestBody.PageSize = defaultPageSize
	}
	return requestBody, nil
}

func SearchByString(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()

	requestBody, err := validateRequestBody(r)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	var results []SearchView
	m1 := regexp.MustCompile(` `)
	phrasesearch := m1.ReplaceAllString(requestBody.SearchString, " <-> ")
	generalsearch := m1.ReplaceAllString(requestBody.SearchString, " | ")

	//Order by authorid to have definitive order (when for examplke some quotes rank the same for plain, phrase, general and similarity)
	err = db.Table("searchview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq,to_tsquery(?) as generalq ",
		requestBody.SearchString, phrasesearch, generalsearch).
		Where("tsv @@ plainq").
		Or("tsv @@ phraseq").
		Or("? % ANY(STRING_TO_ARRAY(name,' '))", requestBody.SearchString).
		Select("*, ts_rank(quotetsv, plainq) as plainrank, ts_rank(quotetsv, phraseq) as phraserank, ts_rank(quotetsv, generalq) as generalrank").
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "phraserank DESC,similarity(name, ?) DESC, plainrank DESC, generalrank DESC, authorid DESC", Vars: []interface{}{requestBody.SearchString}, WithoutParentheses: true},
		}).
		Or("tsv @@ generalq").
		Limit(requestBody.PageSize).
		Offset(requestBody.Page * requestBody.PageSize).
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

	requestBody, err := validateRequestBody(r)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	var results []SearchView

	//Order by authorid to have definitive order (when for examplke some names rank the same for similarity), same for why quoteid
	//% is same as SIMILARITY but with default threshold 0.3
	err = db.Table("searchview").
		Where("nametsv @@ plainto_tsquery(?)", requestBody.SearchString).
		Or("? % ANY(STRING_TO_ARRAY(name,' '))", requestBody.SearchString).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "similarity(name, ?) DESC, authorid DESC, quoteid DESC", Vars: []interface{}{requestBody.SearchString}, WithoutParentheses: true},
		}).
		Limit(requestBody.PageSize).
		Offset(requestBody.Page * requestBody.PageSize).
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
	requestBody, err := validateRequestBody(r)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	var results []SearchView
	m1 := regexp.MustCompile(` `)
	phrasesearch := m1.ReplaceAllString(requestBody.SearchString, " <-> ")
	generalsearch := m1.ReplaceAllString(requestBody.SearchString, " | ")

	//Order by quoteid to have definitive order (when for examplke some quotes rank the same for plain, phrase and general)
	err = db.Table("searchview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq,to_tsquery(?) as generalq ",
		requestBody.SearchString, phrasesearch, generalsearch).
		Where("quotetsv @@ plainq").
		Or("quotetsv @@ phraseq").
		Select("*, ts_rank(quotetsv, plainq) as plainrank, ts_rank(quotetsv, phraseq) as phraserank, ts_rank(quotetsv, generalq) as generalrank").
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "plainrank DESC, phraserank DESC, generalrank DESC, quoteid DESC", Vars: []interface{}{}, WithoutParentheses: true},
		}).
		Or("quotetsv @@ generalq").
		Limit(requestBody.PageSize).
		Offset(requestBody.Page * requestBody.PageSize).
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
