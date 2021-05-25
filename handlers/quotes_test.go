package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

func TestSearchQuotesByString(t *testing.T) {
	t.Run("easy search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {
		var jsonStr = []byte(`{"searchString": "Float like a butterfly sting like a bee"}`)
		request, _ := http.NewRequest(http.MethodPost, "/api/search/quotes", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()

		SearchQuotesByString(response, request)

		var respObj []SearchView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)
		firstAuthor := respObj[0].Name
		want := "Muhammad Ali"
		if firstAuthor != want {
			t.Errorf("got %q, want %q", firstAuthor, want)
		}
	})

	t.Run("intermediate search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {
		var jsonStr = []byte(`{"searchString": "bee sting like a butterfly"}`)
		request, _ := http.NewRequest(http.MethodPost, "/api/search/quotes", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()

		SearchQuotesByString(response, request)

		var respObj []SearchView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)
		firstAuthor := respObj[0].Name
		want := "Muhammad Ali"
		if firstAuthor != want {
			t.Errorf("got %q, want %q", firstAuthor, want)
		}
	})

	t.Run("hard search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {
		var jsonStr = []byte(`{"searchString": "bee butterfly float"}`)
		request, _ := http.NewRequest(http.MethodPost, "/api/search/quotes", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()

		SearchQuotesByString(response, request)

		var respObj []SearchView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)
		firstAuthor := respObj[0].Name
		want := "Muhammad Ali"
		if firstAuthor != want {
			t.Errorf("got %q, want %q", firstAuthor, want)
		}
	})
}

func TestSearchAuthorsByString(t *testing.T) {
	t.Run("easy search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {
		var jsonStr = []byte(`{"searchString": "Friedrich Nietzsche"}`)
		request, _ := http.NewRequest(http.MethodPost, "/api/search/authors", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()

		SearchAuthorsByString(response, request)

		var respObj []SearchView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)
		firstAuthor := respObj[0].Name
		want := "Friedrich Nietzsche"
		if firstAuthor != want {
			t.Errorf("got %q, want %q", firstAuthor, want)
		}
	})

	t.Run("intermediate search should Return list of quotes with Joseph Stalin as first author", func(t *testing.T) {
		var jsonStr = []byte(`{"searchString": "Stalin jseph"}`)
		request, _ := http.NewRequest(http.MethodPost, "/api/search/authors", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()

		SearchAuthorsByString(response, request)

		var respObj []SearchView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)
		firstAuthor := respObj[0].Name
		want := "Joseph Stalin"
		if firstAuthor != want {
			t.Errorf("got %q, want %q", firstAuthor, want)
		}
	})

	t.Run("hard search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {
		var jsonStr = []byte(`{"searchString": "Niet Friedric"}`)
		request, _ := http.NewRequest(http.MethodPost, "/api/search/authors", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()

		SearchAuthorsByString(response, request)

		var respObj []SearchView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)
		firstAuthor := respObj[0].Name
		want := "Friedrich Nietzsche"
		if firstAuthor != want {
			t.Errorf("got %q, want %q", firstAuthor, want)
		}
	})
}

