package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"gorm.io/gorm/clause"
)

const defaultPageSize = 25
const maxPageSize = 200

// swagger:route POST /authors AUTHORS getAuthorsByIds
//
// Get authors by their ids
//
// responses:
//	200: multipleQuotesResponse

// Get Authors handles POST requests to get the authors, and their quotes, that have the given ids
func GetAuthorsById(rw http.ResponseWriter, r *http.Request) {
	requestBody, err := validateRequestBody(r)

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		http.Error(rw, "Could not finish", 404)
		return
	}
	var authors []SearchView
	err = db.Table("searchview").
		Select("*").
		Where("authorid in (?)", requestBody.Ids).
		Find(&authors).
		Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&authors)
}

// swagger:route POST /quotes QUOTES getQuotesByIds
// Get quotes by their ids
//
// responses:
//	200: multipleQuotesResponse

// GetQuotesById handles POST requests to get the quotes, and their authors, that have the given ids
func GetQuotesById(rw http.ResponseWriter, r *http.Request) {
	var requestBody Request
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	log.Println(requestBody)
	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
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
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&quotes)
}

//ValidateRequestBody takes in the request and validates all the input fields, returns an error with reason for validation-failure
//if validation fails.
//TODO: Make validation better! i.e. make it "real"
func validateRequestBody(r *http.Request) (Request, error) {
	var requestBody Request
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		//TODO: Respond with better error -- and add tests
		log.Printf("Got error when decoding: %s", err)
		return Request{}, err
	}

	//TODO: add validation for searchString and page etc.

	if requestBody.PageSize == 0 || requestBody.PageSize > maxPageSize {
		requestBody.PageSize = defaultPageSize
	}
	return requestBody, nil
}

// swagger:route POST /search SEARCH generalSearchByString
// Search for quotes / authors by a general string-search that searches both in the names of the authors and the quotes themselves
//
// responses:
//	200: multipleQuotesResponse

// SearchByString handles POST requests to search for quotes / authors by a search-string
func SearchByString(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()

	requestBody, err := validateRequestBody(r)

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
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
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&results)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Time: %d", elapsed.Milliseconds())
}

// swagger:route POST /search/authors SEARCH searchAuthorsByString
//
// Authors search. Searching authors by a given search string
//
// responses:
//	200: multipleQuotesResponse

// SearchAuthorsByString handles POST requests to search for authors by a search-string
func SearchAuthorsByString(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()

	requestBody, err := validateRequestBody(r)

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
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
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&results)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Time: %d", elapsed.Milliseconds())
}

// swagger:route POST /search/quotes SEARCH searchQuotesByString
// Quotes search. Searching quotes by a given search string
// responses:
//	200: multipleQuotesResponse

// SearchQuotesByString handles POST requests to search for quotes by a search-string
func SearchQuotesByString(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestBody, err := validateRequestBody(r)

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
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
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&results)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Time: %d", elapsed.Milliseconds())
}

// swagger:route POST /topics TOPICS getTopics
// List the available topics, english / icelandic or both
// responses:
//	200: listTopicsResponse

// GetTopics handles POST requests for listing the available quote-topics
func GetTopics(rw http.ResponseWriter, r *http.Request) {
	requestBody, err := validateRequestBody(r)
	log.Println("HERE")
	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	var results []ListItem

	pointer := db.Table("topics")

	if requestBody.Language == "English" {
		pointer = pointer.Where("not isicelandic")
	} else if requestBody.Language == "Icelandic" {
		pointer = pointer.Where("isicelandic")
	}

	err = pointer.Find(&results).Error
	log.Println(results)
	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&results)
}

// swagger:route POST /topic TOPICS getTopic
// Get quotes from a particular topic
// responses:
//	200: multipleQuotesTopicResponse

// GetTopic handles POST requests for getting quotes from a particular topic
func GetTopic(rw http.ResponseWriter, r *http.Request) {
	requestBody, err := validateRequestBody(r)

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	var results []TopicsView

	//Order by quoteid to have definitive order (when for examplke some quotes rank the same for plain, phrase and general)
	dbPoint := db.Table("topicsview")

	if requestBody.Topic != "" {
		dbPoint = dbPoint.Where("lower(topicname) = lower(?)", requestBody.Topic)
	} else {
		dbPoint = dbPoint.Where("topicid = ?", requestBody.Id)
	}

	err = dbPoint.Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "quoteid DESC", Vars: []interface{}{}, WithoutParentheses: true},
	}).
		Limit(requestBody.PageSize).
		Offset(requestBody.Page * requestBody.PageSize).
		Find(&results).Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&results)
}

// swagger:route GET /languages META getLanguages
// Get languages supported by the api
// responses:
//	200: listOfStrings

// GetTopic handles POST requests for getting quotes from a particular topic
func ListLanguagesSupported(rw http.ResponseWriter, r *http.Request) {

	type response = struct {
		Languages []string `json:"languages"`
	}

	json.NewEncoder(rw).Encode(&response{
		Languages: []string{"English", "Icelandic"},
	})
}
