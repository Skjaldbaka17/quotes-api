package routes

import (
	"encoding/json"
	"fmt"
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

// swagger:route POST /authors AUTHORS getAuthorsByIds
//
// Get authors by their ids
//
// responses:
//	200: authorsResponse

// Get Authors handles POST requests to get the authors, and their quotes, that have the given ids
func GetAuthorsById(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var authors []structs.AuthorsView
	err := handlers.Db.Table("authors").
		Where("id in (?)", requestBody.Ids).
		Scan(&authors).
		Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//Update popularity in background!
	go handlers.DirectFetchAuthorsCountIncrement(requestBody.Ids)

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
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var authors []structs.AuthorsView
	dbPointer := handlers.Db.Table("authors")

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
		orderDirection = "DESC"
		if requestBody.OrderConfig.Reverse {
			orderDirection = "ASC"
		}
		dbPointer = dbPointer.Order("count " + orderDirection)
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

	//Update popularity in background!
	go handlers.AuthorsAppearInSearchCountIncrement(authors)

	json.NewEncoder(rw).Encode(&authors)
}

// swagger:route POST /authors/random AUTHORS getRandomAuthor
// Get a random Author, and some of his quotes, according to the given parameters
// responses:
//	200: randomAuthorResponse

// GetRandomAuthor handles POST requests for getting a random author
func GetRandomAuthor(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var result []structs.QuoteView
	var author structs.AuthorsView
	dbPointer := handlers.Db.Table("authors").Where("random() < 0.01")

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

	dbPointer = handlers.Db.Table("searchview").Where("authorid = ?", author.Id)

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

// swagger:route POST /quotes/qod AUTHORS getAuthorOfTheDay
// Gets the author of the day
// responses:
//	200: randomQuoteResponse

//GetAuthorOfTheDay gets the author of the day
func GetAuthorOfTheDay(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	if requestBody.Language == "" {
		requestBody.Language = "English"
	}

	var author structs.QuoteView
	var dbPointer *gorm.DB
	var err error
	switch strings.ToLower(requestBody.Language) {
	case "icelandic":
		dbPointer = handlers.Db.Table("aodiceview")
	default:
		dbPointer = handlers.Db.Table("aodview")
	}

	err = dbPointer.Where("date = current_date").Scan(&author).Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if (structs.QuoteView{}) == author {
		fmt.Println("Setting a brand new AOD for today")
		err = setNewRandomAOD(requestBody.Language)
		if err != nil {
			//TODO: Respond with better error -- and put into swagger -- and add tests
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		switch strings.ToLower(requestBody.Language) {
		case "icelandic":
			err = handlers.Db.Table("qodiceview").Where("date = current_date").Scan(&author).Error
		default:
			err = handlers.Db.Table("qodview").Where("date = current_date").Scan(&author).Error
		}
		log.Println(author)

		if err != nil {
			//TODO: Respond with better error -- and put into swagger -- and add tests
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}

	json.NewEncoder(rw).Encode(author)
}

// swagger:route POST /quotes/qod AUTHORS getAODHistory
// Gets the history for the authors of the day
// responses:
//	200: qodHistoryResponse

//GetAODHistory gets Aod history starting from some point
func GetAODHistory(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	if requestBody.Language == "" {
		requestBody.Language = "English"
	}
	var authors []structs.QuoteView
	var dbPointer *gorm.DB
	var err error
	switch strings.ToLower(requestBody.Language) {
	case "icelandic":
		dbPointer = handlers.Db.Table("aodiceview")
	default:
		dbPointer = handlers.Db.Table("aodview")
	}

	if requestBody.Minimum != "" {
		dbPointer = dbPointer.Where("date >= ?", requestBody.Minimum)
	}

	if requestBody.Maximum != "" {
		dbPointer = dbPointer.Where("date <= ?", requestBody.Maximum)
	}

	dbPointer = dbPointer.Where("date <= current_date")

	err = dbPointer.Order("date DESC").Find(&authors).Error

	log.Println("HERE:", authors)
	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	reg := regexp.MustCompile(time.Now().Format("2006-01-02"))
	if len(authors) == 0 || !reg.Match([]byte(authors[0].Date)) {
		log.Println("Setting a brand new AOD for today")
		err = setNewRandomAOD(requestBody.Language)
		if err != nil {
			log.Println(err)
			//TODO: Respond with better error -- and put into swagger -- and add tests
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		switch strings.ToLower(requestBody.Language) {
		case "icelandic":
			dbPointer = handlers.Db.Table("aodiceview")
		default:
			dbPointer = handlers.Db.Table("aodview")
		}

		if requestBody.Minimum != "" {
			dbPointer = dbPointer.Where("date >= ?", requestBody.Minimum)
		}

		if requestBody.Maximum != "" {
			dbPointer = dbPointer.Where("date <= ?", requestBody.Maximum)
		}

		dbPointer = dbPointer.Where("date <= current_date")

		err = dbPointer.Order("date").Find(&authors).Error

		if err != nil {
			//TODO: Respond with better error -- and put into swagger -- and add tests
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

	}

	json.NewEncoder(rw).Encode(authors)
}

// swagger:route POST /quotes/aod/new AUTHORS setAuthorsOfTheDay
// Sets the author of the day for the given date. It Is password protected TODO: Put in privacy swagger
// responses:
//	200: successResponse

//SetAuthorOfTheDay sets the author of the day (is password protected)
func SetAuthorOfTheDay(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	if len(requestBody.Aods) == 0 {
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "Please supply some authors", StatusCode: http.StatusBadRequest})
		return
	}

	for _, aod := range requestBody.Aods {
		err := setAOD(requestBody.Language, aod.Date, aod.Id)
		if err != nil {
			log.Println(err)
			json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "Some of the authors (ids) you supplied do not have " + requestBody.Language + " quotes", StatusCode: http.StatusBadRequest})
			return
		}
	}

	json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "Successfully inserted quote of the day!", StatusCode: http.StatusOK})
}

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
	var dbPointer *gorm.DB
	switch strings.ToLower(language) {
	case "icelandic":
		dbPointer = handlers.Db.Table("authors").Where("hasicelandicquotes")
	default:
		dbPointer = handlers.Db.Table("authors").Not("hasicelandicquotes")
	}

	log.Println("HEREMatur")
	err := dbPointer.Order("random()").Limit(1).Scan(&authorItem).Error
	if err != nil {
		return err
	}

	log.Println("MASSI:", authorItem)

	return setAOD(language, time.Now().Format("2006-01-02"), authorItem.Id)
}
