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
)

// swagger:route POST /quotes QUOTES GetQuotes
// Get quotes by their ids
//
// responses:
//	200: quotesResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// GetQuotes handles POST requests to get the quotes, and their authors, that have the given ids
func GetQuotes(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	var quotes []structs.SearchViewDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//

	dbPointer := handlers.Db.Table("searchview").Order("quote_id ASC")
	if requestBody.AuthorId > 0 {
		dbPointer = dbPointer.
			Where("author_id = ?", requestBody.AuthorId)
		dbPointer = pagination(requestBody, dbPointer)
	} else {
		dbPointer = dbPointer.Where("quote_id in ?", requestBody.Ids)
	}
	//** ---------- Paramatere configuratino for DB query ends ---------- **//

	err := dbPointer.Find(&quotes).Error

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when querying DB in GetQuotes: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	//Update popularity in background!
	go handlers.DirectFetchQuotesCountIncrement(requestBody.Ids)

	searchViewsAPI := structs.ConvertToSearchViewsAPIModel(quotes)
	json.NewEncoder(rw).Encode(searchViewsAPI)
}

// swagger:route POST /quotes/list QUOTES GetQuotesList
//
// Get list of quotes according to some ordering / parameters
//
// responses:
//	200: quotesResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// GetQuotesList handles POST requests to get the quotes that fit the parameters

func GetQuotesList(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var quotes []structs.SearchViewDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := handlers.Db.Table("searchview")
	dbPointer = quoteLanguageSQL(requestBody.Language, dbPointer)

	orderDirection := "ASC"
	if requestBody.OrderConfig.Reverse {
		orderDirection = "DESC"
	}

	switch strings.ToLower(requestBody.OrderConfig.OrderBy) {
	case "popularity": //TODO: add popularity ordering
		orderDirection = "DESC"
		if requestBody.OrderConfig.Reverse {
			orderDirection = "ASC"
		}
		dbPointer = dbPointer.Order("quote_count " + orderDirection)
	case "length":
		dbPointer = setMaxMinNumber(requestBody.OrderConfig, "length(quote)", orderDirection, dbPointer)
	default:
		dbPointer = setMaxMinNumber(requestBody.OrderConfig, "quote_id", orderDirection, dbPointer)
	}

	//** ---------- Paramatere configuratino for DB query ends ---------- **//

	err := pagination(requestBody, dbPointer).Order("quote_id").
		Find(&quotes).
		Error

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when querying DB in GetQuotesList: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	//Update popularity in background!
	go handlers.QuotesAppearInSearchCountIncrement(quotes)
	searchViewsAPI := structs.ConvertToSearchViewsAPIModel(quotes)
	json.NewEncoder(rw).Encode(searchViewsAPI)
}

// swagger:route POST /quotes/random QUOTES GetRandomQuote
// Get a random quote according to the given parameters
// responses:
//	200: randomQuoteResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// GetRandomQuote handles POST requests for getting a random quote
func GetRandomQuote(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	var dbPointer *gorm.DB
	var topicResult structs.TopicViewDBModel
	var searchViewResult structs.SearchViewDBModel
	// shouldOrderBy := false //Used when there are few rows to choose from and therefore higher probability that random() < 0.005 returns no rows

	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	m1 := regexp.MustCompile(` `)
	phrasesearch := m1.ReplaceAllString(requestBody.SearchString, " <-> ")

	//Random quote from a particular topic
	if requestBody.TopicId > 0 {
		dbPointer = handlers.Db.Table("topicsview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq", requestBody.SearchString, phrasesearch).Where("topic_id = ?", requestBody.TopicId)
		// shouldOrderBy = true
	} else {
		dbPointer = handlers.Db.Table("searchview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq", requestBody.SearchString, phrasesearch)
	}

	//Random quote from a particular author
	if requestBody.AuthorId > 0 {
		dbPointer = dbPointer.Where("author_id = ?", requestBody.AuthorId)
		// shouldOrderBy = true
	}

	//Random quote from a particular language
	dbPointer = quoteLanguageSQL(requestBody.Language, dbPointer)

	if requestBody.SearchString != "" {
		dbPointer = dbPointer.Where("( quote_tsv @@ plainq OR quote_tsv @@ phraseq)")
		// shouldOrderBy = true
	}

	//Order by used to get random quote if there are "few" rows returned
	// if !shouldOrderBy {
	// 	dbPointer = dbPointer.
	// 		Where("random() < 0.005")
	// }

	dbPointer = dbPointer.Order("random()") //Randomized, O( n*log(n) )
	//** ---------- Paramater configuratino for DB query ends ---------- **//
	if requestBody.TopicId > 0 {
		err := dbPointer.Limit(1).Find(&topicResult).Error
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Printf("Got error when querying DB in GetRandomQuote: %s", err)
			json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
			return
		}
		json.NewEncoder(rw).Encode(topicResult.ConvertToAPIModel())
	} else {
		err := dbPointer.Limit(1).Find(&searchViewResult).Error
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Printf("Got error when querying DB in GetRandomQuote: %s", err)
			json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
			return
		}
		json.NewEncoder(rw).Encode(searchViewResult.ConvertToAPIModel())
	}

}

