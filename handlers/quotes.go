package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const defaultPageSize = 25
const maxPageSize = 200
const maxQuotes = 50
const defaultMaxQuotes = 1

var languages = []string{"English", "Icelandic"}

//ValidateRequestBody takes in the request and validates all the input fields, returns an error with reason for validation-failure
//if validation fails.
//TODO: Make validation better! i.e. make it "real"
func getRequestBody(rw http.ResponseWriter, r *http.Request, requestBody *Request) error {
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		err = errors.New("request body is not structured correctly. Please refer to the /docs page for information on how to structure the request body")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(ErrorResponse{Message: err.Error()})
		return err
	}

	if requestBody.PageSize < 1 || requestBody.PageSize > maxPageSize {
		requestBody.PageSize = defaultPageSize
	}

	if requestBody.Page < 0 {
		requestBody.Page = 0
	}

	if requestBody.MaxQuotes < 0 || requestBody.MaxQuotes > maxQuotes {
		requestBody.MaxQuotes = maxQuotes
	}

	if requestBody.MaxQuotes <= 0 {
		requestBody.MaxQuotes = defaultMaxQuotes
	}

	return nil
}

// swagger:route POST /authors AUTHORS getAuthorsByIds
//
// Get authors by their ids
//
// responses:
//	200: authorsResponse

// Get Authors handles POST requests to get the authors, and their quotes, that have the given ids
func GetAuthorsById(rw http.ResponseWriter, r *http.Request) {
	var requestBody Request
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var authors []AuthorsView
	err := db.Table("authors").
		Where("id in (?)", requestBody.Ids).
		Scan(&authors).
		Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&authors)
}

// swagger:route POST /authors/list AUTHORS getAuthorsList
//
// Get list of authors according to some ordering / parameters
//
// responses:
//	200: authorsResponse

// GetAuthorsList handles POST requests to get the authors that fit the parameters

