package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

func TestEasySearchQuotesByString(t *testing.T) {
	t.Run("should Return list of quotes with Muhammad Ali as first author", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/search/quotes", nil)
		response := httptest.NewRecorder()

		request = mux.SetURLVars(request, map[string]string{
			"searchString": "Float like a butterfly sting like a bee",
		})

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

func TestIntermediateSearchQuotesByString(t *testing.T) {
	t.Run("should Return list of quotes with Muhammad Ali as first author", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/search/quotes", nil)
		response := httptest.NewRecorder()

		request = mux.SetURLVars(request, map[string]string{
			"searchString": "bee sting like a butterfly",
		})

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

func TestHardSearchQuotesByString(t *testing.T) {
	t.Run("should Return list of quotes with Muhammad Ali as first author", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/search/quotes", nil)
		response := httptest.NewRecorder()

		request = mux.SetURLVars(request, map[string]string{
			"searchString": "bee butterfly float",
		})

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

func TestEasySearchAuthorsByString(t *testing.T) {
	t.Run("should Return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/search/authors", nil)
		response := httptest.NewRecorder()

		request = mux.SetURLVars(request, map[string]string{
			"searchString": "Friedrich Nietzsche",
		})

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

func TestIntermediateSearchAuthorsByString(t *testing.T) {
	t.Run("should Return list of quotes with Joseph Stalin as first author", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/search/authors", nil)
		response := httptest.NewRecorder()

		request = mux.SetURLVars(request, map[string]string{
			"searchString": "Stalin jseph",
		})

		SearchAuthorsByString(response, request)

		var respObj []SearchView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)
		firstAuthor := respObj[0].Name
		want := "Joseph Stalin"
		if firstAuthor != want {
			t.Errorf("got %q, want %q", firstAuthor, want)
		}
	})
}

func TestHardSearchAuthorsByString(t *testing.T) {
	t.Run("should Return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/search/authors", nil)
		response := httptest.NewRecorder()

		request = mux.SetURLVars(request, map[string]string{
			"searchString": "Niet Friedric",
		})

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

func TestEasySearchByStringForAuthor(t *testing.T) {
	t.Run("should Return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {
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
}

func TestHardSearchByStringForAuthor(t *testing.T) {
	t.Run("should Return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {
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
}

func TestEasySearchByStringForQuote(t *testing.T) {
	t.Run("should Return list of quotes with Martin Luther as first author", func(t *testing.T) {
		var jsonStr = []byte(`{"searchString": "If you are not allowed to Laugh in Heaven"}`)
		request, _ := http.NewRequest(http.MethodGet, "/api/search", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()

		SearchByString(response, request)

		var respObj []SearchView
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)
		log.Println(respObj)
		firstAuthor := respObj[1].Name //Use index 1 because in index 0 there is an author talking extensively about Nietzsche
		want := "Martin Luther"
		if firstAuthor != want {
			t.Errorf("got %q, want %q", firstAuthor, want)
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
	request, _ := http.NewRequest(http.MethodGet, "/api/search/authors", nil)
	response := httptest.NewRecorder()

	request = mux.SetURLVars(request, map[string]string{
		"searchString": searchString,
	})

	SearchAuthorsByString(response, request)

	var respObj []SearchView
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj[0]
}

func getQuotes(searchString string) []SearchView {
	request, _ := http.NewRequest(http.MethodGet, "/api/search/quotes", nil)
	response := httptest.NewRecorder()

	request = mux.SetURLVars(request, map[string]string{
		"searchString": searchString,
	})

	SearchQuotesByString(response, request)

	var respObj []SearchView
	_ = json.Unmarshal(response.Body.Bytes(), &respObj)
	return respObj
}
