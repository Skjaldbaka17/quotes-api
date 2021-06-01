package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	t.Run("Search Quotes By String", func(t *testing.T) {
		t.Run("easy search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "Float like a butterfly sting like a bee"}`)
			respObj := requestAndReturnArray(jsonStr, SearchQuotesByString)
			firstAuthor := respObj[0].Name
			want := "Muhammad Ali"
			if firstAuthor != want {
				t.Errorf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("intermediate search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "bee sting like a butterfly"}`)
			respObj := requestAndReturnArray(jsonStr, SearchQuotesByString)
			firstAuthor := respObj[0].Name
			want := "Muhammad Ali"
			if firstAuthor != want {
				t.Errorf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("hard search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "bee butterfly float"}`)
			respObj := requestAndReturnArray(jsonStr, SearchQuotesByString)
			firstAuthor := respObj[0].Name
			want := "Muhammad Ali"
			if firstAuthor != want {
				t.Errorf("got %q, want %q", firstAuthor, want)
			}
		})

	})

	t.Run("Search Authors By String", func(t *testing.T) {

		t.Run("easy search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "Friedrich Nietzsche"}`)
			respObj := requestAndReturnArray(jsonStr, SearchAuthorsByString)
			firstAuthor := respObj[0].Name
			want := "Friedrich Nietzsche"
			if firstAuthor != want {
				t.Errorf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("intermediate search should Return list of quotes with Joseph Stalin as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "Stalin jseph"}`)
			respObj := requestAndReturnArray(jsonStr, SearchAuthorsByString)
			firstAuthor := respObj[0].Name
			want := "Joseph Stalin"
			if firstAuthor != want {
				t.Errorf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("hard search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "Niet Friedric"}`)
			respObj := requestAndReturnArray(jsonStr, SearchAuthorsByString)
			firstAuthor := respObj[0].Name
			want := "Friedrich Nietzsche"
			if firstAuthor != want {
				t.Errorf("got %q, want %q", firstAuthor, want)
			}
		})

	})

	t.Run("Search by String", func(t *testing.T) {

		t.Run("searching for author", func(t *testing.T) {
			t.Run("easy search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {

				var jsonStr = []byte(`{"searchString": "Friedrich Nietzsche"}`)
				respObj := requestAndReturnArray(jsonStr, SearchByString)
				firstAuthor := respObj[1].Name //Use index 1 because in index 0 there is an author talking extensively about Nietzsche
				want := "Friedrich Nietzsche"
				if firstAuthor != want {
					t.Errorf("got %q, want %q", firstAuthor, want)
				}
			})

			t.Run("hard search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {

				var jsonStr = []byte(`{"searchString": "Nietshe Friedr"}`)
				respObj := requestAndReturnArray(jsonStr, SearchByString)
				firstAuthor := respObj[1].Name //Use index 1 because in index 0 there is an author talking extensively about Nietzsche
				want := "Friedrich Nietzsche"
				if firstAuthor != want {
					t.Errorf("got %q, want %q", firstAuthor, want)
				}
			})
		})

		t.Run("searching for quote", func(t *testing.T) {
			t.Run("easy search should return list of quotes with Martin Luther as first author", func(t *testing.T) {

				var jsonStr = []byte(`{"searchString": "If you are not allowed to Laugh in Heaven"}`)
				respObj := requestAndReturnArray(jsonStr, SearchByString)
				firstAuthor := respObj[1].Name //Use index 1 because in index 0 there is an author talking extensively about Nietzsche
				want := "Martin Luther"
				if firstAuthor != want {
					t.Errorf("got %q, want %q", firstAuthor, want)
				}
			})
		})

	})

	t.Run("Search Pagination Test", func(t *testing.T) {

		t.Run("Search By string pagination", func(t *testing.T) {

			searchString := "Love"
			obj26, err := getObjNr26(searchString, SearchByString)

			if err != nil {
				t.Error(err)
			}
			//Next request to check if same dude in position 0 given that pageSize is 25 and same search parameters
			pageSize := 25
			jsonStr := []byte(fmt.Sprintf(`{"searchString": "%s", "pageSize":%d, "page":1}`, searchString, pageSize))
			request, _ := http.NewRequest(http.MethodPost, "/api/search", bytes.NewBuffer(jsonStr))
			response := httptest.NewRecorder()

			SearchByString(response, request)
			var respObj []QuoteView
			_ = json.Unmarshal(response.Body.Bytes(), &respObj)

			if pageSize != len(respObj) {
				t.Errorf("got list of length %d but expected %d", len(respObj), pageSize)
			}

			if respObj[0] != obj26 {
				t.Errorf("got %+v, want %+v", respObj[0], obj26)
			}
		})

		t.Run("Search Authors By string pagination", func(t *testing.T) {

			searchString := "Martin"
			obj26, err := getObjNr26(searchString, SearchAuthorsByString)

			if err != nil {
				t.Error(err)
			}
			//Next request to check if same dude in position 0 given that pageSize is 25 and same search parameters
			pageSize := 25
			jsonStr := []byte(fmt.Sprintf(`{"searchString": "%s", "pageSize":%d, "page":1}`, searchString, pageSize))
			request, _ := http.NewRequest(http.MethodPost, "/api/search/authors", bytes.NewBuffer(jsonStr))
			response := httptest.NewRecorder()

			SearchAuthorsByString(response, request)

			var respObj []QuoteView
			_ = json.Unmarshal(response.Body.Bytes(), &respObj)

			if pageSize != len(respObj) {
				t.Errorf("got list of length %d but expected %d", len(respObj), pageSize)
			}

			if respObj[0] != obj26 {
				t.Errorf("got %+v, want %+v", respObj[0], obj26)
			}
		})

		t.Run("Search Quotes By string pagination", func(t *testing.T) {

			searchString := "Hate"
			obj26, err := getObjNr26(searchString, SearchQuotesByString)

			if err != nil {
				t.Error(err)
			}
			//Next request to check if same dude in position 0 given that pageSize is 25 and same search parameters
			pageSize := 25
			jsonStr := []byte(fmt.Sprintf(`{"searchString": "%s", "pageSize":%d, "page":1}`, searchString, pageSize))
			request, _ := http.NewRequest(http.MethodPost, "/api/search/quotes", bytes.NewBuffer(jsonStr))
			response := httptest.NewRecorder()

			SearchQuotesByString(response, request)

			var respObj []QuoteView
			_ = json.Unmarshal(response.Body.Bytes(), &respObj)

			if pageSize != len(respObj) {
				t.Errorf("got list of length %d but expected %d", len(respObj), pageSize)
			}

			if respObj[0] != obj26 {
				t.Errorf("got %+v, want %+v", respObj[0], obj26)
			}
		})

	})

}

func TestAuthors(t *testing.T) {

	t.Run("should return Author with id 1", func(t *testing.T) {

		authorId := Set{1}
		var jsonStr = []byte(fmt.Sprintf(`{"ids": [%s]}`, authorId.toString()))
		respObj := requestAndReturnArray(jsonStr, GetAuthorsById)
		firstAuthor := respObj[0]
		if firstAuthor.Authorid != authorId[0] {
			t.Errorf("got %d, want %d", firstAuthor.Authorid, authorId[0])
		}

	})

	t.Run("Random author", func(t *testing.T) {
		t.Run("Should return a random author", func(t *testing.T) { t.Skip() })
		t.Run("Should return a random Icelandic author", func(t *testing.T) { t.Skip() })
		t.Run("Should return a random English Author", func(t *testing.T) { t.Skip() })
	})

}

func TestQuotes(t *testing.T) {
	t.Run("should return Quotes with id 1, 2 and 3...", func(t *testing.T) {

		var quoteIds = Set{1, 2, 3}
		var jsonStr = []byte(fmt.Sprintf(`{"ids":  [%s]}`, quoteIds.toString()))
		respObj := requestAndReturnArray(jsonStr, GetQuotesById)

		if len(respObj) != len(quoteIds) {
			t.Errorf("got list of length %d but expected list of length %d", len(respObj), len(quoteIds))
		}

		for idx, quote := range respObj {
			if quote.Quoteid != quoteIds[idx] {
				t.Errorf("got %d, expected %d", quote.Quoteid, quoteIds[idx])
			}
		}
	})

	t.Run("Random Quote", func(t *testing.T) {

		//The test calls the function twice to test if the function returns two different quotes
		t.Run("Should return a random quote", func(t *testing.T) {

			var jsonStr = []byte(`{}`)
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)

			if firstRespObj.Quote == "" {
				t.Errorf("Expected a random quote but got an empty quote")
			}

			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)

			if secondRespObj.Quoteid == firstRespObj.Quoteid {
				t.Errorf("Expected two different quotes but got the same quote twice which is higly improbable")
			}
		})

		t.Run("Should return a random quote from Teddy Roosevelt (given authorId)", func(t *testing.T) {

			teddyName := "Theodore Roosevelt"
			teddyAuthor := getAuthor(teddyName)
			var jsonStr = []byte(fmt.Sprintf(`{"authorId": %d}`, teddyAuthor.Authorid))
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)

			if firstRespObj.Name != teddyName {
				t.Errorf("got %s, expected %s", firstRespObj.Name, teddyName)
			}

			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)

			if secondRespObj.Authorid != firstRespObj.Authorid {
				t.Errorf("got author with id %d, expected author with id %d", secondRespObj.Authorid, firstRespObj.Authorid)
			}

			if secondRespObj.Quoteid == firstRespObj.Quoteid {
				t.Errorf("got quote %s, expected a random different quote", secondRespObj.Quote)
			}

		})

		t.Run("Should return a random quote from topic 'motivational' (given topicId)", func(t *testing.T) {

			topicName := "motivational"
			topicId := getTopicId(topicName)
			var jsonStr = []byte(fmt.Sprintf(`{"topicId": %d}`, topicId))
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if firstRespObj.Topicname != topicName {
				t.Errorf("got %s, expected %s", firstRespObj.Topicname, topicName)
			}
			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if secondRespObj.Topicid != firstRespObj.Topicid {
				t.Errorf("got topic with id %d, expected topic with id %d", secondRespObj.Topicid, firstRespObj.Topicid)
			}

			if secondRespObj.Quoteid == firstRespObj.Quoteid {
				t.Errorf("got quote %s, expected a random different quote", secondRespObj.Quote)
			}
		})

		t.Run("Should return a random English quote", func(t *testing.T) {

			language := "english"
			var jsonStr = []byte(fmt.Sprintf(`{"language": "%s"}`, language))
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if firstRespObj.Isicelandic {
				t.Errorf("first response, got an IcelandicQuote but expected an English quote")
			}
			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if secondRespObj.Isicelandic {
				t.Errorf("second response, got an IcelandicQuote but expected an English quote")
			}

			if secondRespObj.Quoteid == firstRespObj.Quoteid {
				t.Errorf("got quote %s, expected a random different quote", secondRespObj.Quote)
			}
		})

		t.Run("Should return a random Icelandic quote", func(t *testing.T) {

			language := "Icelandic"
			var jsonStr = []byte(fmt.Sprintf(`{"language": "%s"}`, language))
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if !firstRespObj.Isicelandic {
				t.Errorf("first response, got an EnglishQuote but expected an Icelandic quote")
			}
			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if !secondRespObj.Isicelandic {
				t.Errorf("second response, got an EnglishQuote, %+v, but expected an Icelandic quote", secondRespObj)
			}

			if secondRespObj.Quoteid == firstRespObj.Quoteid {
				t.Errorf("got quote %s, expected a random different quote", secondRespObj.Quote)
			}
		})

		t.Run("Should return a random quote containing the searchString 'love'", func(t *testing.T) {

			searchString := "love"
			var jsonStr = []byte(fmt.Sprintf(`{"searchString":"%s"}`, searchString))
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			m1 := regexp.MustCompile(searchString)
			if !m1.Match([]byte(firstRespObj.Quote)) {
				t.Errorf("first response, got the quote %+v that does not contain the searchString %s", firstRespObj, searchString)
			}

			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if !m1.Match([]byte(secondRespObj.Quote)) {
				t.Errorf("second response, got the quote %+v that does not contain the searchString %s", secondRespObj, searchString)
			}

			if secondRespObj.Quoteid == firstRespObj.Quoteid {
				t.Errorf("got quote %s, expected a random different quote", secondRespObj.Quote)
			}

		})

		t.Run("Should return a random Icelandic quote containing the searchString 'þitt'", func(t *testing.T) {

			searchString := "þitt"
			var jsonStr = []byte(fmt.Sprintf(`{"searchString":"%s"}`, searchString))
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			m1 := regexp.MustCompile(searchString)
			if !m1.Match([]byte(firstRespObj.Quote)) {
				t.Errorf("first response, got the quote %+v that does not contain the searchString %s", firstRespObj, searchString)
			}

			if !firstRespObj.Isicelandic {
				t.Errorf("first response, got the quote %+v which is in English but expected it to be in icelandic", firstRespObj)
			}

			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if !m1.Match([]byte(secondRespObj.Quote)) {
				t.Errorf("second response, got the quote %+v that does not contain the searchString %s", secondRespObj, searchString)
			}

			if secondRespObj.Quoteid == firstRespObj.Quoteid {
				t.Errorf("got quote %s, expected a random different quote", secondRespObj.Quote)
			}
		})

		t.Run("Should return a random quote containing the searchString 'strong' from the topic 'inspirational' (given topicId)", func(t *testing.T) {

			topicName := "inspirational"
			topicId := getTopicId(topicName)
			searchString := "strong"
			var jsonStr = []byte(fmt.Sprintf(`{"searchString":"%s","topicId": %d}`, searchString, topicId))
			firstRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)

			if firstRespObj.Topicname != topicName {
				t.Errorf("got %s, expected %s", firstRespObj.Topicname, topicName)
			}

			m1 := regexp.MustCompile(searchString)
			if !m1.Match([]byte(firstRespObj.Quote)) {
				t.Errorf("first response, got the quote %+v that does not contain the searchString %s", firstRespObj, searchString)
			}

			secondRespObj := requestAndReturnSingle(jsonStr, GetRandomQuote)
			if !m1.Match([]byte(secondRespObj.Quote)) {
				t.Errorf("second response, got the quote %+v that does not contain the searchString %s", secondRespObj, searchString)
			}

			if secondRespObj.Topicid != firstRespObj.Topicid {
				t.Errorf("got topic with id %d, expected topic with id %d", secondRespObj.Topicid, firstRespObj.Topicid)
			}

			if secondRespObj.Quoteid == firstRespObj.Quoteid {
				t.Errorf("got quote %s, expected a random different quote... Remember that this is a random function and therefore there is a chance the same quote is fetched twice.", secondRespObj.Quote)
			}
		})

	})

}

