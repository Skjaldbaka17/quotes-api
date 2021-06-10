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

// swagger:route POST /quotes QUOTES getQuotes
// Get quotes by their ids
//
// responses:
//	200: multipleQuotesResponse

//
//Params: in Body {ids:[]int, authorId: int}
//
//

// GetQuotes handles POST requests to get the quotes, and their authors, that have the given ids
func GetQuotes(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	var quotes []structs.QuoteView
	//** ---------- Paramatere configuratino for DB query begins ---------- **//

	dbPointer := handlers.Db.Table("searchview").Order("quoteid ASC")
	if requestBody.AuthorId > 0 {
		dbPointer = dbPointer.
			Where("authorid = ?", requestBody.AuthorId)
		dbPointer = pagination(requestBody, dbPointer)
	} else {
		dbPointer = dbPointer.Where("quoteid in ?", requestBody.Ids)
	}
	//** ---------- Paramatere configuratino for DB query ends ---------- **//

	err := dbPointer.Find(&quotes).Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//Update popularity in background!
	go handlers.DirectFetchQuotesCountIncrement(requestBody.Ids)

	json.NewEncoder(rw).Encode(&quotes)
}

// swagger:route POST /quotes/list QUOTES getQuotesList
//
// Get list of quotes according to some ordering / parameters
//
// responses:
//	200: multipleQuoteResponse

// GetQuotesList handles POST requests to get the quotes that fit the parameters

func GetQuotesList(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var quotes []structs.QuoteView
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
		dbPointer = dbPointer.Order("quotecount " + orderDirection)
	case "length":
		dbPointer = setMaxMinNumber(requestBody.OrderConfig, "length(quote)", orderDirection, dbPointer)
	default:
		dbPointer = setMaxMinNumber(requestBody.OrderConfig, "quoteid", orderDirection, dbPointer)
	}

	//** ---------- Paramatere configuratino for DB query ends ---------- **//

	err := pagination(requestBody, dbPointer).Order("quoteid").
		Find(&quotes).
		Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//Update popularity in background!
	go handlers.QuotesAppearInSearchCountIncrement(quotes)

	json.NewEncoder(rw).Encode(&quotes)
}

// swagger:route POST /quotes/random QUOTES getRandomQuote
// Get a random quote according to the given parameters
// responses:
//	200: randomQuoteResponse

// GetRandomQuote handles POST requests for getting a random quote
func GetRandomQuote(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	var dbPointer *gorm.DB
	var result structs.QuoteView
	shouldOrderBy := false //Used when there are few rows to choose from and therefore higher probability that random() < 0.005 returns no rows

	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	m1 := regexp.MustCompile(` `)
	phrasesearch := m1.ReplaceAllString(requestBody.SearchString, " <-> ")

	//Random quote from a particular topic
	if requestBody.TopicId > 0 {
		dbPointer = handlers.Db.Table("topicsview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq", requestBody.SearchString, phrasesearch).Where("topicid = ?", requestBody.TopicId)
		shouldOrderBy = true
	} else {
		dbPointer = handlers.Db.Table("searchview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq", requestBody.SearchString, phrasesearch)
	}

	//Random quote from a particular author
	if requestBody.AuthorId > 0 {
		dbPointer = dbPointer.Where("authorid = ?", requestBody.AuthorId)
		shouldOrderBy = true
	}

	//Random quote from a particular language
	dbPointer = quoteLanguageSQL(requestBody.Language, dbPointer)

	if requestBody.SearchString != "" {
		dbPointer = dbPointer.Where("( quotetsv @@ plainq OR quotetsv @@ phraseq)")
		shouldOrderBy = true
	}

	//Order by used to get random quote if there are "few" rows returned
	if !shouldOrderBy {
		dbPointer = dbPointer.
			Where("random() < 0.005")
	}

	dbPointer = dbPointer.Order("random()") //Randomized, O( n*log(n) )
	//** ---------- Paramater configuratino for DB query ends ---------- **//

	err := dbPointer.Limit(1).Find(&result).Error
	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(result)

}

// swagger:route POST /quotes/qod/new QUOTES setQuoteOfTheDay
// Sets the quote of the day for the given date. It Is password protected TODO: Put in privacy swagger
// responses:
//	200: successResponse

//SetQuoteOfTheyDay sets the quote of the day (is password protected)
func SetQuoteOfTheDay(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	if requestBody.Language == "" {
		requestBody.Language = "English"
	}

	if len(requestBody.Qods) == 0 {
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "Please supply some quotes", StatusCode: http.StatusBadRequest})
		return
	}

	for _, qod := range requestBody.Qods {
		err := setQOD(requestBody.Language, qod.Date, qod.Id)
		if err != nil {
			log.Println(err)
			json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "Some of the quotes (ids) you supplied are not in " + requestBody.Language, StatusCode: http.StatusBadRequest})
			return
		}
	}

	json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "Successfully inserted quote of the day!", StatusCode: http.StatusOK})
}

//setQOD inserts a new row into qod/qodice table
func setQOD(language string, date string, quoteId int) error {
	switch strings.ToLower(language) {
	case "icelandic":
		return handlers.Db.Exec("insert into qodice (quoteid, date) values((select id from quotes where id = ? and isicelandic), ?) on conflict (date) do update set quoteid = ?", quoteId, date, quoteId).Error
	default:
		return handlers.Db.Exec("insert into qod (quoteid, date) values((select id from quotes where id = ? and not isicelandic), ?) on conflict (date) do update set quoteid = ?", quoteId, date, quoteId).Error
	}
}

//SetNewRandomQOD sets a random quote as the qod for today (if language=icelandic is supplied then it adds the random qod to the icelandic qod table)
func setNewRandomQOD(language string) error {
	var quoteItem structs.ListItem
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

// swagger:route POST /quotes/qod QUOTES getQuoteOfTheDay
// Gets the quote of the day
// responses:
//	200: randomQuoteResponse

//GetQuoteOfTheyDay gets the quote of the day
func GetQuoteOfTheDay(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	if requestBody.Language == "" {
		requestBody.Language = "English"
	}

	var quote structs.QuoteView
	var err error
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := qodLanguageSQL(requestBody.Language).Where("date = current_date")
	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err = dbPointer.Scan(&quote).Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if (structs.QuoteView{}) == quote {
		fmt.Println("Setting a brand new QOD for today")
		err = setNewRandomQOD(requestBody.Language)
		if err != nil {
			//TODO: Respond with better error -- and put into swagger -- and add tests
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		GetQuoteOfTheDay(rw, r)
		return
	}

	json.NewEncoder(rw).Encode(quote)
}

// swagger:route POST /quotes/qod QUOTES getQODHistory
// Gets the history for the quotes of the day
// responses:
//	200: qodHistoryResponse

//GetQODHistory gets Qod history starting from some point
func GetQODHistory(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	if requestBody.Language == "" {
		requestBody.Language = "English"
	}

	var quotes []structs.QuoteView
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
		//TODO: Respond with better error -- and put into swagger -- and add tests
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if len(quotes) == 0 {
		err = setNewRandomQOD(requestBody.Language)
		if err != nil {
			//TODO: Respond with better error -- and put into swagger -- and add tests
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		GetQODHistory(rw, r)
		return
	}

	json.NewEncoder(rw).Encode(quotes)
}
