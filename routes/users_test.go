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

func TestUsers(t *testing.T) {
	t.Run("Create User", func(t *testing.T) {
		t.Run("Should create user with name: Þórður Ágústsson, password: 1234567890 and email: random@gmail.com", func(t *testing.T) {

			password := "1234567890"
			passwordConfirmation := "1234567890"
			random, _ := uuid.NewRandom()
			name := "Þórður Ágústsson"
			email := random.String() + "@gmail.com"
			var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "passwordConfirmation":"%s", "email":"%s"}`, name, password, passwordConfirmation, email))
			userResponse, response := basicRequestReturnSingle(jsonStr, CreateUser)

			if response.Result().StatusCode != 200 {
				t.Fatalf("Expected a status code of 200 but got %d", response.Result().StatusCode)
			}

			if userResponse.ApiKey == "" {
				t.Fatalf("Expected an Api key but got %s", userResponse.ApiKey)
			}
		})

		t.Run("Should try to create user with email already in use but get error", func(t *testing.T) {
			password := "1234567890"
			passwordConfirmation := "1234567890"
			random, _ := uuid.NewRandom()
			name := "Þórður Ágústsson"
			email := random.String() + "@gmail.com"
			var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "passwordConfirmation":"%s", "email":"%s"}`, name, password, passwordConfirmation, email))

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
			password := "1234567890"
			passwordConfirmation := "12345678901"
			random, _ := uuid.NewRandom()
			name := "Þórður Ágústsson"
			email := random.String() + "@gmail.com"
			var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "passwordConfirmation":"%s", "email":"%s"}`, name, password, passwordConfirmation, email))
			_, response := basicRequestReturnSingle(jsonStr, CreateUser)

			if response.Result().StatusCode != 400 {
				t.Fatalf("Expected a status code of 200 but got %d", response.Result().StatusCode)
			}
		})

		t.Run("Should try to create user with weak password, and get error", func(t *testing.T) {
			password := "123"
			passwordConfirmation := "123"
			random, _ := uuid.NewRandom()
			name := "Þórður Ágústsson"
			email := random.String() + "@gmail.com"
			var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "passwordConfirmation":"%s", "email":"%s"}`, name, password, passwordConfirmation, email))
			_, response := basicRequestReturnSingle(jsonStr, CreateUser)

			if response.Result().StatusCode != 400 {
				t.Fatalf("Expected a status code of 200 but got %d", response.Result().StatusCode)
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