// swagger:route POST /quotes/qod/new QUOTES SetQuoteOfTheDay
// Sets the quote of the day for the given dates
// responses:
//	200: successResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

//SetQuoteOfTheyDay sets the quote of the day (is password protected)
func SetQuoteOfTheDay(rw http.ResponseWriter, r *http.Request) {
	if err := handlers.AuthorizeGODApiKey(rw, r); err != nil {
		return
	}

	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	if requestBody.Language == "" {
		requestBody.Language = "English"
	}

	if len(requestBody.Qods) == 0 {
		log.Println("Not QODS supplied when setting quote of the day")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "Please supply some quotes", StatusCode: http.StatusBadRequest})
		return
	}

	for _, qod := range requestBody.Qods {
		err := setQOD(requestBody.Language, qod.Date, qod.Id)
		if err != nil {
			log.Printf("Got error when settin the qod %+v as QOD: %s", qod, err)
			json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "Some of the quotes (ids) you supplied are not in " + requestBody.Language, StatusCode: http.StatusBadRequest})
			return
		}

	}

	json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "Successfully inserted quote of the day!", StatusCode: http.StatusOK})
}

// swagger:route POST /quotes/qod QUOTES GetQuoteOfTheDay
// Gets the quote of the day
// responses:
//	200: quoteOfTheDayResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

//GetQuoteOfTheyDay gets the quote of the day
func GetQuoteOfTheDay(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	if requestBody.Language == "" {
		requestBody.Language = "English"
	}

	var quote structs.QodViewDBModel
	var err error
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := qodLanguageSQL(requestBody.Language).Where("date = current_date")
	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err = dbPointer.Scan(&quote).Error

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when querying DB in GetQODs: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	if (structs.QodViewDBModel{}) == quote {
		fmt.Println("Setting a brand new QOD for today")
		err = setNewRandomQOD(requestBody.Language)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Printf("Got error when setting new random qod: %s", err)
			json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
			return
		}

		GetQuoteOfTheDay(rw, r)
		return
	}

	json.NewEncoder(rw).Encode(quote.ConvertToAPIModel())
}

// swagger:route POST /quotes/qod QUOTES GetQODHistory
// Gets the history for the quotes of the day
// responses:
//	200: qodHistoryResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

//GetQODHistory gets Qod history starting from some point
func GetQODHistory(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	if requestBody.Language == "" {
		requestBody.Language = "English"
	}

	var quotes []structs.QodViewDBModel
	var err error
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := qodLanguageSQL(requestBody.Language)

	//Not maximum because then possibility of endless cycle with the if statement below!
	if requestBody.Minimum != "" {
		dbPointer = dbPointer.Where("date >= ?", requestBody.Minimum)
	}
	dbPointer = dbPointer.Where("date <= current_date").Order("date DESC")
	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err = dbPointer.Find(&quotes).Error

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when querying DB in GetQODHistory: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	if len(quotes) == 0 {
		err = setNewRandomQOD(requestBody.Language)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Printf("Got error when querying setting new Random QOD in history: %s", err)
			json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
			return
		}
		GetQODHistory(rw, r)
		return
	}

	qodHistoryAPI := structs.ConvertToQodViewsAPIModel(quotes)
	json.NewEncoder(rw).Encode(qodHistoryAPI)
}

//setQOD inserts a new row into qod/qodice table
func setQOD(language string, date string, quoteId int) error {
	switch strings.ToLower(language) {
	case "icelandic":
		return handlers.Db.Exec("insert into qodice (quote_id, date) values((select id from quotes where id = ? and is_icelandic), ?) on conflict (date) do update set quote_id = ?", quoteId, date, quoteId).Error
	default:
		return handlers.Db.Exec("insert into qod (quote_id, date) values((select id from quotes where id = ? and not is_icelandic), ?) on conflict (date) do update set quote_id = ?", quoteId, date, quoteId).Error
	}
}

//SetNewRandomQOD sets a random quote as the qod for today (if language=icelandic is supplied then it adds the random qod to the icelandic qod table)
func setNewRandomQOD(language string) error {
	var quoteItem structs.QuoteDBModel
	var dbPointer *gorm.DB
	dbPointer = handlers.Db.Table("quotes")
	dbPointer = quoteLanguageSQL(language, dbPointer)
	if strings.ToLower(language) != "icelandic" {
		dbPointer = dbPointer.Where("Random() < 0.005")
	}

	err := dbPointer.Order("random()").Limit(1).Scan(&quoteItem).Error
	if err != nil {
		return err
	}

	return setQOD(language, time.Now().Format("2006-01-02"), quoteItem.Id)
}
