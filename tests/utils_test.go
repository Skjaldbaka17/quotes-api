package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
	"GetQuotes":             routes.GetQuotes,
	"GetQuotesList":         routes.GetQuotesList,
	"GetRandomQuote":        routes.GetRandomQuote,
	"GetQuoteOfTheDay":      routes.GetQuoteOfTheDay,
	"GetQODHistory":         routes.GetQODHistory,
	"SearchByString":        routes.SearchByString,
	"SearchAuthorsByString": routes.SearchAuthorsByString,
	"SearchQuotesByString":  routes.SearchQuotesByString,
	"GetTopics":             routes.GetTopics,
	"GetTopic":              routes.GetTopic,
}

func deleteUser(id int) {
	handlers.Db.Table("users").Delete(&structs.User{Id: id})
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

			defer deleteUser(userResponse.Id)
			defer func() {
				// Delete from request History
				handlers.Db.Exec("delete from requesthistory where user_id = ?", userResponse.Id)
			}()

			var userStruct structs.User
			handlers.Db.Table("users").Where("api_key = ?", userResponse.ApiKey).First(&userStruct)

			//Creating events to insert, use up all allowed requests for this apikey
			requestEvents := []structs.RequestEvent{}
			for i := 0.0; i < handlers.REQUESTS_PER_HOUR[handlers.TIERS[0]]; i++ {
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

		t.Run("Should save a requestEvent in requestHistory for all routes", func(t *testing.T) { t.Skip() })
		t.Run("Should save a errorEvent in requestHistory for all routes, making errors happen", func(t *testing.T) { t.Skip() })

		t.Run("Should not be able to setQOD nor setAOD", func(t *testing.T) {
			user := getBasicUser()
			var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "passwordConfirmation":"%s", "email":"%s"}`, user.Name, user.Password, user.PasswordConfirmation, user.Email))
			userResponse, _ := basicRequestReturnSingle(jsonStr, routes.CreateUser)
			defer deleteUser(userResponse.Id)
			jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s", "aods": [{"id":%d, "date":""}]}`, userResponse.ApiKey, 1))

			response := basicRequestReturnResponse(jsonStr, routes.SetAuthorOfTheDay)
			if response.Result().StatusCode != http.StatusUnauthorized {
				t.Fatalf("Expected to be declined, unauthorized user trying to setAOD, and receive status code %d, but got %d ", http.StatusUnauthorized, response.Result().StatusCode)
			}

			jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s", "qods": [{"id":%d, "date":""}]}`, userResponse.ApiKey, 1))

			response = basicRequestReturnResponse(jsonStr, routes.SetQuoteOfTheDay)
			if response.Result().StatusCode != http.StatusUnauthorized {
				t.Fatalf("Expected to be declined, unauthorized user trying to setQOD, and receive status code %d, but got %d ", http.StatusUnauthorized, response.Result().StatusCode)
			}

		})
	})

	t.Cleanup(func() {
		log.Println("CLEANUP TestUtils!")
		// Set popularity of authors to 0
		handlers.Db.Exec("Update authors set count = 0 where count > 0")
		// Set popularity of quotes to 0
		handlers.Db.Exec("Update quotes set count = 0 where count > 0")
		// Set popularity of topics to 0
		handlers.Db.Exec("Update topics set count = 0 where count > 0")
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