func TestTopics(t *testing.T) {

	t.Run("Should return the possible English topics as a list of objects", func(t *testing.T) {

		var nrOfEnglishTopics int = 13
		var language string = "English"
		var jsonStr = []byte(fmt.Sprintf(`{"language": "%s"}`, language))
		request, _ := http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()

		GetTopics(response, request)

		var respObj []ListItem
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)

		if len(respObj) != nrOfEnglishTopics {
			t.Errorf("got %d number of topics, expected %d", len(respObj), nrOfEnglishTopics)
		}
	})

	t.Run("Should return the possible Icelandic topics as a list of objects", func(t *testing.T) {

		var nrOfIcelandicTopics int = 7
		var language string = "Icelandic"
		var jsonStr = []byte(fmt.Sprintf(`{"language": "%s"}`, language))
		request, _ := http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()
		GetTopics(response, request)
		var respObj []ListItem
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)

		if len(respObj) != nrOfIcelandicTopics {
			t.Errorf("got %d number of topics, expected %d", len(respObj), nrOfIcelandicTopics)
		}
	})

	t.Run("Should return the first 25 quotes from a topic 'nameOfTopic'", func(t *testing.T) {

		var nameOfTopic string = "inspirational"
		var pageSize int = 25
		var page int = 0
		var jsonStr = []byte(fmt.Sprintf(`{"topic": "%s", "pageSize":%d, "page":%d}`, nameOfTopic, pageSize, page))
		request, _ := http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()
		GetTopic(response, request)
		var respObj []QuoteView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)

		if len(respObj) != pageSize {
			t.Errorf("got %d number of quotes, expected %d as the pagesize", len(respObj), pageSize)
		}

		for _, obj := range respObj {
			if respObj[0].Topicname != nameOfTopic {
				t.Errorf("got %+v but expected a quote with topic %s", obj, nameOfTopic)
			}
		}

	})

	t.Run("Should return the first 25 quotes from a topic with id", func(t *testing.T) {

		topicId := getTopicId("inspirational")
		var pageSize int = 26
		var page int = 0
		var jsonStr = []byte(fmt.Sprintf(`{"id": %d, "pageSize":%d, "page":%d}`, topicId, pageSize, page))
		request, _ := http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()
		GetTopic(response, request)
		var respObj []QuoteView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)

		if len(respObj) != pageSize {
			t.Errorf("got %d number of quotes, expected %d as the pagesize", len(respObj), pageSize)
		}

		for _, obj := range respObj {
			if respObj[0].Topicid != topicId {
				t.Errorf("got %+v but expected a quote with topicId %d", obj, topicId)
			}
		}

	})

	t.Run("Should test pagination for a specific topic, by id", func(t *testing.T) {

		topicId := getTopicId("inspirational")
		//First get the 2nd page's first quote, where pagesize is 25
		var pageSize int = 25
		var page int = 1
		var jsonStr = []byte(fmt.Sprintf(`{"id": %d, "pageSize":%d, "page":%d}`, topicId, pageSize, page))
		request, _ := http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()
		GetTopic(response, request)

		var respObj []QuoteView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)

		obj26 := respObj[0]

		// Then get the first 100 quotes, i.e. first page with pagesize 100
		pageSize = 100
		page = 0
		jsonStr = []byte(fmt.Sprintf(`{"id": %d, "pageSize":%d, "page":%d}`, topicId, pageSize, page))
		request, _ = http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
		response = httptest.NewRecorder()
		GetTopic(response, request)

		_ = json.Unmarshal(response.Body.Bytes(), &respObj)

		if len(respObj) != pageSize {
			t.Errorf("got %d number of quotes, expected %d as the pagesize", len(respObj), pageSize)
		}

		//Compare the 26th object from the 100pagesize request with the 1st object from the 2nd page where pagesize is 25.
		if respObj[25] != obj26 {
			t.Errorf("got %+v but expected %+v", respObj[25], obj26)
		}

	})

}

