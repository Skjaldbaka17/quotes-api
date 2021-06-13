package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Skjaldbaka17/quotes-api/handlers"
	"github.com/Skjaldbaka17/quotes-api/routes"
	"github.com/Skjaldbaka17/quotes-api/structs"
	"github.com/google/uuid"
)

type httpRequest func(http.ResponseWriter, *http.Request)

var functions map[string]httpRequest = map[string]httpRequest{
	"GetAuthorsById":        routes.GetAuthorsById,
	"GetAuthorsList":        routes.GetAuthorsList,
	"GetRandomAuthor":       routes.GetRandomAuthor,
	"GetAuthorOfTheDay":     routes.GetAuthorOfTheDay,
	"GetAODHistory":         routes.GetAODHistory,
	"SetAuthorOfTheDay":     routes.SetAuthorOfTheDay,
	"GetQuotes":             routes.GetQuotes,
	"GetQuotesList":         routes.GetQuotesList,
	"GetRandomQuote":        routes.GetRandomQuote,
	"SetQuoteOfTheDay":      routes.SetQuoteOfTheDay,
	"GetQuoteOfTheDay":      routes.GetQuoteOfTheDay,
	"GetQODHistory":         routes.GetQODHistory,
	"SearchByString":        routes.SearchByString,
	"SearchAuthorsByString": routes.SearchAuthorsByString,
	"SearchQuotesByString":  routes.SearchQuotesByString,
	"GetTopics":             routes.GetTopics,
	"GetTopic":              routes.GetTopic,
}

func TestUtils(t *testing.T) {
	t.Run("Get Request Body Test", func(t *testing.T) {
		t.Run("Running all protected functions without ApiKey, should return error for all", func(t *testing.T) {
			var jsonStr = []byte(`{}`)
			for name, fn := range functions {
				response := basicRequestReturnResponse(jsonStr, fn)
				if response.Result().StatusCode != http.StatusForbidden {
					t.Fatalf("Expected to be declined and receive status code %d, but got %d when callin the route/function %s", http.StatusForbidden, response.Result().StatusCode, name)
				}
			}

		})
		t.Run("Running all protected functions with non-existant ApiKey, should return error for all", func(t *testing.T) {
			var jsonStr = []byte(`{"apiKey":"MyApiKeyB"}`)
			for name, fn := range functions {
				response := basicRequestReturnResponse(jsonStr, fn)
				if response.Result().StatusCode != http.StatusForbidden {
					t.Fatalf("Expected to be declined and receive status code %d, but got %d when callin the route/function %s", http.StatusForbidden, response.Result().StatusCode, name)
				}
			}
		})

		t.Run("Running all protected functions with existant ApiKey but all requests-per hour for this user are used, should return error for all", func(t *testing.T) {
			user := getBasicUser()
			var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "passwordConfirmation":"%s", "email":"%s"}`, user.Name, user.Password, user.PasswordConfirmation, user.Email))
			userResponse, _ := basicRequestReturnSingle(jsonStr, routes.CreateUser)

			if userResponse.ApiKey == "" {
				t.Fatalf("Expected an Api key but got %s", userResponse.ApiKey)
			}

			var userStruct structs.User
			handlers.Db.Table("users").Where("api_key = ?", userResponse.ApiKey).First(&userStruct)

			//Creating events to insert, use up all allowed requests for this apikey
			requestEvents := []structs.RequestEvent{}
			for i := 0; i < handlers.REQUESTS_PER_HOUR[handlers.TIERS[0]]; i++ {
				requestEvents = append(requestEvents, structs.RequestEvent{
					UserId:      userStruct.Id,
					ApiKey:      userResponse.ApiKey,
					RequestBody: `{}`,
					Route:       "/",
				})
			}

			//Inserting the created events into the history
			err := handlers.Db.Table("requesthistory").Create(&requestEvents).Error
			if err != nil {
				t.Fatalf("Failed inserting the request events into the table, %s", err)
			}

			jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s"}`, userResponse.ApiKey))
			for name, fn := range functions {
				response := basicRequestReturnResponse(jsonStr, fn)
				if response.Result().StatusCode != http.StatusUnauthorized {
					t.Fatalf("Expected to be declined and receive status code %d, but got %d when callin the route/function %s", http.StatusUnauthorized, response.Result().StatusCode, name)
				}
			}
		})

	})
}

func basicRequestReturnResponse(jsonStr []byte, fn httpRequest) *httptest.ResponseRecorder {
	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	fn(response, request)
	return response
}

func basicRequestReturnSingle(jsonStr []byte, fn httpRequest) (structs.UserResponse, *httptest.ResponseRecorder) {
	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	fn(response, request)
	var userResponse structs.UserResponse
	_ = json.Unmarshal(response.Body.Bytes(), &userResponse)
	return userResponse, response
}

func getBasicUser() structs.UserRequest {
	password := "1234567890"
	passwordConfirmation := "1234567890"
	random, _ := uuid.NewRandom()
	name := "Þórður Ágústsson"
	email := random.String() + "@gmail.com"
	return structs.UserRequest{
		Name:                 name,
		Email:                email,
		Password:             password,
		PasswordConfirmation: passwordConfirmation,
	}
}
