package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Skjaldbaka17/quotes-api/handlers"
	"github.com/Skjaldbaka17/quotes-api/structs"
)

type httpRequest func(http.ResponseWriter, *http.Request)

type Set []int

func TestQuotes(t *testing.T) {
	user := createUser(t)
	godUser := getGODModeUser(t)
	t.Run("Get Quotes", func(t *testing.T) {
		t.Run("should return Quotes with id 1, 2 and 3...", func(t *testing.T) {

			var quoteIds = Set{1, 2, 3}
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","ids":  [%s]}`, user.ApiKey, quoteIds.toString()))
			respObj, _ := requestAndReturnArray(jsonStr, GetQuotes)

			if len(respObj) != len(quoteIds) {
				t.Fatalf("got list of length %d but expected list of length %d", len(respObj), len(quoteIds))
			}

			for idx, quote := range respObj {
				if quote.QuoteId != quoteIds[idx] {
					t.Fatalf("got %d, expected %d", quote.QuoteId, quoteIds[idx])
				}
			}
		})

		t.Run("should get Quotes for author with id 1", func(t *testing.T) {
			authorId := 1
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","authorId":  %d}`, user.ApiKey, authorId))
			respObj, _ := requestAndReturnArray(jsonStr, GetQuotes)

			if len(respObj) == 0 {
				t.Fatalf("got list of length 0 but expected some quotes, response : %+v", respObj)
			}

			if respObj[0].AuthorId != authorId {
				t.Fatalf("got quotes for author with id %d but expected quotes for the author with id %d, respObj: %+v", respObj[0].AuthorId, authorId, respObj)
			}
		})

	})

	t.Run("Quoteslist Test", func(t *testing.T) {

		t.Run("Should return first 50 quotes (by quoteId)", func(t *testing.T) {

			pageSize := 50
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","pageSize": %d}`, user.ApiKey, pageSize))

			respObj, _ := requestAndReturnArray(jsonStr, GetQuotesList)

			if len(respObj) != 50 {
				t.Fatalf("got list of length %d, but expected list of length %d", len(respObj), pageSize)
			}

			firstQuote := respObj[0]
			if firstQuote.QuoteId != 1 {
				t.Fatalf("got %d, want quote with id 1. Resp: %+v", firstQuote.QuoteId, firstQuote)
			}

		})

		t.Run("Should return first quotes, in Icelandic", func(t *testing.T) {

			language := "icelandic"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language": "%s"}`, user.ApiKey, language))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetQuotesList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstQuote := respObj[0]

			if !firstQuote.IsIcelandic {
				t.Fatalf("got %+v, but expected a quote in Icelandic.", firstQuote)
			}

		})

		t.Run("Should return first quotes in reverse quoteId order (i.e. first quote has id larger than 639.028)", func(t *testing.T) {

			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","orderConfig":{"reverse":%s}}`, user.ApiKey, "true"))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetQuotesList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstQuote := respObj[0]

			if firstQuote.QuoteId < 639028 {
				t.Fatalf("got %+v, but want quote with larger quoteid i.e. want last quote in db", firstQuote)
			}

		})

		t.Run("Should return first quotes starting from id 300.000  (i.e. greater than or equal to 300.000)", func(t *testing.T) {
			minimum := 300000
			orderBy := "id"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","orderConfig":{"orderBy":"%s","minimum":"%d"}}`, user.ApiKey, orderBy, minimum))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetQuotesList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstQuote := respObj[0]

			if firstQuote.QuoteId < minimum {
				t.Fatalf("got %+v, want quote that has id larger or equal to 300.000", firstQuote)
			}

		})

		t.Run("Should return quotes with less than or equal to 5 letters in the quote", func(t *testing.T) {

			maximum := 5
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","orderConfig":{"orderBy":"length","maximum":"%d"}}`, user.ApiKey, maximum))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetQuotesList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstQuote := respObj[0]

			if len(firstQuote.Quote) > 5 {
				t.Fatalf("got %+v, but expected a quote that has no more than 5 letters", firstQuote)
			}

		})

		t.Run("Should return first quotes with quote-length at least 10 an most 11", func(t *testing.T) {

			minimum := 10
			maximum := 11
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","orderConfig":{"orderBy":"length","maximum":"%d", "minimum":"%d"}}`, user.ApiKey, maximum, minimum))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetQuotesList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstQuote := respObj[0]

			if len(firstQuote.Quote) != 10 {
				t.Fatalf("got %+v, but expected a quote that has no fewer than 10 letters", firstQuote)
			}

		})

		t.Run("Should return first Quotes with less than letters in the quote in total in reversed order (start with those quotes of length 10)", func(t *testing.T) {

			maximum := 10
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","orderConfig":{"orderBy":"length","maximum":"%d","reverse":true}}`, user.ApiKey, maximum))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetQuotesList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstQuote := respObj[0]

			if len(firstQuote.Quote) != 10 {
				t.Fatalf("got %+v, but expected a quote that has 10 letters", firstQuote)
			}

		})

		t.Run("Should return first 50 quotes (ordered by most popular, i.e. DESC count)", func(t *testing.T) {
			handlers.DirectFetchQuotesCountIncrement([]int{1})

			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","orderConfig":{"orderBy":"%s"}}`, user.ApiKey, "popularity"))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetQuotesList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstQuote := respObj[0]

			if firstQuote.QuoteCount == 0 {
				t.Fatalf("got %+v, but expected a quote that has more than 0 popularity count", firstQuote)
			}

		})

		t.Run("Should return first 50 quotes in reverse popularity order (i.e. least popular first i.e. ASC count)", func(t *testing.T) {

			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","orderConfig":{"orderBy":"%s","reverse":true}}`, user.ApiKey, "popularity"))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetQuotesList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstQuote := respObj[0]

			if firstQuote.QuoteCount != 0 {
				t.Fatalf("got %+v, but expected an author that has 0 popularity count", firstQuote)
			}

		})

		t.Run("Should return first 100 Quotes", func(t *testing.T) {
			pageSize := 100
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","pageSize":%d}}`, user.ApiKey, pageSize))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetQuotesList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			if len(respObj) != 100 {
				t.Fatalf("got %d nr of quotes, but expected %d quotes", len(respObj), pageSize)
			}
		})

		t.Run("Should return the next 50 quotes starting from quoteId 250.000 (i.e. pagination, page 1, quoteId order)", func(t *testing.T) {

			pageSize := 100
			minimum := 250000
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","pageSize":%d, "orderConfig":{"minimum":"%d"}}}`, user.ApiKey, pageSize, minimum))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetQuotesList)

			objToFetch := respObj[50]

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			if respObj[0].QuoteId < minimum {
				t.Fatalf("got %+v, but expected quote with a higher quoteid than %d", len(respObj), minimum)
			}

			pageSize = 50
			page := 1
			jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","pageSize":%d, "page":%d, "orderConfig":{"minimum":"%d"}}}`, user.ApiKey, pageSize, page, minimum))

			respObj, errResponse = requestAndReturnArray(jsonStr, GetQuotesList)

			if objToFetch.QuoteId != respObj[0].QuoteId {
				t.Fatalf("got %+v, but expected %+v", respObj[0], objToFetch)
			}

		})

	})

	t.Run("Quote of the day", func(t *testing.T) {

		t.Run("Should set / Overwrite Quote of the day", func(t *testing.T) {

			quoteId := 1
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","qods": [{"id":%d, "date":""}]}`, godUser.ApiKey, quoteId))
			_, response := requestAndReturnArray(jsonStr, SetQuoteOfTheDay)
			if response.StatusCode != 200 {
				t.Fatalf("Expected a succesful insert but got %+v", response)
			}

		})

		t.Run("Should set QOD for 12-22-2020 and 12-21-2020", func(t *testing.T) {
			//TODO: add to test that the quotes where actually input into the DB
			quoteId1 := 2
			date1 := "2020-12-22"
			date2 := "2020-12-21"
			quoteId2 := 3
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","qods": [{"id":%d, "date":"%s"},{"id":%d, "date":"%s"}]}`, godUser.ApiKey, quoteId1, date1, quoteId2, date2))
			_, response := requestAndReturnArray(jsonStr, SetQuoteOfTheDay)
			if response.StatusCode != 200 {
				t.Fatalf("Expected a succesful insert but got %+v", response)
			}

		})

		t.Run("Should get Quote of the day", func(t *testing.T) {
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language":"%s"}`, user.ApiKey, "english"))
			quote := requestAndReturnSingle(jsonStr, GetQuoteOfTheDay)

			if quote.Quote == "" {
				t.Fatalf("Expected the quote of the day but got an empty quote %+v", quote)
			}

			const layout = "2006-01-02T15:04:05Z" //The date needed for reference always
			date, _ := time.Parse(layout, quote.Date)
			if date.Format("01-02-2006") != time.Now().Format("01-02-2006") {
				t.Fatalf("Expected the quote for the date %s but got QOD for date %s i.e. %+v", time.Now().Format("01-02-2006"), date.Format("01-02-2006"), quote)
			}

		})

		t.Run("Should get complete history of quote of the day", func(t *testing.T) {
			//Input a quote in history for testing
			quoteId := 1111
			date := "1998-06-16"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","qods": [{"id":%d, "date":"%s"}]}`, godUser.ApiKey, quoteId, date))
			_, response := requestAndReturnArray(jsonStr, SetQuoteOfTheDay)
			if response.StatusCode != 200 {
				t.Fatalf("Expected a succesful insert but got %+v", response)
			}

			//Get History:

			jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language":"%s"}`, user.ApiKey, "english"))
			quotes, _ := requestAndReturnArray(jsonStr, GetQODHistory)

			if len(quotes) == 0 {
				t.Fatalf("Expected the history of QOD but got an empty list: %+v", quotes)
			}

			containsBirfdayQuote := false
			containsTodayQuote := false
			const layout = "2006-01-02T15:04:05Z" //The date needed for reference always
			for _, quote := range quotes {
				date, _ := time.Parse(layout, quote.Date)
				if date.Format("01-02-2006") == time.Now().Format("01-02-2006") {
					containsTodayQuote = true
				}

				if date.Format("01-02-2006") == "06-16-1998" {
					containsBirfdayQuote = true
				}
			}

			if !containsBirfdayQuote {
				t.Fatalf("QOD history should contain the QOD for birfday but does not: %+v", quotes)
			}

			if !containsTodayQuote {
				t.Fatalf("QOD history should contain the QOD for today but does not: %+v", quotes)
			}

		})

		t.Run("Should get history of QOD starting from June 4th 2021", func(t *testing.T) {

			//Input a quote in history for testing
			quoteId := 666
			date := "2021-06-04"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","qods": [{"id":%d, "date":"%s"}]}`, godUser.ApiKey, quoteId, date))
			_, response := requestAndReturnArray(jsonStr, SetQuoteOfTheDay)
			if response.StatusCode != 200 {
				t.Fatalf("Expected a succesful insert but got %+v", response)
			}

			//Get History:

			minimum := "2021-06-04"
			jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language":"%s", "minimum":"%s"}`, user.ApiKey, "english", minimum))
			quotes, _ := requestAndReturnArray(jsonStr, GetQODHistory)

			if len(quotes) == 0 {
				t.Fatalf("Expected the history of QOD but got an empty list: %+v", quotes)
			}

			const layout = "2006-01-02T15:04:05Z" //The date needed for reference always
			compareDate, _ := time.Parse(layout, "2021-06-04")
			compareYear := compareDate.Year()
			compareMonth := compareDate.Month()
			compareDay := compareDate.Day()
			containsQuoteNotInRange := false
			containsFourthOfJuneQuote := false
			for _, quote := range quotes {
				date, _ := time.Parse(layout, quote.Date)
				yearOfQuote := date.Year()
				monthOfQuote := date.Month()
				dayOfQuote := date.Day()

				if yearOfQuote < compareYear || (yearOfQuote == compareYear && monthOfQuote < compareMonth) || (yearOfQuote == compareYear && monthOfQuote == compareMonth && dayOfQuote < compareDay) {
					containsQuoteNotInRange = true
				}

				if date.Format("01-02-2006") == "06-04-2021" {
					containsFourthOfJuneQuote = true
				}
			}

			if containsQuoteNotInRange {
				t.Fatalf("QOD history contains an earlier quote than was requested: %+v", quotes)
			}

			if !containsFourthOfJuneQuote {
				t.Fatalf("QOD history should contain the QOD for 4th of june 20201 but does not: %+v", quotes)
			}

		})

	})

	t.Run("Random Quotes", func(t *testing.T) {

		//The test calls the function twice to test if the function returns two different quotes
		t.Run("Should return a random quote", func(t *testing.T) {

			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s"}`, user.ApiKey))
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)

			if firstRespObj.Quote == "" {
				t.Fatalf("Expected a random quote but got an empty quote")
			}

			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)

			if secondRespObj.QuoteId == firstRespObj.QuoteId {
				t.Fatalf("Expected two different quotes but got the same quote twice which is higly improbable")
			}
		})

		t.Run("Should return a random quote from Teddy Roosevelt (given authorId)", func(t *testing.T) {

			teddyName := "Theodore Roosevelt"
			teddyAuthor := getAuthor(teddyName, user.ApiKey)
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","authorId": %d}`, user.ApiKey, teddyAuthor.Id))
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)

			if firstRespObj.Name != teddyName {
				t.Fatalf("got %s, expected %s", firstRespObj.Name, teddyName)
			}

			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)

			if secondRespObj.AuthorId != firstRespObj.AuthorId {
				t.Fatalf("got author with id %d, expected author with id %d", secondRespObj.AuthorId, firstRespObj.AuthorId)
			}

			if secondRespObj.QuoteId == firstRespObj.QuoteId {
				t.Fatalf("got quote %s, expected a random different quote", secondRespObj.Quote)
			}

		})

		t.Run("Should return a random quote from topic 'motivational' (given topicId)", func(t *testing.T) {

			topicName := "motivational"
			topicId := getTopicId(topicName, user.ApiKey)
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","topicId": %d}`, user.ApiKey, topicId))
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if firstRespObj.TopicName != topicName {
				t.Fatalf("got topicname: %s, expected %s", firstRespObj.TopicName, topicName)
			}
			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if secondRespObj.TopicId != firstRespObj.TopicId {
				t.Fatalf("got topic with id %d, expected topic with id %d", secondRespObj.TopicId, firstRespObj.TopicId)
			}

			if secondRespObj.QuoteId == firstRespObj.QuoteId {
				t.Fatalf("got quote %s, expected a random different quote", secondRespObj.Quote)
			}
		})

		t.Run("Should return a random English quote", func(t *testing.T) {

			language := "english"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language": "%s"}`, user.ApiKey, language))
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if firstRespObj.IsIcelandic {
				t.Fatalf("first response, got an IcelandicQuote but expected an English quote")
			}
			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if secondRespObj.IsIcelandic {
				t.Fatalf("second response, got an IcelandicQuote but expected an English quote")
			}

			if secondRespObj.QuoteId == firstRespObj.QuoteId {
				t.Fatalf("got quote %s, expected a random different quote", secondRespObj.Quote)
			}
		})

		t.Run("Should return a random Icelandic quote", func(t *testing.T) {

			language := "Icelandic"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language": "%s"}`, user.ApiKey, language))
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if !firstRespObj.IsIcelandic {
				t.Fatalf("first response, got an EnglishQuote but expected an Icelandic quote")
			}
			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if !secondRespObj.IsIcelandic {
				t.Fatalf("second response, got an EnglishQuote, %+v, but expected an Icelandic quote", secondRespObj)
			}

			if secondRespObj.QuoteId == firstRespObj.QuoteId {
				t.Fatalf("got quote %s, expected a random different quote", secondRespObj.Quote)
			}
		})

		t.Run("Should return a random quote containing the searchString 'love'", func(t *testing.T) {

			searchString := "love"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","searchString":"%s"}`, user.ApiKey, searchString))
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			regexStub := searchString[:3]
			m1 := regexp.MustCompile(regexStub)
			if !m1.Match([]byte(strings.ToLower(firstRespObj.Quote))) {
				t.Fatalf("first response, got the quote %+v that does not contain the searchString %s", firstRespObj, regexStub)
			}

			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if !m1.Match([]byte(strings.ToLower(secondRespObj.Quote))) {
				t.Fatalf("second response, got the quote %+v that does not contain the searchString %s", secondRespObj, regexStub)
			}

			if secondRespObj.QuoteId == firstRespObj.QuoteId {
				t.Fatalf("got quote %s, expected a random different quote", secondRespObj.Quote)
			}

		})

		t.Run("Should return a random Icelandic quote containing the searchString '??itt'", func(t *testing.T) {

			searchString := "??itt"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","searchString":"%s"}`, user.ApiKey, searchString))
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			m1 := regexp.MustCompile(searchString)
			if !m1.Match([]byte(firstRespObj.Quote)) {
				t.Fatalf("first response, got the quote %+v that does not contain the searchString %s", firstRespObj, searchString)
			}

			if !firstRespObj.IsIcelandic {
				t.Fatalf("first response, got the quote %+v which is in English but expected it to be in icelandic", firstRespObj)
			}

			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if !m1.Match([]byte(secondRespObj.Quote)) {
				t.Fatalf("second response, got the quote %+v that does not contain the searchString %s", secondRespObj, searchString)
			}

			if secondRespObj.QuoteId == firstRespObj.QuoteId {
				t.Fatalf("got quote %s, expected a random different quote", secondRespObj.Quote)
			}
		})

		t.Run("Should return a random quote containing the searchString 'strong' from the topic 'inspirational' (given topicId)", func(t *testing.T) {

			topicName := "inspirational"
			topicId := getTopicId(topicName, user.ApiKey)
			searchString := "strong"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","searchString":"%s","topicId": %d}`, user.ApiKey, searchString, topicId))
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)

			if firstRespObj.TopicName != topicName {
				t.Fatalf("got %s, expected %s", firstRespObj.TopicName, topicName)
			}

			m1 := regexp.MustCompile(searchString)
			if !m1.Match([]byte(firstRespObj.Quote)) {
				t.Fatalf("first response, got the quote %+v that does not contain the searchString %s", firstRespObj, searchString)
			}

			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if !m1.Match([]byte(secondRespObj.Quote)) {
				t.Fatalf("second response, got the quote %+v that does not contain the searchString %s", secondRespObj, searchString)
			}

			if secondRespObj.TopicId != firstRespObj.TopicId {
				t.Fatalf("got topic with id %d, expected topic with id %d", secondRespObj.TopicId, firstRespObj.TopicId)
			}

			if secondRespObj.QuoteId == firstRespObj.QuoteId {
				t.Fatalf("got quote %s, expected a random different quote... Remember that this is a random function and therefore there is a chance the same quote is fetched twice.", secondRespObj.Quote)
			}
		})

	})

	t.Cleanup(func() {
		log.Println("CLEANUP TestQuotes!")
		// Delete from qod
		handlers.Db.Exec("DELETE FROM qod")
		handlers.Db.Exec("DELETE FROM qodice")
	})
}

