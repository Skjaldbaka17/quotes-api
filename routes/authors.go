package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Skjaldbaka17/quotes-api/handlers"
	"github.com/Skjaldbaka17/quotes-api/structs"

	"gorm.io/gorm"
)

// swagger:route POST /authors AUTHORS GetAuthors
// Get the authors by their ids
//
// responses:
//	200: authorsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// Get Authors handles POST requests to get the authors, and their quotes, that have the given ids
func GetAuthorsById(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var authors []structs.AuthorsView
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	err := handlers.Db.Table("authors").
		Where("id in (?)", requestBody.Ids).
		Scan(&authors).
		Error
	//** ---------- Paramatere configuratino for DB query ends ---------- **//

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when querying DB in GetAuthorsById: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	//Update popularity in background!
	go handlers.DirectFetchAuthorsCountIncrement(requestBody.Ids)

	json.NewEncoder(rw).Encode(&authors)
}

// swagger:route POST /authors/list AUTHORS ListAuthors
//
// Get a list of authors according to some ordering / parameters
//
// responses:
//	200: authorsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// GetAuthorsList handles POST requests to get the authors that fit the parameters
func GetAuthorsList(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var authors []structs.AuthorsView
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := handlers.Db.Table("authors")

	dbPointer = authorLanguageSQL(requestBody.Language, dbPointer)

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
		dbPointer = dbPointer.Order("count " + orderDirection)
	case "nrofquotes":
		switch strings.ToLower(requestBody.Language) {
		case "english":
			dbPointer = setMaxMinNumber(requestBody.OrderConfig, "nrofenglishquotes", orderDirection, dbPointer)
		case "icelandic":
			dbPointer = setMaxMinNumber(requestBody.OrderConfig, "nroficelandicquotes", orderDirection, dbPointer)
		default:
			dbPointer = setMaxMinNumber(requestBody.OrderConfig, "nroficelandicquotes + nrofenglishquotes", orderDirection, dbPointer)
		}

	default:
		//Minimum letter to start with (i.e. start from given minimum letter of the alphabet)
		if requestBody.OrderConfig.Minimum != "" {
			dbPointer = dbPointer.Where("initcap(name) >= ?", strings.ToUpper(requestBody.OrderConfig.Minimum))
		}
		//Maximum letter to start with (i.e. end at the given maximum letter of the alphabet)
		if requestBody.OrderConfig.Maximum != "" {
			dbPointer = dbPointer.Where("initcap(name) <= ?", strings.ToUpper(requestBody.OrderConfig.Maximum))
		}
		dbPointer = dbPointer.Order("initcap(name) " + orderDirection)
	}

	//** ---------- Paramatere configuratino for DB query ends---------- **//
	err := pagination(requestBody, dbPointer).Order("id").
		Find(&authors).
		Error

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when querying DB in GetAuthors: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	//Update popularity in background!
	go handlers.AuthorsAppearInSearchCountIncrement(authors)

	json.NewEncoder(rw).Encode(&authors)
}

// swagger:route POST /authors/random AUTHORS GetRandomAuthor
// Get a random Author, and some of his quotes, according to the given parameters
// responses:
//	200: randomAuthorResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// GetRandomAuthor handles POST requests for getting a random author
func GetRandomAuthor(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var result []structs.QuoteView
	var author structs.AuthorsView
	//** ---------- Paramatere configuratino for DB query begins ---------- **//

	//Get Random author
	dbPointer := handlers.Db.Table("authors").Order("random()")

	//author from a particular language
	dbPointer = authorLanguageSQL(requestBody.Language, dbPointer)
	//** ---------- Paramatere configuratino for DB query ends ---------- **//

	err := dbPointer.First(&author).Error

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when querying DB, first one, in GetRandomAuthor: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	dbPointer = handlers.Db.Table("searchview").Where("authorid = ?", author.Id)

	//An icelandic quote from the particular/random author
	dbPointer = quoteLanguageSQL(requestBody.Language, dbPointer)

	err = dbPointer.Limit(requestBody.MaxQuotes).Find(&result).Error

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when querying DB, second one, in GetAuthors: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	json.NewEncoder(rw).Encode(result)
}

// swagger:route POST /authors/aod AUTHORS GetAuthorOfTheDay
// Gets the author of the day
// responses:
//	200: authorOfTheDayResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

//GetAuthorOfTheDay gets the author of the day
func GetAuthorOfTheDay(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	if requestBody.Language == "" {
		requestBody.Language = "English"
	}

	var author structs.AuthorsView
	var err error

	//** ---------- Paramatere configuratino for DB query begins ---------- **//

	//Which table to look for quotes (ice table has icelandic quotes)
	dbPointer := aodLanguageSQL(requestBody.Language).
		Where("date = current_date")
	//** ---------- Paramatere configuratino for DB query ends ---------- **//

	err = dbPointer.Scan(&author).Error

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when querying DB in GetAuthorOfTheDay: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	if (structs.AuthorsView{}) == author {
		err = setNewRandomAOD(requestBody.Language)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Printf("Got error when setting new random AOD in GetAuthorOfTheDay: %s", err)
			json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
			return
		}

		GetAuthorOfTheDay(rw, r) //Dangerous? possibility of endless cycle? Only iff the setNewRandomAOD fails in some way. Or the date is not saved correctly into the DB?
		return
	}

	json.NewEncoder(rw).Encode(author)
}

// swagger:route POST /authors/aod/history AUTHORS GetAODHistory
// Gets the history for the authors of the day
// responses:
//	200: aodHistoryResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