type Set []int

type httpRequest func(http.ResponseWriter, *http.Request)

func getTopicId(topicName string) int {

	var jsonStr = []byte(fmt.Sprintf(`{"topic": "%s"}`, topicName))
	request, _ := http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	GetTopic(response, request)

	var respObj []QuoteView
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj[0].Topicid
}

func (set *Set) toString() string {
	var IDs []string
	for _, i := range *set {
		IDs = append(IDs, strconv.Itoa(i))
	}

	return strings.Join(IDs, ", ")
}

func getAuthor(searchString string) QuoteView {

	var jsonStr = []byte(fmt.Sprintf(`{"searchString": "%s"}`, searchString))
	request, _ := http.NewRequest(http.MethodGet, "/api/search/authors", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	SearchAuthorsByString(response, request)

	var respObj []QuoteView
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj[0]
}

func getQuotes(searchString string) []QuoteView {

	var jsonStr = []byte(fmt.Sprintf(`{"searchString": "%s"}`, searchString))
	request, _ := http.NewRequest(http.MethodGet, "/api/search/quotes", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	SearchQuotesByString(response, request)

	var respObj []QuoteView
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj
}

func getRequestAndResponseForTest(jsonStr []byte) (*httptest.ResponseRecorder, *http.Request) {
	request, _ := http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()
	return response, request
}

//TODO: Give a better name,more intuitive
func getObjNr26(searchString string, fn httpRequest) (QuoteView, error) {
	pageSize := 100
	var jsonStr = []byte(fmt.Sprintf(`{"searchString": "%s", "pageSize":%d}`, searchString, pageSize))
	request, _ := http.NewRequest(http.MethodPost, "/api/search", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	fn(response, request)

	var respObj []QuoteView
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)

	if pageSize != len(respObj) {
		return QuoteView{}, fmt.Errorf("got list of length %d but expected %d", len(respObj), pageSize)
	}

	return respObj[25], nil
}

func requestAndReturnSingle(jsonStr []byte, fn httpRequest) QuoteView {
	response, request := getRequestAndResponseForTest(jsonStr)
	fn(response, request)

	var respObj QuoteView

	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj
}

func requestAndReturnArray(jsonStr []byte, fn httpRequest) []QuoteView {
	response, request := getRequestAndResponseForTest(jsonStr)
	fn(response, request)
	var respObj []QuoteView
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj
}