func TestSearchByString(t *testing.T) {
	t.Run("searching for author", func(t *testing.T) {
		t.Run("easy search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {
			var jsonStr = []byte(`{"searchString": "Friedrich Nietzsche"}`)
			request, _ := http.NewRequest(http.MethodPost, "/api/search", bytes.NewBuffer(jsonStr))
			response := httptest.NewRecorder()

			SearchByString(response, request)

			var respObj []SearchView
			_ = json.Unmarshal(response.Body.Bytes(), &respObj)
			firstAuthor := respObj[1].Name //Use index 1 because in index 0 there is an author talking extensively about Nietzsche
			want := "Friedrich Nietzsche"
			if firstAuthor != want {
				t.Errorf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("hard search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {
			var jsonStr = []byte(`{"searchString": "Nietshe Friedr"}`)
			request, _ := http.NewRequest(http.MethodGet, "/api/search", bytes.NewBuffer(jsonStr))
			response := httptest.NewRecorder()

			SearchByString(response, request)

			var respObj []SearchView
			_ = json.Unmarshal(response.Body.Bytes(), &respObj)
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
			request, _ := http.NewRequest(http.MethodGet, "/api/search", bytes.NewBuffer(jsonStr))
			response := httptest.NewRecorder()

			SearchByString(response, request)

			var respObj []SearchView
			_ = json.Unmarshal(response.Body.Bytes(), &respObj)
			firstAuthor := respObj[1].Name //Use index 1 because in index 0 there is an author talking extensively about Nietzsche
			want := "Martin Luther"
			if firstAuthor != want {
				t.Errorf("got %q, want %q", firstAuthor, want)
			}
		})
	})
}

func TestPagination(t *testing.T) {
	t.Run("Search By string pagination", func(t *testing.T) {
		pageSize := 100
		var jsonStr = []byte(fmt.Sprintf(`{"searchString": "Love", "pageSize":%d}`, pageSize))
		request, _ := http.NewRequest(http.MethodPost, "/api/search", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()

		SearchByString(response, request)

		var respObj []SearchView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)

		if pageSize != len(respObj) {
			t.Errorf("got list of length %d but expected %d", len(respObj), pageSize)
		}
		obj26 := respObj[25]
		//Next request to check if same dude in position 0 given that pageSize is 25 and same search parameters
		pageSize = 25
		jsonStr = []byte(fmt.Sprintf(`{"searchString": "Love", "pageSize":%d, "page":1}`, pageSize))
		request, _ = http.NewRequest(http.MethodPost, "/api/search", bytes.NewBuffer(jsonStr))
		response = httptest.NewRecorder()

		SearchByString(response, request)
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)

		if pageSize != len(respObj) {
			t.Errorf("got list of length %d but expected %d", len(respObj), pageSize)
		}

		if respObj[0] != obj26 {
			t.Errorf("got %+v, want %+v", respObj[0], obj26)
		}
	})

	t.Run("Search Authors By string pagination", func(t *testing.T) {
		pageSize := 100
		var jsonStr = []byte(fmt.Sprintf(`{"searchString": "Martin", "pageSize":%d}`, pageSize))
		request, _ := http.NewRequest(http.MethodPost, "/api/search/authors", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()

		SearchAuthorsByString(response, request)

		var respObj []SearchView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)
		if pageSize != len(respObj) {
			t.Errorf("got list of length %d but expected %d", len(respObj), pageSize)
		}
		obj26 := respObj[25]
		//Next request to check if same dude in position 0 given that pageSize is 25 and same search parameters
		pageSize = 25
		jsonStr = []byte(fmt.Sprintf(`{"searchString": "Martin", "pageSize":%d, "page":1}`, pageSize))
		request, _ = http.NewRequest(http.MethodPost, "/api/search/authors", bytes.NewBuffer(jsonStr))
		response = httptest.NewRecorder()

		SearchAuthorsByString(response, request)
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)

		if pageSize != len(respObj) {
			t.Errorf("got list of length %d but expected %d", len(respObj), pageSize)
		}

		if respObj[0] != obj26 {
			t.Errorf("got %+v, want %+v", respObj[0], obj26)
		}
	})

	t.Run("Search Quotes By string pagination", func(t *testing.T) {
		pageSize := 100
		var jsonStr = []byte(fmt.Sprintf(`{"searchString": "Hate", "pageSize":%d}`, pageSize))
		request, _ := http.NewRequest(http.MethodPost, "/api/search/quotes", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()

		SearchQuotesByString(response, request)

		var respObj []SearchView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)
		if pageSize != len(respObj) {
			t.Errorf("got list of length %d but expected %d", len(respObj), pageSize)
		}
		obj26 := respObj[25]
		//Next request to check if same dude in position 0 given that pageSize is 25 and same search parameters
		pageSize = 25
		jsonStr = []byte(fmt.Sprintf(`{"searchString": "Hate", "pageSize":%d, "page":1}`, pageSize))
		request, _ = http.NewRequest(http.MethodPost, "/api/search/quotes", bytes.NewBuffer(jsonStr))
		response = httptest.NewRecorder()

		SearchQuotesByString(response, request)
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)

		if pageSize != len(respObj) {
			t.Errorf("got list of length %d but expected %d", len(respObj), pageSize)
		}

		if respObj[0] != obj26 {
			t.Errorf("got %+v, want %+v", respObj[0], obj26)
		}
	})
}

func TestGetAuthorById(t *testing.T) {
	t.Run("should return Muhammad Ali", func(t *testing.T) {
		//Need to get the authorid first, because the id can be different depending on DB.
		author := getAuthor("Muhammad Ali")

		request, _ := http.NewRequest(http.MethodGet, "/api/authors/", nil)
		response := httptest.NewRecorder()

		request = mux.SetURLVars(request, map[string]string{
			"id": strconv.Itoa(author.Authorid),
		})

		GetAuthorById(response, request)

		var respObj SearchView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)
		got := respObj.Name
		want := "Muhammad Ali"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestGetQuoteById(t *testing.T) {
	t.Run("should return muhammad alis quote float like a butterfly...", func(t *testing.T) {
		want := "Float like a butterfly, sting like a bee."
		quote := getQuotes(want)
		request, _ := http.NewRequest(http.MethodGet, "/api/authors/", nil)
		response := httptest.NewRecorder()

		request = mux.SetURLVars(request, map[string]string{
			"id": strconv.Itoa(quote[0].Quoteid),
		})

		GetQuoteById(response, request)

		var respObj SearchView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)
		got := respObj.Quote

		//Because Muhammad Ali has two quotes with "Float like a butterfly, sting like a bee.", the returned quotes
		// from getQuotes have in index 0 the quote where he is talking about the butterfly quote, and the real quote
		// is in index 1 therefore we test the got-quote with regexp
		m1 := regexp.MustCompile(want)
		if m1.MatchString(got) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func getAuthor(searchString string) SearchView {
	var jsonStr = []byte(fmt.Sprintf(`{"searchString": "%s"}`, searchString))
	request, _ := http.NewRequest(http.MethodGet, "/api/search/authors", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	SearchAuthorsByString(response, request)

	var respObj []SearchView
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj[0]
}

func getQuotes(searchString string) []SearchView {
	var jsonStr = []byte(fmt.Sprintf(`{"searchString": "%s"}`, searchString))
	request, _ := http.NewRequest(http.MethodGet, "/api/search/quotes", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	SearchQuotesByString(response, request)

	var respObj []SearchView
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj
}
