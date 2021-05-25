package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestSearchQuotesByString(t *testing.T) {
	t.Run("should Return list of quotes with Muhammad Ali as first author", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/search/quotes", nil)
		response := httptest.NewRecorder()

		request = mux.SetURLVars(request, map[string]string{
			"searchString": "butterfly sting bee",
		})

		SearchQuotesByString(response, request)

		var respObj []SearchView
		got := json.Unmarshal(response.Body.Bytes(), &respObj)
		firstAuthor := respObj[0].Name
		want := "Muhammad Ali"
		log.Println("responÞ", got)
		if firstAuthor != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