func getTopicId(topicName string, apiKey string) int {

	var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","topic": "%s"}`, apiKey, topicName))
	request, _ := http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	GetTopic(response, request)

	var respObj []structs.TestApiResponse
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj[0].TopicId
}

func (set *Set) toString() string {
	var IDs []string
	for _, i := range *set {
		IDs = append(IDs, strconv.Itoa(i))
	}

	return strings.Join(IDs, ", ")
}

func getAuthor(searchString string, apiKey string) structs.TestApiResponse {

	var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","searchString": "%s"}`, apiKey, searchString))
	request, _ := http.NewRequest(http.MethodGet, "/api/search/authors", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	SearchAuthorsByString(response, request)

	var respObj []structs.TestApiResponse
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj[0]
}

func getAuthorsById(authorIds Set, apiKey string) []structs.TestApiResponse {

	var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","ids": [%s]}`, apiKey, authorIds.toString()))
	request, _ := http.NewRequest(http.MethodGet, "/api/search/authors", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	GetAuthorsById(response, request)

	var respObj []structs.TestApiResponse
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj
}

func getQuotes(searchString string, apiKey string) []structs.TestApiResponse {

	var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","searchString": "%s"}`, apiKey, searchString))
	request, _ := http.NewRequest(http.MethodGet, "/api/search/quotes", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	SearchQuotesByString(response, request)

	var respObj []structs.TestApiResponse
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj
}

func getRequestAndResponseForTest(jsonStr []byte) (*httptest.ResponseRecorder, *http.Request) {
	request, _ := http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()
	return response, request
}

//TODO: Give a better name,more intuitive
func getObjNr26(searchString string, fn httpRequest, apiKey string) (structs.TestApiResponse, error) {
	pageSize := 100
	var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","searchString": "%s", "pageSize":%d}`, apiKey, searchString, pageSize))
	request, _ := http.NewRequest(http.MethodPost, "/api/search", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	fn(response, request)

	var respObj []structs.TestApiResponse
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)

	if pageSize != len(respObj) {
		return structs.TestApiResponse{}, fmt.Errorf("got list of length %d but expected %d", len(respObj), pageSize)
	}

	return respObj[25], nil
}

func requestAndReturnSingle(jsonStr []byte, fn httpRequest) structs.TestApiResponse {
	response, request := getRequestAndResponseForTest(jsonStr)
	fn(response, request)

	var respObj structs.TestApiResponse

	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj
}

func requestAndReturnArray(jsonStr []byte, fn httpRequest) ([]structs.TestApiResponse, structs.ErrorResponse) {
	response, request := getRequestAndResponseForTest(jsonStr)
	fn(response, request)
	var respObj []structs.TestApiResponse
	var errorResp structs.ErrorResponse
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	_ = json.Unmarshal(response.Body.Bytes(), &errorResp)
	if errorResp.StatusCode == 0 {
		errorResp.StatusCode = response.Result().StatusCode
	}
	return respObj, errorResp
}

func requestAndReturnArrayAuthors(jsonStr []byte, fn httpRequest) ([]structs.TestApiResponse, structs.ErrorResponse) {
	response, request := getRequestAndResponseForTest(jsonStr)
	fn(response, request)
	var respObj []structs.TestApiResponse
	var errorResp structs.ErrorResponse
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	_ = json.Unmarshal(response.Body.Bytes(), &errorResp)
	if errorResp.StatusCode == 0 {
		errorResp.StatusCode = response.Result().StatusCode
	}
	return respObj, errorResp
}
