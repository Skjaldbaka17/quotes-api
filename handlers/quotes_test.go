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
			respObj, _ := requestAndReturnArray(jsonStr, SearchQuotesByString)
			firstAuthor := respObj[0].Name
			want := "Muhammad Ali"
			if firstAuthor != want {
				t.Errorf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("intermediate search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "bee sting like a butterfly"}`)
			respObj, _ := requestAndReturnArray(jsonStr, SearchQuotesByString)
			firstAuthor := respObj[0].Name
			want := "Muhammad Ali"
			if firstAuthor != want {
				t.Errorf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("hard search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "bee butterfly float"}`)
			respObj, _ := requestAndReturnArray(jsonStr, SearchQuotesByString)
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
			respObj, _ := requestAndReturnArray(jsonStr, SearchAuthorsByString)
			firstAuthor := respObj[0].Name
			want := "Friedrich Nietzsche"
			if firstAuthor != want {
				t.Errorf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("intermediate search should Return list of quotes with Joseph Stalin as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "Stalin jseph"}`)
			respObj, _ := requestAndReturnArray(jsonStr, SearchAuthorsByString)
			firstAuthor := respObj[0].Name
			want := "Joseph Stalin"
			if firstAuthor != want {
				t.Errorf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("hard search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "Niet Friedric"}`)
			respObj, _ := requestAndReturnArray(jsonStr, SearchAuthorsByString)
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
				respObj, _ := requestAndReturnArray(jsonStr, SearchByString)
				firstAuthor := respObj[1].Name //Use index 1 because in index 0 there is an author talking extensively about Nietzsche
				want := "Friedrich Nietzsche"
				if firstAuthor != want {
					t.Errorf("got %q, want %q", firstAuthor, want)
				}
			})

			t.Run("hard search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {

				var jsonStr = []byte(`{"searchString": "Nietshe Friedr"}`)
				respObj, _ := requestAndReturnArray(jsonStr, SearchByString)
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
				respObj, _ := requestAndReturnArray(jsonStr, SearchByString)
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

		respObj, _ := requestAndReturnArray(jsonStr, GetAuthorsById)
		firstAuthor := respObj[0]
		if firstAuthor.Id != authorId[0] {
			t.Errorf("got %d, want %d", firstAuthor.Authorid, authorId[0])
		}

	})

	t.Run("Authorlist Test", func(t *testing.T) {

		t.Run("Should return first 50 authors (alphabetically)", func(t *testing.T) {

			pageSize := 50
			var jsonStr = []byte(fmt.Sprintf(`{"pageSize": %d}`, pageSize))

			respObj, _ := requestAndReturnArray(jsonStr, GetAuthorsList)

			if len(respObj) != 50 {
				t.Errorf("got list of length %d, but expected list of length %d", len(respObj), pageSize)
			}

			firstAuthor := respObj[0]
			if firstAuthor.Name[0] != 'A' {
				t.Errorf("got %s, want name that starts with 'A'", firstAuthor.Name)
			}

		})

		t.Run("Should return first authors, with only English quotes, (alphabetically)", func(t *testing.T) {

			language := "english"
			var jsonStr = []byte(fmt.Sprintf(`{"language": "%s"}`, language))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Hasicelandicquotes {
				t.Errorf("got %+v, but expected an author that has no icelandic quotes", firstAuthor)
			}

			if firstAuthor.Name[0] != 'A' {
				t.Errorf("got %s, want name that starts with 'A'", firstAuthor.Name)
			}

		})

		t.Run("Should return first English authors in reverse alphabetical order (i.e. first author starts with Z)", func(t *testing.T) {

			language := "english"
			var jsonStr = []byte(fmt.Sprintf(`{"language": "%s", "orderConfig":{"orderBy":"alphabetical", "reverse":true}}`, language))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Hasicelandicquotes {
				t.Errorf("got %+v, but expected an author that has no icelandic quotes", firstAuthor)
			}

			if firstAuthor.Name[0] != 'Z' {
				t.Errorf("got %s, want name that starts with 'Z'", firstAuthor.Name)
			}

		})

		t.Run("Should return first authors starting from 'F' (i.e. greater than or equal to 'F' alphabetically)", func(t *testing.T) {
			language := "english"
			minimum := "f"
			var jsonStr = []byte(fmt.Sprintf(`{"language": "%s", "orderConfig":{"orderBy":"alphabetical","minimum":"%s"}}`, language, minimum))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Hasicelandicquotes {
				t.Errorf("got %+v, but expected an author that has no icelandic quotes", firstAuthor)
			}

			if firstAuthor.Name[0] != strings.ToUpper(minimum)[0] {
				t.Errorf("got %s, want name that starts with 'F'", firstAuthor.Name)
			}

		})

		t.Run("Should return first authors starting from 'Y' in reverse order (i.e. first authors gotten should start with Z and the last will end with Y)", func(t *testing.T) { t.Skip() })

		t.Run("Should return first authors starting from 'F' and Ending at (including) 'H' in reverse order, i.e. start at H and end at F", func(t *testing.T) { t.Skip() })

		t.Run("Should return authors with less than or equal to 1 quotes in total", func(t *testing.T) {

			maximum := 1
			var jsonStr = []byte(fmt.Sprintf(`{"orderConfig":{"orderBy":"nrOfQuotes","maximum":"%d"}}`, maximum))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Nroficelandicquotes+firstAuthor.Nrofenglishquotes > 1 {
				t.Errorf("got %+v, but expected an author that has no more than 1 quotes", firstAuthor)
			}

		})

		t.Run("Should return first authors with more than 10 quotes but less than or equal to 11 in total", func(t *testing.T) {

			minimum := 10
			maximum := 11
			var jsonStr = []byte(fmt.Sprintf(`{"orderConfig":{"orderBy":"nrOfQuotes","maximum":"%d", "minimum":"%d"}}`, maximum, minimum))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Nroficelandicquotes+firstAuthor.Nrofenglishquotes != 10 {
				t.Errorf("got %+v, but expected an author that has no fewer than 10 quotes", firstAuthor)
			}

		})

		t.Run("Should return first authors with less than 10 quotes in total in reversed order (start with those with 10 quotes)", func(t *testing.T) {

			maximum := 10
			var jsonStr = []byte(fmt.Sprintf(`{"orderConfig":{"orderBy":"nrOfQuotes","maximum":"%d","reverse":true}}`, maximum))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Nroficelandicquotes+firstAuthor.Nrofenglishquotes != 10 {
				t.Errorf("got %+v, but expected an author that has 10 quotes", firstAuthor)
			}

		})

		t.Run("Should return first authors (reverse order DESC by nr of quotes) only icelandic quotes", func(t *testing.T) {
			language := "icelandic"
			var jsonStr = []byte(fmt.Sprintf(`{"language":"%s", "orderConfig":{"orderBy":"nrOfQuotes","reverse":true}}`, language))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Nroficelandicquotes != 160 {
				t.Errorf("got %+v, but expected an author that has 10 quotes", firstAuthor)
			}
		})

		t.Run("Should return first 50 authors (ordered by most popular, i.e. DESC count)", func(t *testing.T) {
			directFetchAuthorsCountIncrement([]int{1})

			var jsonStr = []byte(fmt.Sprintf(`{"orderConfig":{"orderBy":"%s"}}`, "popularity"))

			respObj, errResponse := requestAndReturnArrayAuthors(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Count == 0 {
				t.Errorf("got %+v, but expected an author that does not have 0 popularity count", firstAuthor)
			}

		})

		t.Run("Should return first 50 authors in reverse popularity order (i.e. least popular first i.e. ASC count)", func(t *testing.T) {

			var jsonStr = []byte(fmt.Sprintf(`{"orderConfig":{"orderBy":"%s","reverse":true}}`, "popularity"))

			respObj, errResponse := requestAndReturnArrayAuthors(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Count != 0 {
				t.Errorf("got %+v, but expected an author that has 0 popularity count", firstAuthor)
			}

		})

		t.Run("Should return first 100 authors", func(t *testing.T) {
			pageSize := 100
			var jsonStr = []byte(fmt.Sprintf(`{"pageSize":%d}}`, pageSize))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			if len(respObj) != 100 {
				t.Errorf("got %d nr of authors, but expected %d authors", len(respObj), pageSize)
			}
		})

		t.Run("Should return the next 50 authors starting from 'F' (i.e. pagination, page 1, alphabetical order)", func(t *testing.T) {

			pageSize := 100
			minimum := "F"
			var jsonStr = []byte(fmt.Sprintf(`{"pageSize":%d, "orderConfig":{"minimum":"%s"}}}`, pageSize, minimum))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			objToFetch := respObj[50]

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			if respObj[0].Name[0] != minimum[0] {
				t.Errorf("got %+v, but expected author starting with '%s'", len(respObj), minimum)
			}

			pageSize = 50
			page := 1
			jsonStr = []byte(fmt.Sprintf(`{"pageSize":%d, "page":%d, "orderConfig":{"minimum":"%s"}}}`, pageSize, page, minimum))

			respObj, errResponse = requestAndReturnArray(jsonStr, GetAuthorsList)

			if objToFetch != respObj[0] {
				t.Fatalf("got %+v, but expected %+v", respObj[0], objToFetch)
			}

		})

	})

	t.Run("Random author", func(t *testing.T) {
		t.Run("Should return a random author with only a single quote (i.e. default)", func(t *testing.T) {

			var jsonStr = []byte(`{}`)
			firstRespObj, _ := requestAndReturnArray(jsonStr, GetRandomAuthor)

			if len(firstRespObj) != 1 {
				t.Errorf("Expected only a single quote from the random author but got %d", len(firstRespObj))
			}

			firstAuthor := firstRespObj[0]
			if firstAuthor.Name == "" {
				t.Errorf("Expected a random author but got an empty name for author")
			}

			secondRespObj, _ := requestAndReturnArray(jsonStr, GetRandomAuthor)
			secondAuthor := secondRespObj[0]
			if firstAuthor.Authorid == secondAuthor.Authorid {
				t.Errorf("Expected two different authors but got the same author twice which is higly improbable, got author with id %d and name %s", firstAuthor.Authorid, firstAuthor.Name)
			}

		})

		t.Run("Should return a random Author with only quotes from him in Icelandic", func(t *testing.T) {
			language := "icelandic"
			var jsonStr = []byte(fmt.Sprintf(`{"language":"%s"}`, language))
			firstRespObj, _ := requestAndReturnArray(jsonStr, GetRandomAuthor)

			firstAuthor := firstRespObj[0]
			if firstAuthor.Name == "" {
				t.Errorf("Expected a random author but got an empty name for author")
			}

			if !firstAuthor.Isicelandic {
				t.Errorf("Expected the quotes returned to be in icelandic")
			}

			secondRespObj, _ := requestAndReturnArray(jsonStr, GetRandomAuthor)
			secondAuthor := secondRespObj[0]
			if firstAuthor.Authorid == secondAuthor.Authorid {
				t.Errorf("Expected two different authors but got the same author twice which is higly improbable, got author with id %d and name %s", firstAuthor.Authorid, firstAuthor.Name)
			}
		})

		t.Run("Should return a random Author with only quotes from him in English", func(t *testing.T) {

			language := "english"
			var jsonStr = []byte(fmt.Sprintf(`{"language":"%s"}`, language))
			firstRespObj, _ := requestAndReturnArray(jsonStr, GetRandomAuthor)

			firstAuthor := firstRespObj[0]
			if firstAuthor.Name == "" {
				t.Errorf("Expected a random author but got an empty name for author")
			}

			if firstAuthor.Isicelandic {
				t.Errorf("Expected the quotes returned to be in English")
			}

			secondRespObj, _ := requestAndReturnArray(jsonStr, GetRandomAuthor)
			secondAuthor := secondRespObj[0]
			if firstAuthor.Authorid == secondAuthor.Authorid {
				t.Errorf("Expected two different authors but got the same author twice which is higly improbable, got author with id %d and name %s", firstAuthor.Authorid, firstAuthor.Name)
			}

		})

		t.Run("Should return author with a maximum of 2 of his quotes", func(t *testing.T) {
			maxQuotes := 2
			var jsonStr = []byte(fmt.Sprintf(`{"maxQuotes":%d}`, maxQuotes))
			firstRespObj, _ := requestAndReturnArray(jsonStr, GetRandomAuthor)

			firstAuthor := firstRespObj[0]
			if firstAuthor.Name == "" {
				t.Errorf("Expected a random author but got an empty name for author")
			}

			if len(firstRespObj) != 2 {
				t.Errorf("Expected 2 quotes but got %d", len(firstRespObj))
			}
		})

	})

	t.Run("Author of the day", func(t *testing.T) {

		t.Run("Should set / Overwrite Author of the day", func(t *testing.T) { t.Skip() })

		t.Run("Should get Author of the day", func(t *testing.T) { t.Skip() })

		t.Run("Should get complete history of Author of the day", func(t *testing.T) { t.Skip() })

		t.Run("Should get history of AOD starting from June 4th 2021", func(t *testing.T) { t.Skip() })

	})

}

func TestQuotes(t *testing.T) {
	t.Run("should return Quotes with id 1, 2 and 3...", func(t *testing.T) {

		var quoteIds = Set{1, 2, 3}
		var jsonStr = []byte(fmt.Sprintf(`{"ids":  [%s]}`, quoteIds.toString()))
		respObj, _ := requestAndReturnArray(jsonStr, GetQuotesById)

		if len(respObj) != len(quoteIds) {
			t.Errorf("got list of length %d but expected list of length %d", len(respObj), len(quoteIds))
		}

		for idx, quote := range respObj {
			if quote.Quoteid != quoteIds[idx] {
				t.Errorf("got %d, expected %d", quote.Quoteid, quoteIds[idx])
			}
		}
	})

	t.Run("Quote of the day", func(t *testing.T) {

		t.Run("Should set / Overwrite Quote of the day", func(t *testing.T) { t.Skip() })

		t.Run("Should get Quote of the day", func(t *testing.T) { t.Skip() })

		t.Run("Should get complete history of quote of the day", func(t *testing.T) { t.Skip() })

		t.Run("Should get history of QOD starting from June 4th 2021", func(t *testing.T) { t.Skip() })

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
			var jsonStr = []byte(fmt.Sprintf(`{"authorId": %d}`, teddyAuthor.Id))
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

func TestIncrementCount(t *testing.T) {
	t.Run("Should Increment Authors count from direct fetch by ids", func(t *testing.T) {
		authorIds := Set{1, 2, 3}
		err := directFetchAuthorsCountIncrement(authorIds)

		if err != nil {
			t.Fatalf("Expected no error but got %s", err.Error())
		}

		authors := getAuthorsById(authorIds)
		if authors[0].Count == 0 {
			t.Fatalf("Expected count of authors given to increase to above 0 but got count 0 for author: %+v", authors[0])
		}
	})

	t.Run("Should Increment Authors count from appearing in a search", func(t *testing.T) {

		quotes := []QuoteView{
			{
				Authorid: 100,
				Quoteid:  100,
			},
		}
		err := appearInSearchCountIncrement(quotes)

		if err != nil {
			t.Fatalf("Expected no error but got %s", err.Error())
		}

		authors := getAuthorsById([]int{quotes[0].Authorid})
		if authors[0].Count == 0 {
			t.Fatalf("Expected count of authors given to increase to above 0 but got count 0 for author: %+v", authors[0])
		}

	})
	t.Run("Should Increment Quotes count from direct fetch by ids", func(t *testing.T) { t.Skip() })
	t.Run("Should Increment Quotes count from appearing in a search", func(t *testing.T) { t.Skip() })
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

func getAuthor(searchString string) AuthorsView {

	var jsonStr = []byte(fmt.Sprintf(`{"searchString": "%s"}`, searchString))
	request, _ := http.NewRequest(http.MethodGet, "/api/search/authors", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	SearchAuthorsByString(response, request)

	var respObj []AuthorsView
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj[0]
}

func getAuthorsById(authorIds Set) []AuthorsView {

	var jsonStr = []byte(fmt.Sprintf(`{"ids": [%s]}`, authorIds.toString()))
	request, _ := http.NewRequest(http.MethodGet, "/api/search/authors", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	GetAuthorsById(response, request)

	var respObj []AuthorsView
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj
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

func requestAndReturnArray(jsonStr []byte, fn httpRequest) ([]QuoteView, ErrorResponse) {
	response, request := getRequestAndResponseForTest(jsonStr)
	fn(response, request)
	var respObj []QuoteView
	var errorResp ErrorResponse
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	_ = json.Unmarshal(response.Body.Bytes(), &errorResp)
	errorResp.StatusCode = response.Result().StatusCode
	return respObj, errorResp
}

func requestAndReturnArrayAuthors(jsonStr []byte, fn httpRequest) ([]AuthorsView, ErrorResponse) {
	response, request := getRequestAndResponseForTest(jsonStr)
	fn(response, request)
	var respObj []AuthorsView
	var errorResp ErrorResponse
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	_ = json.Unmarshal(response.Body.Bytes(), &errorResp)
	errorResp.StatusCode = response.Result().StatusCode
	return respObj, errorResp
}
