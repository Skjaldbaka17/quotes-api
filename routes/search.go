package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Skjaldbaka17/quotes-api/handlers"
	"github.com/Skjaldbaka17/quotes-api/structs"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// swagger:route POST /search SEARCH generalSearchByString
// Search for quotes / authors by a general string-search that searches both in the names of the authors and the quotes themselves
//
// responses:
//	200: multipleQuotesResponse

// SearchByString handles POST requests to search for quotes / authors by a search-string
func SearchByString(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	start := time.Now()

	var results []structs.QuoteView
	m1 := regexp.MustCompile(` `)
	phrasesearch := m1.ReplaceAllString(requestBody.SearchString, " <-> ")
	generalsearch := m1.ReplaceAllString(requestBody.SearchString, " | ")
	var dbPointer *gorm.DB
	//TODO: Validate that this topicId exists
	if requestBody.TopicId > 0 {
		dbPointer = handlers.Db.Table("topicsview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq,to_tsquery(?) as generalq ",
			requestBody.SearchString, phrasesearch, generalsearch)
	} else {
		dbPointer = handlers.Db.Table("searchview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq,to_tsquery(?) as generalq ",
			requestBody.SearchString, phrasesearch, generalsearch)
	}

	dbPointer = dbPointer.
		Select("*, ts_rank(quotetsv, plainq) as plainrank, ts_rank(quotetsv, phraseq) as phraserank, ts_rank(quotetsv, generalq) as generalrank").
		Where("( tsv @@ plainq OR tsv @@ phraseq OR ? % ANY(STRING_TO_ARRAY(name,' ')) OR tsv @@ generalq)", requestBody.SearchString)

	if requestBody.TopicId > 0 {
		dbPointer = dbPointer.Where("topicid = ?", requestBody.TopicId)
	}

	//Order by authorid to have definitive order (when for examplke some quotes rank the same for plain, phrase, general and similarity)
	dbPointer = dbPointer.Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "phraserank DESC,similarity(name, ?) DESC, plainrank DESC, generalrank DESC, authorid DESC", Vars: []interface{}{requestBody.SearchString}, WithoutParentheses: true},
	})

	//Particular language search
	switch strings.ToLower(requestBody.Language) {
	case "english":
		dbPointer = dbPointer.Not("isicelandic")
	case "icelandic":
		dbPointer = dbPointer.Where("isicelandic")
	}

	err := dbPointer.Limit(requestBody.PageSize).
		Offset(requestBody.Page * requestBody.PageSize).
		Find(&results).Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//Update popularity in background!
	go handlers.AppearInSearchCountIncrement(results)

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
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	start := time.Now()

	var results []structs.AuthorsView

	//Order by authorid to have definitive order (when for examplke some names rank the same for similarity), same for why quoteid
	//% is same as SIMILARITY but with default threshold 0.3
	dbPointer := handlers.Db.Table("authors").
		Where("( tsv @@ plainto_tsquery(?) OR ? % ANY(STRING_TO_ARRAY(name,' ')) )", requestBody.SearchString, requestBody.SearchString).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "similarity(name, ?) DESC, id DESC", Vars: []interface{}{requestBody.SearchString}, WithoutParentheses: true},
		})

		//Particular language search
	switch strings.ToLower(requestBody.Language) {
	case "english":
		dbPointer = dbPointer.Not("hasicelandicquotes")
	case "icelandic":
		dbPointer = dbPointer.Where("hasicelandicquotes")
	}

	err := dbPointer.Limit(requestBody.PageSize).
		Offset(requestBody.Page * requestBody.PageSize).
		Find(&results).Error
	if err != nil {

		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//Update popularity in background!
	go handlers.AuthorsAppearInSearchCountIncrement(results)

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
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	start := time.Now()

	var results []structs.QuoteView
	m1 := regexp.MustCompile(` `)
	phrasesearch := m1.ReplaceAllString(requestBody.SearchString, " <-> ")
	generalsearch := m1.ReplaceAllString(requestBody.SearchString, " | ")
	var dbPointer *gorm.DB
	//TODO: Validate that this topicId exists
	if requestBody.TopicId > 0 {
		dbPointer = handlers.Db.Table("topicsview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq,to_tsquery(?) as generalq ",
			requestBody.SearchString, phrasesearch, generalsearch)
	} else {
		dbPointer = handlers.Db.Table("searchview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq,to_tsquery(?) as generalq ",
			requestBody.SearchString, phrasesearch, generalsearch)
	}

	dbPointer = dbPointer.Select("*, ts_rank(quotetsv, plainq) as plainrank, ts_rank(quotetsv, phraseq) as phraserank, ts_rank(quotetsv, generalq) as generalrank").
		Where("( quotetsv @@ plainq OR quotetsv @@ phraseq OR quotetsv @@ generalq)")

	if requestBody.TopicId > 0 {
		dbPointer = dbPointer.Where("topicid = ?", requestBody.TopicId)
	}

	//Order by quoteid to have definitive order (when for examplke some quotes rank the same for plain, phrase and general)
	dbPointer = dbPointer.
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "plainrank DESC, phraserank DESC, generalrank DESC, quoteid DESC", Vars: []interface{}{}, WithoutParentheses: true},
		})

	//Particular language search
	switch strings.ToLower(requestBody.Language) {
	case "english":
		dbPointer = dbPointer.Not("isicelandic")
	case "icelandic":
		dbPointer = dbPointer.Where("isicelandic")
	}

	err := dbPointer.Limit(requestBody.PageSize).
		Offset(requestBody.Page * requestBody.PageSize).
		Find(&results).Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//Update popularity in background!
	go handlers.AppearInSearchCountIncrement(results)

	json.NewEncoder(rw).Encode(&results)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Time: %d", elapsed.Milliseconds())
}
