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
//  200: topicViewsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// SearchByString handles POST requests to search for quotes / authors by a search-string
func SearchByString(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var topicResults []structs.TopicViewDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := getBasePointer(requestBody)
	//Order by authorid to have definitive order (when for examplke some quotes rank the same for plain, phrase, general and similarity)
	dbPointer = dbPointer.
		Where("( tsv @@ plainq OR tsv @@ phraseq OR ? % ANY(STRING_TO_ARRAY(name,' ')) OR tsv @@ generalq)", requestBody.SearchString).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "phraserank DESC,similarity(name, ?) DESC, plainrank DESC, generalrank DESC, author_id DESC", Vars: []interface{}{requestBody.SearchString}, WithoutParentheses: true},
		})

	//Particular language search
	dbPointer = quoteLanguageSQL(requestBody.Language, dbPointer)
	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err := pagination(requestBody, dbPointer).
		Find(&topicResults).Error

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when querying DB in SearchByString: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	//Update popularity in background!
	go handlers.TopicViewAppearInSearchCountIncrement(topicResults)
	apiResults := structs.ConvertToTopicViewsAPIModel(topicResults)
	json.NewEncoder(rw).Encode(apiResults)
}

// swagger:route POST /search/authors SEARCH SearchAuthorsByString
//
// Authors search. Searching authors by a given search string
//
// responses:
//	200: authorsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// SearchAuthorsByString handles POST requests to search for authors by a search-string
func SearchAuthorsByString(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var results []structs.AuthorDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	//Order by authorid to have definitive order (when for examplke some names rank the same for similarity), same for why quote_id
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
	//  500: internalServerErrorResponse
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when querying DB in SearchAuthorsByString: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	//Update popularity in background!
	go handlers.AuthorsAppearInSearchCountIncrement(results)

	authorsAPI := structs.ConvertToAuthorsAPIModel(results)
	json.NewEncoder(rw).Encode(authorsAPI)
}

// swagger:route POST /search/quotes SEARCH SearchQuotesByString
// Quotes search. Searching quotes by a given search string
// responses:
//  200: topicViewsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// SearchQuotesByString handles POST requests to search for quotes by a search-string
func SearchQuotesByString(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var topicResults []structs.TopicViewDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := getBasePointer(requestBody)
	dbPointer = dbPointer.Where("( quote_tsv @@ plainq OR quote_tsv @@ phraseq OR quote_tsv @@ generalq)")

	if requestBody.AuthorId > 0 {
		dbPointer = dbPointer.Where("author_id = ?", requestBody.AuthorId)
	}

	//Order by quote_id to have definitive order (when for examplke some quotes rank the same for plain, phrase and general)
	dbPointer = dbPointer.
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "plainrank DESC, phraserank DESC, generalrank DESC, quote_id DESC", Vars: []interface{}{}, WithoutParentheses: true},
		})

	//Particular language search
	dbPointer = quoteLanguageSQL(requestBody.Language, dbPointer)
	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err := pagination(requestBody, dbPointer).
		Find(&topicResults).Error

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when querying DB in SearchQuotesByString: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	//Update popularity in background!
	go handlers.TopicViewAppearInSearchCountIncrement(topicResults)
	apiResults := structs.ConvertToTopicViewsAPIModel(topicResults)
	json.NewEncoder(rw).Encode(apiResults)

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
		requestBody.SearchString, phrasesearch, generalsearch).Select("*, ts_rank(quote_tsv, plainq) as plainrank, ts_rank(quote_tsv, phraseq) as phraserank, ts_rank(quote_tsv, generalq) as generalrank")

	if requestBody.TopicId > 0 {
		dbPointer = dbPointer.Where("topic_id = ?", requestBody.TopicId)
	}
	return dbPointer
}

func pagination(requestBody structs.Request, dbPointer *gorm.DB) *gorm.DB {
	return dbPointer.Limit(requestBody.PageSize).
		Offset(requestBody.Page * requestBody.PageSize)
}
