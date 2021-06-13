package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Skjaldbaka17/quotes-api/structs"
	"github.com/google/uuid"
)

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
func TestUsers(t *testing.T) {
	t.Run("Create User", func(t *testing.T) {
		t.Run("Should create user with name: Þórður Ágústsson, password: 1234567890 and email: random@gmail.com", func(t *testing.T) {
			user := getBasicUser()
			var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "passwordConfirmation":"%s", "email":"%s"}`, user.Name, user.Password, user.PasswordConfirmation, user.Email))
			userResponse, response := basicRequestReturnSingle(jsonStr, CreateUser)

			if response.Result().StatusCode != 200 {
				t.Fatalf("Expected a status code of 200 but got %d", response.Result().StatusCode)
			}

			if userResponse.ApiKey == "" {
				t.Fatalf("Expected an Api key but got %s", userResponse.ApiKey)
			}
		})

		t.Run("Should try to create user with email already in use but get error", func(t *testing.T) {
			user := getBasicUser()
			var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "passwordConfirmation":"%s", "email":"%s"}`, user.Name, user.Password, user.PasswordConfirmation, user.Email))

			_, _ = basicRequestReturnSingle(jsonStr, CreateUser)

			//Do it again with same info
			_, response := basicRequestReturnSingle(jsonStr, CreateUser)

			if response.Result().StatusCode != 400 {
				t.Fatalf("Expected a status code of 400 but got %d", response.Result().StatusCode)
			}
		})

		t.Run("Should try to create user with non-regex-email and get error", func(t *testing.T) {
			t.Skip()
		})

		t.Run("Should try to create user with password and confirm password not the same and get error", func(t *testing.T) {
			user := getBasicUser()
			user.PasswordConfirmation = "12345678901"
			var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "passwordConfirmation":"%s", "email":"%s"}`, user.Name, user.Password, user.PasswordConfirmation, user.Email))
			_, response := basicRequestReturnSingle(jsonStr, CreateUser)

			if response.Result().StatusCode != 400 {
				t.Fatalf("Expected a status code of 200 but got %d", response.Result().StatusCode)
			}
		})

		t.Run("Should try to create user with weak password, and get error", func(t *testing.T) {
			user := getBasicUser()
			user.PasswordConfirmation = "123"
			user.Password = "123"
			var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "passwordConfirmation":"%s", "email":"%s"}`, user.Name, user.Password, user.PasswordConfirmation, user.Email))
			_, response := basicRequestReturnSingle(jsonStr, CreateUser)

			if response.Result().StatusCode != 400 {
				t.Fatalf("Expected a status code of 200 but got %d", response.Result().StatusCode)
			}
		})
	})

	t.Run("Login", func(t *testing.T) {
		t.Run("Should login to account", func(t *testing.T) {
			//First create user

			user := getBasicUser()
			var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "passwordConfirmation":"%s", "email":"%s"}`, user.Name, user.Password, user.PasswordConfirmation, user.Email))
			firstUserResponse, _ := basicRequestReturnSingle(jsonStr, CreateUser)

			apiKey := firstUserResponse.ApiKey

			if apiKey == "" {
				t.Fatalf("Expected an Api key as response to creating the user but got %s", apiKey)
			}

			//Login to createdUser
			jsonStr = []byte(fmt.Sprintf(`{"password":"%s", "email":"%s"}`, user.Password, user.Email))
			userResponse, response := basicRequestReturnSingle(jsonStr, Login)

			if response.Result().StatusCode != 200 {
				t.Fatalf("Expected a status code of 200 but got %d", response.Result().StatusCode)
			}

			if userResponse.ApiKey == "" {
				t.Fatalf("Expected an Api key but got %s", userResponse.ApiKey)
			}

			if userResponse.ApiKey != apiKey {
				t.Fatalf("Expected the same api key, %s, as the created user has but got %s ", apiKey, userResponse.ApiKey)
			}

		})

		t.Run("Should try and login to an account with a wrong email and get error", func(t *testing.T) {
			//First create user
			user := getBasicUser()
			//Login to createdUser
			var jsonStr = []byte(fmt.Sprintf(`{"password":"%s", "email":"%s"}`, user.Password, user.Email))
			userResponse, response := basicRequestReturnSingle(jsonStr, Login)

			if userResponse.ApiKey != "" {
				t.Fatalf("Expected an empty Api Key but got %s!", userResponse.ApiKey)
			}
			if response.Result().StatusCode != 400 {
				t.Fatalf("Expected a status code of 400 but got %d", response.Result().StatusCode)
			}
		})

		t.Run("Should try and login to account, correct email, but with wrong password and get error", func(t *testing.T) {
			//First create user

			user := getBasicUser()
			var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "passwordConfirmation":"%s", "email":"%s"}`, user.Name, user.Password, user.PasswordConfirmation, user.Email))
			firstUserResponse, _ := basicRequestReturnSingle(jsonStr, CreateUser)

			apiKey := firstUserResponse.ApiKey

			if apiKey == "" {
				t.Fatalf("Expected an Api key as response to creating the user but got %s", apiKey)
			}

			//Login to createdUser
			jsonStr = []byte(fmt.Sprintf(`{"password":"wrooooooooong", "email":"%s"}`, user.Email))
			userResponse, response := basicRequestReturnSingle(jsonStr, Login)

			if userResponse.ApiKey != "" {
				t.Fatalf("Expected an empty Api Key but got %s!", userResponse.ApiKey)
			}

			if response.Result().StatusCode != 401 {
				t.Fatalf("Expected a status code of 401 but got %d", response.Result().StatusCode)
			}
		})
	})
}

func basicRequestReturnSingle(jsonStr []byte, fn httpRequest) (structs.UserResponse, *httptest.ResponseRecorder) {
	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()

	fn(response, request)
	var userResponse structs.UserResponse
	_ = json.Unmarshal(response.Body.Bytes(), &userResponse)
	return userResponse, response
}