func GetAuthorsList(rw http.ResponseWriter, r *http.Request) {
	var requestBody Request
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var authors []AuthorsView
	dbPointer := db.Table("authors")

	switch strings.ToLower(requestBody.Language) {
	case "english":
		dbPointer = dbPointer.Not("hasicelandicquotes")
	case "icelandic":
		dbPointer = dbPointer.Where("hasicelandicquotes")
	}

	orderDirection := "ASC"
	if requestBody.OrderConfig.Reverse {
		orderDirection = "DESC"
	}

	switch strings.ToLower(requestBody.OrderConfig.OrderBy) {
	case "popularity": //TODO: add popularity ordering
	case "nrofquotes":
		switch strings.ToLower(requestBody.Language) {
		case "english":
			if nr, err := strconv.Atoi(requestBody.OrderConfig.Maximum); err == nil {
				dbPointer = dbPointer.Where("nrofenglishquotes <= ?", nr)
			}
			if nr, err := strconv.Atoi(requestBody.OrderConfig.Minimum); err == nil {
				dbPointer = dbPointer.Where("nrofenglishquotes >= ?", nr)
			}
			dbPointer = dbPointer.Order("nrofenglishquotes " + orderDirection)
		case "icelandic":
			if nr, err := strconv.Atoi(requestBody.OrderConfig.Maximum); err == nil {
				dbPointer = dbPointer.Where("nroficelandicquotes <= ?", nr)
			}
			if nr, err := strconv.Atoi(requestBody.OrderConfig.Minimum); err == nil {
				dbPointer = dbPointer.Where("nroficelandicquotes >= ?", nr)
			}
			dbPointer = dbPointer.Order("nroficelandicquotes " + orderDirection)
		default:
			if nr, err := strconv.Atoi(requestBody.OrderConfig.Maximum); err == nil {
				dbPointer = dbPointer.Where("nroficelandicquotes + nrofenglishquotes <= ?", nr)
			}
			if nr, err := strconv.Atoi(requestBody.OrderConfig.Minimum); err == nil {
				dbPointer = dbPointer.Where("nroficelandicquotes + nrofenglishquotes >= ?", nr)
			}
			dbPointer = dbPointer.Order("nroficelandicquotes + nrofenglishquotes " + orderDirection)
		}

	default:
		if requestBody.OrderConfig.Minimum != "" {
			dbPointer = dbPointer.Where("initcap(name) >= ?", strings.ToUpper(requestBody.OrderConfig.Minimum))
		}
		if requestBody.OrderConfig.Maximum != "" {
			dbPointer = dbPointer.Where("initcap(name) <= ?", strings.ToUpper(requestBody.OrderConfig.Maximum))
		}
		dbPointer = dbPointer.Order("initcap(name) " + orderDirection)
	}

	err := dbPointer.Limit(requestBody.PageSize).Order("id").
		Offset(requestBody.Page * requestBody.PageSize).
		Find(&authors).
		Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Authors:", len(authors))

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
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	var quotes []QuoteView
	err := db.Table("searchview").
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

// swagger:route POST /search SEARCH generalSearchByString
// Search for quotes / authors by a general string-search that searches both in the names of the authors and the quotes themselves
//
// responses:
//	200: multipleQuotesResponse

// SearchByString handles POST requests to search for quotes / authors by a search-string
func SearchByString(rw http.ResponseWriter, r *http.Request) {
	var requestBody Request
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	start := time.Now()

	var results []QuoteView
	m1 := regexp.MustCompile(` `)
	phrasesearch := m1.ReplaceAllString(requestBody.SearchString, " <-> ")
	generalsearch := m1.ReplaceAllString(requestBody.SearchString, " | ")

	//Order by authorid to have definitive order (when for examplke some quotes rank the same for plain, phrase, general and similarity)
	dbPointer := db.Table("searchview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq,to_tsquery(?) as generalq ",
		requestBody.SearchString, phrasesearch, generalsearch).
		Select("*, ts_rank(quotetsv, plainq) as plainrank, ts_rank(quotetsv, phraseq) as phraserank, ts_rank(quotetsv, generalq) as generalrank").
		Where("( tsv @@ plainq OR tsv @@ phraseq OR ? % ANY(STRING_TO_ARRAY(name,' ')) OR tsv @@ generalq)", requestBody.SearchString).
		Clauses(clause.OrderBy{
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
	var requestBody Request
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	start := time.Now()

	var results []QuoteView

	//Order by authorid to have definitive order (when for examplke some names rank the same for similarity), same for why quoteid
	//% is same as SIMILARITY but with default threshold 0.3
	dbPointer := db.Table("searchview").
		Where("( nametsv @@ plainto_tsquery(?) OR ? % ANY(STRING_TO_ARRAY(name,' ')) )", requestBody.SearchString, requestBody.SearchString).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "similarity(name, ?) DESC, authorid DESC, quoteid DESC", Vars: []interface{}{requestBody.SearchString}, WithoutParentheses: true},
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
	var requestBody Request
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	start := time.Now()

	var results []QuoteView
	m1 := regexp.MustCompile(` `)
	phrasesearch := m1.ReplaceAllString(requestBody.SearchString, " <-> ")
	generalsearch := m1.ReplaceAllString(requestBody.SearchString, " | ")

	//Order by quoteid to have definitive order (when for examplke some quotes rank the same for plain, phrase and general)
	dbPointer := db.Table("searchview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq,to_tsquery(?) as generalq ",
		requestBody.SearchString, phrasesearch, generalsearch).
		Select("*, ts_rank(quotetsv, plainq) as plainrank, ts_rank(quotetsv, phraseq) as phraserank, ts_rank(quotetsv, generalq) as generalrank").
		Where("( quotetsv @@ plainq OR quotetsv @@ phraseq OR quotetsv @@ generalq)").
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
	var requestBody Request
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	var results []ListItem

	pointer := db.Table("topics")

	switch strings.ToLower(requestBody.Language) {
	case "english":
		pointer = pointer.Not("isicelandic")
	case "icelandic":
		pointer = pointer.Where("isicelandic")
	}

	err := pointer.Find(&results).Error
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
	var requestBody Request
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	var results []QuoteView

	//Order by quoteid to have definitive order (when for examplke some quotes rank the same for plain, phrase and general)
	dbPoint := db.Table("topicsview")

	if requestBody.Topic != "" {
		dbPoint = dbPoint.Where("lower(topicname) = lower(?)", requestBody.Topic)
	} else {
		dbPoint = dbPoint.Where("topicid = ?", requestBody.Id)
	}

	err := dbPoint.Clauses(clause.OrderBy{
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

// swagger:route POST /quotes/random QUOTES getRandomQuote
// Get a random quote according to the given parameters
// responses:
//	200: randomQuoteResponse

// GetRandomQuote handles POST requests for getting a random quote
func GetRandomQuote(rw http.ResponseWriter, r *http.Request) {
	var requestBody Request
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	var dbPointer *gorm.DB
	var result []QuoteView
	shouldOrderBy := false //Used when there are few rows to choose from and therefore higher probability that random() < 0.005 returns no rows

	m1 := regexp.MustCompile(` `)
	phrasesearch := m1.ReplaceAllString(requestBody.SearchString, " <-> ")

	//Random quote from a particular topic
	if requestBody.TopicId > 0 {
		dbPointer = db.Table("topicsview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq", requestBody.SearchString, phrasesearch).Where("topicid = ?", requestBody.TopicId)
		shouldOrderBy = true
	} else {
		dbPointer = db.Table("searchview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq", requestBody.SearchString, phrasesearch)
	}

	//Random quote from a particular author
	if requestBody.AuthorId > 0 {
		dbPointer = dbPointer.Where("authorid = ?", requestBody.AuthorId)
		shouldOrderBy = true
	}

	//Random quote from a particular language
	if requestBody.Language != "" {
		switch strings.ToLower(requestBody.Language) {
		case "english":
			dbPointer = dbPointer.Not("isicelandic")
		case "icelandic":
			dbPointer = dbPointer.Where("isicelandic")
		}
	}

	if requestBody.SearchString != "" {
		dbPointer = dbPointer.Where("( quotetsv @@ plainq OR quotetsv @@ phraseq)")
		shouldOrderBy = true
	}

	//Order by used to get random quote if there are "few" rows returned
	if shouldOrderBy {
		dbPointer = dbPointer.Order("random()") //Randomized, O( n*log(n) )
	} else {
		dbPointer = dbPointer.
			Where("random() < 0.005") //Randomized, O(n)
	}

	err := dbPointer.Find(&result).Error
	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if len(result) == 0 {
		json.NewEncoder(rw).Encode(result)
	} else {
		json.NewEncoder(rw).Encode(result[0])
	}

}

// swagger:route POST /authors/random AUTHORS getRandomAuthor
// Get a random Author, and some of his quotes, according to the given parameters
// responses:
//	200: randomAuthorResponse

// GetRandomAuthor handles POST requests for getting a random author
func GetRandomAuthor(rw http.ResponseWriter, r *http.Request) {
	var requestBody Request
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var result []QuoteView
	var author AuthorsView
	dbPointer := db.Table("authors").Where("random() < 0.01")

	//Random author from a particular language
	if requestBody.Language != "" {
		switch strings.ToLower(requestBody.Language) {
		case "english":
			dbPointer = dbPointer.Not("hasicelandicquotes")
		case "icelandic":
			dbPointer = dbPointer.Where("hasicelandicquotes")
		}
	}

	err := dbPointer.First(&author).Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	dbPointer = db.Table("searchview").Where("authorid = ?", author.Id)

	//An icelandic quote from the particular/random author
	if requestBody.Language != "" {
		switch strings.ToLower(requestBody.Language) {
		case "english":
			dbPointer = dbPointer.Not("isicelandic")
		case "icelandic":
			dbPointer = dbPointer.Where("isicelandic")
		}
	}

	err = dbPointer.Limit(requestBody.MaxQuotes).Find(&result).Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(result)

}

// swagger:route GET /languages META getLanguages
// Get languages supported by the api
// responses:
//	200: listOfStrings

// ListLanguages handles GET requests for getting the languages supported by the api
func ListLanguagesSupported(rw http.ResponseWriter, r *http.Request) {

	type response = struct {
		Languages []string `json:"languages"`
	}

	json.NewEncoder(rw).Encode(&response{
		Languages: languages,
	})
}
