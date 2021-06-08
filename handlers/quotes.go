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

	const layout = "2006-01-02"
	//Set date into correct format, if supplied, otherwise input today's date in the correct format for all qods
	if len(requestBody.Qods) != 0 {
		for idx, _ := range requestBody.Qods {
			if requestBody.Qods[idx].Date == "" {
				requestBody.Qods[idx].Date = time.Now().UTC().Format(layout)
			} else {
				var parsedDate time.Time
				parsedDate, err = time.Parse(layout, requestBody.Qods[idx].Date)
				if err != nil {
					log.Printf("Got error when decoding: %s", err)
					err = fmt.Errorf("the date is not structured correctly, should be in %s format", layout)
					rw.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(rw).Encode(ErrorResponse{Message: err.Error()})
					return err
				}

				requestBody.Qods[idx].Date = parsedDate.UTC().Format(layout)
			}
		}
	}

	//Set date into correct format, if supplied, otherwise input today's date in the correct format for all qods
	if len(requestBody.Aods) != 0 {
		for idx, _ := range requestBody.Aods {
			if requestBody.Aods[idx].Date == "" {
				requestBody.Aods[idx].Date = time.Now().UTC().Format(layout)
			} else {
				var parsedDate time.Time
				parsedDate, err = time.Parse(layout, requestBody.Aods[idx].Date)
				if err != nil {
					log.Printf("Got error when decoding: %s", err)
					err = fmt.Errorf("the date is not structured correctly, should be in %s format", layout)
					rw.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(rw).Encode(ErrorResponse{Message: err.Error()})
					return err
				}

				requestBody.Aods[idx].Date = parsedDate.UTC().Format(layout)
			}
		}
	}

	if requestBody.Minimum != "" {

		parseDate, err := time.Parse(layout, requestBody.Minimum)
		if err != nil {
			log.Printf("Got error when decoding: %s", err)
			err = fmt.Errorf("the minimum date is not structured correctly, should be in %s format", layout)
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(ErrorResponse{Message: err.Error()})
			return err
		}
		requestBody.Minimum = parseDate.Format("01-02-2006")
	}

	if requestBody.Maximum != "" {

		parseDate, err := time.Parse(layout, requestBody.Maximum)
		if err != nil {
			log.Printf("Got error when decoding: %s", err)
			err = fmt.Errorf("the maximum date is not structured correctly, should be in %s format", layout)
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(ErrorResponse{Message: err.Error()})
			return err
		}
		requestBody.Minimum = parseDate.Format("01-02-2006")
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

	//Update popularity in background!
	go directFetchAuthorsCountIncrement(requestBody.Ids)

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
	go authorsAppearInSearchCountIncrement(authors)

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

	//Update popularity in background!
	go directFetchQuotesCountIncrement(requestBody.Ids)

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

	//Update popularity in background!
	go appearInSearchCountIncrement(results)

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

	var results []AuthorsView

	//Order by authorid to have definitive order (when for examplke some names rank the same for similarity), same for why quoteid
	//% is same as SIMILARITY but with default threshold 0.3
	dbPointer := db.Table("authors").
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
	go authorsAppearInSearchCountIncrement(results)

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

	//Update popularity in background!
	go appearInSearchCountIncrement(results)

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

	//Update popularity in background!
	go directFetchTopicCountIncrement(requestBody.Id, requestBody.Topic)

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

	err := dbPointer.Limit(100).Find(&result).Error
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

// swagger:route POST /quotes/aod/new AUTHORS setAuthorsOfTheDay
// Sets the author of the day for the given date. It Is password protected TODO: Put in privacy swagger
// responses:
//	200: successResponse

//SetAuthorOfTheDay sets the author of the day (is password protected)
func SetAuthorOfTheDay(rw http.ResponseWriter, r *http.Request) {
	var requestBody Request
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	if len(requestBody.Aods) == 0 {
		json.NewEncoder(rw).Encode(ErrorResponse{Message: "Please supply some authors", StatusCode: http.StatusBadRequest})
		return
	}

	for _, aod := range requestBody.Aods {
		err := setAOD(requestBody.Language, aod.Date, aod.Id)
		if err != nil {
			log.Println(err)
			json.NewEncoder(rw).Encode(ErrorResponse{Message: "Some of the authors (ids) you supplied do not have " + requestBody.Language + " quotes", StatusCode: http.StatusBadRequest})
			return
		}
	}

	json.NewEncoder(rw).Encode(ErrorResponse{Message: "Successfully inserted quote of the day!", StatusCode: http.StatusOK})
}

func setAOD(language string, date string, authorId int) error {
	switch strings.ToLower(language) {
	case "icelandic":
		return db.Exec("insert into aodice (authorid, date) values((select id from authors where id = ? and hasicelandicquotes), ?) on conflict (date) do update set authorid = ?", authorId, date, authorId).Error
	default:
		return db.Exec("insert into aod (authorid, date) values((select id from authors where id = ? and not hasicelandicquotes), ?) on conflict (date) do update set authorid = ?", authorId, date, authorId).Error
	}
}

//SetNewRandomQOD sets a random quote as the qod for today (if language=icelandic is supplied then it adds the random qod to the icelandic qod table)
func setNewRandomAOD(language string) error {
	var authorItem ListItem
	var dbPointer *gorm.DB
	switch strings.ToLower(language) {
	case "icelandic":
		dbPointer = db.Table("authors").Where("hasicelandicquotes")
	default:
		dbPointer = db.Table("authors").Not("hasicelandicquotes")
	}

	log.Println("HEREMatur")
	err := dbPointer.Order("random()").Limit(1).Scan(&authorItem).Error
	if err != nil {
		return err
	}

	log.Println("MASSI:", authorItem)

	return setAOD(language, time.Now().Format("2006-01-02"), authorItem.Id)
}

// swagger:route POST /quotes/qod AUTHORS getAuthorOfTheDay
// Gets the author of the day
// responses:
//	200: randomQuoteResponse

//GetAuthorOfTheDay gets the author of the day
func GetAuthorOfTheDay(rw http.ResponseWriter, r *http.Request) {
	var requestBody Request
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	if requestBody.Language == "" {
		requestBody.Language = "English"
	}

	var author QuoteView
	var dbPointer *gorm.DB
	var err error
	switch strings.ToLower(requestBody.Language) {
	case "icelandic":
		dbPointer = db.Table("aodiceview")
	default:
		dbPointer = db.Table("aodview")
	}

	err = dbPointer.Where("date = current_date").Scan(&author).Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if (QuoteView{}) == author {
		fmt.Println("Setting a brand new AOD for today")
		err = setNewRandomAOD(requestBody.Language)
		if err != nil {
			//TODO: Respond with better error -- and put into swagger -- and add tests
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		switch strings.ToLower(requestBody.Language) {
		case "icelandic":
			err = db.Table("qodiceview").Where("date = current_date").Scan(&author).Error
		default:
			err = db.Table("qodview").Where("date = current_date").Scan(&author).Error
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
	var requestBody Request
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	if requestBody.Language == "" {
		requestBody.Language = "English"
	}
	var authors []QuoteView
	var dbPointer *gorm.DB
	var err error
	switch strings.ToLower(requestBody.Language) {
	case "icelandic":
		dbPointer = db.Table("aodiceview")
	default:
		dbPointer = db.Table("aodview")
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
			dbPointer = db.Table("aodiceview")
		default:
			dbPointer = db.Table("aodview")
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

// swagger:route POST /quotes/qod/new QUOTES setQuoteOfTheDay
// Sets the quote of the day for the given date. It Is password protected TODO: Put in privacy swagger
// responses:
//	200: successResponse

//SetQuoteOfTheyDay sets the quote of the day (is password protected)
func SetQuoteOfTheDay(rw http.ResponseWriter, r *http.Request) {
	var requestBody Request
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	if requestBody.Language == "" {
		requestBody.Language = "English"
	}

	if len(requestBody.Qods) == 0 {
		json.NewEncoder(rw).Encode(ErrorResponse{Message: "Please supply some quotes", StatusCode: http.StatusBadRequest})
		return
	}

	for _, qod := range requestBody.Qods {
		err := setQOD(requestBody.Language, qod.Date, qod.Id)
		if err != nil {
			log.Println(err)
			json.NewEncoder(rw).Encode(ErrorResponse{Message: "Some of the quotes (ids) you supplied are not in " + requestBody.Language, StatusCode: http.StatusBadRequest})
			return
		}
	}

	json.NewEncoder(rw).Encode(ErrorResponse{Message: "Successfully inserted quote of the day!", StatusCode: http.StatusOK})
}

func setQOD(language string, date string, quoteId int) error {
	switch strings.ToLower(language) {
	case "icelandic":
		return db.Exec("insert into qodice (quoteid, date) values((select id from quotes where id = ? and isicelandic), ?) on conflict (date) do update set quoteid = ?", quoteId, date, quoteId).Error
	default:
		return db.Exec("insert into qod (quoteid, date) values((select id from quotes where id = ? and not isicelandic), ?) on conflict (date) do update set quoteid = ?", quoteId, date, quoteId).Error
	}
}

//SetNewRandomQOD sets a random quote as the qod for today (if language=icelandic is supplied then it adds the random qod to the icelandic qod table)
func setNewRandomQOD(language string) error {
	var quoteItem ListItem
	var dbPointer *gorm.DB
	switch strings.ToLower(language) {
	case "icelandic":
		dbPointer = db.Table("quotes").Where("isicelandic")
	default:
		dbPointer = db.Table("quotes").Not("isicelandic").Where("Random() < 0.005")
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
	var requestBody Request
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	if requestBody.Language == "" {
		requestBody.Language = "English"
	}

	var quote QuoteView
	var dbPointer *gorm.DB
	var err error
	switch strings.ToLower(requestBody.Language) {
	case "icelandic":
		dbPointer = db.Table("qodiceview")
	default:
		dbPointer = db.Table("qodview")
	}

	err = dbPointer.Where("date = current_date").Scan(&quote).Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if (QuoteView{}) == quote {
		fmt.Println("Setting a brand new QOD for today")
		err = setNewRandomQOD(requestBody.Language)
		if err != nil {
			//TODO: Respond with better error -- and put into swagger -- and add tests
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		switch strings.ToLower(requestBody.Language) {
		case "icelandic":
			err = db.Table("qodiceview").Where("date = current_date").Scan(&quote).Error
		default:
			err = db.Table("qodview").Where("date = current_date").Scan(&quote).Error
		}
		log.Println(quote)

		if err != nil {
			//TODO: Respond with better error -- and put into swagger -- and add tests
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}

	json.NewEncoder(rw).Encode(quote)
}

// swagger:route POST /quotes/qod QUOTES getQODHistory
// Gets the history for the quotes of the day
// responses:
//	200: qodHistoryResponse

//GetQODHistory gets Qod history starting from some point
func GetQODHistory(rw http.ResponseWriter, r *http.Request) {
	var requestBody Request
	if err := getRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	if requestBody.Language == "" {
		requestBody.Language = "English"
	}

	var quotes []QuoteView
	var dbPointer *gorm.DB
	var err error
	switch strings.ToLower(requestBody.Language) {
	case "icelandic":
		dbPointer = db.Table("qodiceview")
	default:
		dbPointer = db.Table("qodview")
	}

	if requestBody.Minimum != "" {
		dbPointer = dbPointer.Where("date >= ?", requestBody.Minimum)
	}

	if requestBody.Maximum != "" {
		dbPointer = dbPointer.Where("date <= ?", requestBody.Maximum)
	}

	dbPointer = dbPointer.Where("date <= current_date")

	err = dbPointer.Order("date DESC").Find(&quotes).Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if len(quotes) == 0 {
		fmt.Println("Setting a brand new QOD for today")
		err = setNewRandomQOD(requestBody.Language)
		if err != nil {
			//TODO: Respond with better error -- and put into swagger -- and add tests
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		switch strings.ToLower(requestBody.Language) {
		case "icelandic":
			dbPointer = db.Table("qodiceview")
		default:
			dbPointer = db.Table("qodview")
		}

		if requestBody.Minimum != "" {
			dbPointer = dbPointer.Where("date >= ?", requestBody.Minimum)
		}

		if requestBody.Maximum != "" {
			dbPointer = dbPointer.Where("date <= ?", requestBody.Maximum)
		}

		dbPointer = dbPointer.Where("date <= current_date")

		err = dbPointer.Order("date").Find(&quotes).Error

		if err != nil {
			//TODO: Respond with better error -- and put into swagger -- and add tests
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

	}

	json.NewEncoder(rw).Encode(quotes)
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
