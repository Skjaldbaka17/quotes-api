package routes

import (
	"encoding/json"
	"net/http"
)

var languages = []string{"English", "Icelandic"}

// swagger:route GET /languages META GetLanguages
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
