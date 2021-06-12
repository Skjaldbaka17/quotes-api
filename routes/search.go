package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/Skjaldbaka17/quotes-api/handlers"
	"github.com/Skjaldbaka17/quotes-api/structs"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// swagger:route POST /search SEARCH SearchByString
// Search for quotes / authors by a general string-search that searches both in the names of the authors and the quotes themselves
//
// responses:
//	200: quotesResponse

// SearchByString handles POST requests to search for quotes / authors by a search-string
func SearchByString(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var results []structs.QuoteView
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := getBasePointer(requestBody)
	//Order by authorid to have definitive order (when for examplke some quotes rank the same for plain, phrase, general and similarity)
	dbPointer = dbPointer.
		Where("( tsv @@ plainq OR tsv @@ phraseq OR ? % ANY(STRING_TO_ARRAY(name,' ')) OR tsv @@ generalq)", requestBody.SearchString).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "phraserank DESC,similarity(name, ?) DESC, plainrank DESC, generalrank DESC, authorid DESC", Vars: []interface{}{requestBody.SearchString}, WithoutParentheses: true},
		})

	//Particular language search
	dbPointer = quoteLanguageSQL(requestBody.Language, dbPointer)
	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err := pagination(requestBody, dbPointer).
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
}

// swagger:route POST /search/authors SEARCH SearchAuthorsByString
//
// Authors search. Searching authors by a given search string
//
// responses:
//	200: authorsResponse

// SearchAuthorsByString handles POST requests to search for authors by a search-string
func SearchAuthorsByString(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var results []structs.AuthorsView
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	//Order by authorid to have definitive order (when for examplke some names rank the same for similarity), same for why quoteid
	//% is same as SIMILARITY but with default threshold 0.3
	dbPointer := handlers.Db.Table("authors").
		Where("( tsv @@ plainto_tsquery(?) OR (?) % ANY(STRING_TO_ARRAY(name,' ')) )", requestBody.SearchString, requestBody.SearchString).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "similarity(name, ?) DESC, id DESC", Vars: []interface{}{requestBody.SearchString}, WithoutParentheses: true},
		})

	//Particular language search
	dbPointer = authorLanguageSQL(requestBody.Language, dbPointer)
	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err := pagination(requestBody, dbPointer).
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
}

// swagger:route POST /search/quotes SEARCH SearchQuotesByString
// Quotes search. Searching quotes by a given search string
// responses:
//	200: quotesResponse

// SearchQuotesByString handles POST requests to search for quotes by a search-string
func SearchQuotesByString(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var results []structs.QuoteView
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := getBasePointer(requestBody)
	dbPointer = dbPointer.Where("( quotetsv @@ plainq OR quotetsv @@ phraseq OR quotetsv @@ generalq)")

	if requestBody.AuthorId > 0 {
		dbPointer = dbPointer.Where("authorid = ?", requestBody.AuthorId)
	}

	//Order by quoteid to have definitive order (when for examplke some quotes rank the same for plain, phrase and general)
	dbPointer = dbPointer.
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "plainrank DESC, phraserank DESC, generalrank DESC, quoteid DESC", Vars: []interface{}{}, WithoutParentheses: true},
		})

	//Particular language search
	dbPointer = quoteLanguageSQL(requestBody.Language, dbPointer)
	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err := pagination(requestBody, dbPointer).
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
}

//getBasePointer returns a base DB pointer for a table for a thorough full text search
func getBasePointer(requestBody structs.Request) *gorm.DB {
	table := "searchview"
	//TODO: Validate that this topicId exists
	if requestBody.TopicId > 0 {
		table = "topicsview"
	}
	m1 := regexp.MustCompile(` `)
	phrasesearch := m1.ReplaceAllString(requestBody.SearchString, " <-> ")
	generalsearch := m1.ReplaceAllString(requestBody.SearchString, " | ")
	dbPointer := handlers.Db.Table(table+", plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq,to_tsquery(?) as generalq ",
		requestBody.SearchString, phrasesearch, generalsearch).Select("*, ts_rank(quotetsv, plainq) as plainrank, ts_rank(quotetsv, phraseq) as phraserank, ts_rank(quotetsv, generalq) as generalrank")

	if requestBody.TopicId > 0 {
		dbPointer = dbPointer.Where("topicid = ?", requestBody.TopicId)
	}
	return dbPointer
}

func pagination(requestBody structs.Request, dbPointer *gorm.DB) *gorm.DB {
	return dbPointer.Limit(requestBody.PageSize).
		Offset(requestBody.Page * requestBody.PageSize)
}
