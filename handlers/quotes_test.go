package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		got := json.Unmarshal(response.Body.Bytes(), &respObj)
		firstAuthor := respObj[0].Name
		want := "Muhammad Ali"
		if firstAuthor != want {
			t.Errorf("got %q, want %q", got, want)
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
		got := json.Unmarshal(response.Body.Bytes(), &respObj)
		firstAuthor := respObj[0].Name
		want := "Muhammad Ali"
		if firstAuthor != want {
			t.Errorf("got %q, want %q", got, want)
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
		got := json.Unmarshal(response.Body.Bytes(), &respObj)
		firstAuthor := respObj[0].Name
		want := "Muhammad Ali"
		if firstAuthor != want {
			t.Errorf("got %q, want %q", got, want)
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
		got := json.Unmarshal(response.Body.Bytes(), &respObj)
		firstAuthor := respObj[0].Name
		want := "Friedrich Nietzsche"
		if firstAuthor != want {
			t.Errorf("got %q, want %q", got, want)
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
		got := json.Unmarshal(response.Body.Bytes(), &respObj)
		firstAuthor := respObj[0].Name
		want := "Joseph Stalin"
		if firstAuthor != want {
			t.Errorf("got %q, want %q", got, want)
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
		got := json.Unmarshal(response.Body.Bytes(), &respObj)
		firstAuthor := respObj[0].Name
		want := "Friedrich Nietzsche"
		if firstAuthor != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
