package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Skjaldbaka17/quotes-api/structs"
)

func TestSearch(t *testing.T) {
	t.Run("Search Quotes By String", func(t *testing.T) {
		t.Run("easy search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "Float like a butterfly sting like a bee"}`)
			respObj, _ := requestAndReturnArray(jsonStr, SearchQuotesByString)
			firstAuthor := respObj[0].Name
			want := "Muhammad Ali"
			if firstAuthor != want {
				t.Fatalf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("intermediate search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "bee sting like a butterfly"}`)
			respObj, _ := requestAndReturnArray(jsonStr, SearchQuotesByString)
			firstAuthor := respObj[0].Name
			want := "Muhammad Ali"
			if firstAuthor != want {
				t.Fatalf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("hard search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "bee butterfly float"}`)
			respObj, _ := requestAndReturnArray(jsonStr, SearchQuotesByString)
			firstAuthor := respObj[0].Name
			want := "Muhammad Ali"
			if firstAuthor != want {
				t.Fatalf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("Search for quote 'Happiness resides not in possessions...' inside topic 'inspirational' by supplying its topicid", func(t *testing.T) {
			topicId := getTopicId("inspirational")
			var jsonStr = []byte(fmt.Sprintf(`{"searchString": "Happiness resides not in possessions", "topicId":%d}`, topicId))
			respObj, _ := requestAndReturnArray(jsonStr, SearchQuotesByString)
			firstAuthorName := respObj[0].Name
			want_author := "Democritus"
			if firstAuthorName != want_author {
				t.Fatalf("got %q, want %q", firstAuthorName, want_author)
			}

			firstAuthorQuote := respObj[0].Quote
			want_quote := "Happiness resides not in possessions, and not in gold, happiness dwells in the soul."
			if firstAuthorQuote != want_quote {
				t.Fatalf("got %q, want %q", firstAuthorQuote, want_quote)
			}

			if respObj[0].Topicid != topicId {
				t.Fatalf("got quote with topicId %d, but expected with topicID %d. Quote got: %+v", respObj[0].Topicid, topicId, respObj[0])
			}
		})
	})

	t.Run("Search Authors By String", func(t *testing.T) {
		//Michael Jordan
		t.Run("easy search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "Friedrich Nietzsche"}`)
			respObj, _ := requestAndReturnArray(jsonStr, SearchAuthorsByString)
			firstAuthor := respObj[0].Name
			want := "Friedrich Nietzsche"
			if firstAuthor != want {
				t.Fatalf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("intermediate search should Return list of quotes with Joseph Stalin as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "Stalin jseph"}`)
			respObj, _ := requestAndReturnArray(jsonStr, SearchAuthorsByString)
			firstAuthor := respObj[0].Name
			want := "Joseph Stalin"
			if firstAuthor != want {
				t.Fatalf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("hard search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {

			var jsonStr = []byte(`{"searchString": "Niet Friedric"}`)
			respObj, _ := requestAndReturnArray(jsonStr, SearchAuthorsByString)
			firstAuthor := respObj[0].Name
			want := "Friedrich Nietzsche"
			if firstAuthor != want {
				t.Fatalf("got %q, want %q", firstAuthor, want)
			}
		})

		//Dont think this application is necessary
		t.Run("Search Authors inside topic 'inspirational' for 'Michael Jordan' by supplying its topicid", func(t *testing.T) { t.Skip() })
	})

	t.Run("Search by String", func(t *testing.T) {

		t.Run("searching for author", func(t *testing.T) {
			t.Run("easy search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {

				var jsonStr = []byte(`{"searchString": "Friedrich Nietzsche"}`)
				respObj, _ := requestAndReturnArray(jsonStr, SearchByString)
				firstAuthor := respObj[1].Name //Use index 1 because in index 0 there is an author talking extensively about Nietzsche
				want := "Friedrich Nietzsche"
				if firstAuthor != want {
					t.Fatalf("got %q, want %q", firstAuthor, want)
				}
			})

			t.Run("hard search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {

				var jsonStr = []byte(`{"searchString": "Nietshe Friedr"}`)
				respObj, _ := requestAndReturnArray(jsonStr, SearchByString)
				firstAuthor := respObj[1].Name //Use index 1 because in index 0 there is an author talking extensively about Nietzsche
				want := "Friedrich Nietzsche"
				if firstAuthor != want {
					t.Fatalf("got %q, want %q", firstAuthor, want)
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
					t.Fatalf("got %q, want %q", firstAuthor, want)
				}
			})
		})

		t.Run("General Search inside topic 'inspirational' by supplying its id, should return 'Michael Jordan' Quote", func(t *testing.T) {
			topicId := getTopicId("inspirational")
			var jsonStr = []byte(fmt.Sprintf(`{"searchString": "Jordan Michel", "topicId":%d}`, topicId))
			respObj, _ := requestAndReturnArray(jsonStr, SearchByString)
			firstAuthorName := respObj[0].Name
			want_author := "Michael Jordan"
			if firstAuthorName != want_author {
				t.Fatalf("got %q, want %q", firstAuthorName, want_author)
			}

			if respObj[0].Topicid != topicId {
				t.Fatalf("got quote with topicId %d, but expected with topicID %d. Quote got: %+v", respObj[0].Topicid, topicId, respObj[0])
			}
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
			var respObj []structs.QuoteView
			_ = json.Unmarshal(response.Body.Bytes(), &respObj)

			if pageSize != len(respObj) {
				t.Fatalf("got list of length %d but expected %d", len(respObj), pageSize)
			}

			if respObj[0].Quoteid != obj26.Quoteid {
				t.Fatalf("got %+v, want %+v", respObj[0], obj26)
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

			var respObj []structs.QuoteView
			_ = json.Unmarshal(response.Body.Bytes(), &respObj)

			if pageSize != len(respObj) {
				t.Fatalf("got list of length %d but expected %d", len(respObj), pageSize)
			}

			if respObj[0] != obj26 {
				t.Fatalf("got %+v, want %+v", respObj[0], obj26)
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

			var respObj []structs.QuoteView
			_ = json.Unmarshal(response.Body.Bytes(), &respObj)

			if pageSize != len(respObj) {
				t.Fatalf("got list of length %d but expected %d", len(respObj), pageSize)
			}

			if respObj[0].Quoteid != obj26.Quoteid {
				t.Fatalf("got %+v, want %+v", respObj[0], obj26)
			}
		})

	})

}