//GetAODHistory gets Aod history starting from some point
func GetAODHistory(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	if requestBody.Language == "" {
		requestBody.Language = "English"
	}
	var authors []structs.AuthorsView
	var err error
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := aodLanguageSQL(requestBody.Language)

	if requestBody.Minimum == "" {
		requestBody.Minimum = "1900-12-21"
	}
	now := time.Now()
	minDate, err := time.Parse("2006-01-02", requestBody.Minimum)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Printf("Got error when parsing mindate in GetAODHistory: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "Please supply date in '2020-12-21' format"})
		return
	}

	if !now.After(minDate) {
		rw.WriteHeader(http.StatusBadRequest)
		log.Printf("Got error when comparing mindate to today in GetAodHistory: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "Please send a minimum date that is before today"})
		return
	}

	//Not maximum because then possibility of endless cycle with the if statement below!
	if requestBody.Minimum != "" {
		dbPointer = dbPointer.Where("date >= ?", requestBody.Minimum)
	}
	dbPointer = dbPointer.Where("date <= current_date").Order("date DESC")
	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err = dbPointer.Find(&authors).Error

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when querying DB in GetAodHistory: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	reg := regexp.MustCompile(time.Now().Format("2006-01-02"))

	if len(authors) == 0 || !reg.Match([]byte(authors[0].Date)) {
		err = setNewRandomAOD(requestBody.Language)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Printf("Got error when setting newRandomAOD in getAODHistory: %s", err)
			json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
			return
		}

		GetAODHistory(rw, r)
		return
	}

	json.NewEncoder(rw).Encode(authors)
}

// swagger:route POST /authors/aod/new AUTHORS SetAuthorOfTheDay
//
// sets the author of the day for the given dates
//
// responses:
//	200: successResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

//SetAuthorOfTheDay sets the author of the day.
func SetAuthorOfTheDay(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	if len(requestBody.Aods) == 0 {
		log.Println("No author supplied")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "Please supply some authors", StatusCode: http.StatusBadRequest})
		return
	}

	for _, aod := range requestBody.Aods {
		err := setAOD(requestBody.Language, aod.Date, aod.Id)
		if err != nil {
			log.Printf("Got err when setting the AOD for %+v: %s", aod, err)
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "Some of the authors (ids) you supplied do not have " + requestBody.Language + " quotes", StatusCode: http.StatusBadRequest})
			return
		}
	}

	json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "Successfully inserted quote of the day!", StatusCode: http.StatusOK})
}

//setAOD inserts a new row into the aod/aodice table
func setAOD(language string, date string, authorId int) error {
	switch strings.ToLower(language) {
	case "icelandic":
		return handlers.Db.Exec("insert into aodice (authorid, date) values((select id from authors where id = ? and hasicelandicquotes), ?) on conflict (date) do update set authorid = ?", authorId, date, authorId).Error
	default:
		return handlers.Db.Exec("insert into aod (authorid, date) values((select id from authors where id = ? and not hasicelandicquotes), ?) on conflict (date) do update set authorid = ?", authorId, date, authorId).Error
	}
}

//SetNewRandomQOD sets a random quote as the qod for today (if language=icelandic is supplied then it adds the random qod to the icelandic qod table)
func setNewRandomAOD(language string) error {
	var authorItem structs.ListItem

	if language == "" {
		language = "english"
	}
	dbPointer := handlers.Db.Table("authors")
	dbPointer = authorLanguageSQL(language, dbPointer)

	log.Println("HEREMatur")
	err := dbPointer.Order("random()").Limit(1).Scan(&authorItem).Error
	if err != nil {
		return err
	}

	log.Println("MASSI:", authorItem)

	return setAOD(language, time.Now().Format("2006-01-02"), authorItem.Id)
}

//aodLanguageSQL adds to the sql query for the authors db a condition of whether the authors to be fetched have quotes in a particular language
func aodLanguageSQL(language string) *gorm.DB {
	switch strings.ToLower(language) {
	case "icelandic":
		return handlers.Db.Table("aodiceview")
	default:
		return handlers.Db.Table("aodview")
	}
}

//qodLanguageSQL adds to the sql query for the quotes db a condition of whether the quotes to be fetched are quotes in a particular language
func qodLanguageSQL(language string) *gorm.DB {
	switch strings.ToLower(language) {
	case "icelandic":
		return handlers.Db.Table("qodiceview")
	default:
		return handlers.Db.Table("qodview")
	}
}

//authorLanguageSQL adds to the sql query for the authors db a condition of whether the authors to be fetched have quotes in a particular language
func authorLanguageSQL(language string, dbPointer *gorm.DB) *gorm.DB {
	if language != "" {
		switch strings.ToLower(language) {
		case "english":
			dbPointer = dbPointer.Not("hasicelandicquotes")
		case "icelandic":
			dbPointer = dbPointer.Where("hasicelandicquotes")
		}
	}
	return dbPointer
}

//quoteLanguageSQL adds to the sql query for the quotes db a condition of whether the quotes to be fetched are in a particular language
func quoteLanguageSQL(language string, dbPointer *gorm.DB) *gorm.DB {
	if language != "" {
		switch strings.ToLower(language) {
		case "english":
			dbPointer = dbPointer.Not("isicelandic")
		case "icelandic":
			dbPointer = dbPointer.Where("isicelandic")
		}
	}
	return dbPointer
}

//setMaxMinNumber sets the condition for which authors to return
func setMaxMinNumber(orderConfig structs.OrderConfig, column string, orderDirection string, dbPointer *gorm.DB) *gorm.DB {
	if nr, err := strconv.Atoi(orderConfig.Maximum); err == nil {
		dbPointer = dbPointer.Where(column+" <= ?", nr)
	}
	if nr, err := strconv.Atoi(orderConfig.Minimum); err == nil {
		dbPointer = dbPointer.Where(column+" >= ?", nr)
	}
	return dbPointer.Order(column + " " + orderDirection)
}
