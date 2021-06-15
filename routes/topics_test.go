package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Skjaldbaka17/quotes-api/structs"
)

func TestTopics(t *testing.T) {
	user := createUser(t)

	t.Run("Should return the possible English topics as a list of objects", func(t *testing.T) {

		var nrOfEnglishTopics int = 13
		var language string = "English"
		var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language": "%s"}`, user.ApiKey, language))
		request, _ := http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()

		GetTopics(response, request)

		var respObj []structs.TestApiResponse
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)

		if len(respObj) != nrOfEnglishTopics {
			t.Fatalf("got %d number of topics, expected %d", len(respObj), nrOfEnglishTopics)
		}
	})

	t.Run("Should return the possible Icelandic topics as a list of objects", func(t *testing.T) {

		var nrOfIcelandicTopics int = 7
		var language string = "Icelandic"
		var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language": "%s"}`, user.ApiKey, language))
		request, _ := http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()
		GetTopics(response, request)
		var respObj []structs.TestApiResponse
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)

		if len(respObj) != nrOfIcelandicTopics {
			t.Fatalf("got %d number of topics, expected %d", len(respObj), nrOfIcelandicTopics)
		}
	})

	t.Run("Should return the first 25 quotes from a topic 'nameOfTopic'", func(t *testing.T) {

		var nameOfTopic string = "inspirational"
		var pageSize int = 25
		var page int = 0
		var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s", "topic": "%s", "pageSize":%d, "page":%d}`, user.ApiKey, nameOfTopic, pageSize, page))
		request, _ := http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()
		GetTopic(response, request)
		var respObj []structs.TestApiResponse
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)

		if len(respObj) != pageSize {
			t.Fatalf("got %d number of quotes, expected %d as the pagesize", len(respObj), pageSize)
		}

		for _, obj := range respObj {
			if respObj[0].TopicName != nameOfTopic {
				t.Fatalf("got %+v but expected a quote with topic %s", obj, nameOfTopic)
			}
		}

	})

	t.Run("Should return the first 25 quotes from a topic with id", func(t *testing.T) {

		topicId := getTopicId("inspirational", user.ApiKey)
		var pageSize int = 26
		var page int = 0
		var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","id": %d, "pageSize":%d, "page":%d}`, user.ApiKey, topicId, pageSize, page))
		request, _ := http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()
		GetTopic(response, request)
		var respObj []structs.TestApiResponse
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)

		if len(respObj) != pageSize {
			t.Fatalf("got %d number of quotes, expected %d as the pagesize", len(respObj), pageSize)
		}

		for _, obj := range respObj {
			if respObj[0].TopicId != topicId {
				t.Fatalf("got %+v but expected a quote with topicId %d", obj, topicId)
			}
		}

	})

	t.Run("Should test pagination for a specific topic, by id", func(t *testing.T) {

		topicId := getTopicId("inspirational", user.ApiKey)
		//First get the 2nd page's first quote, where pagesize is 25
		var pageSize int = 25
		var page int = 1
		var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","id": %d, "pageSize":%d, "page":%d}`, user.ApiKey, topicId, pageSize, page))
		request, _ := http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()
		GetTopic(response, request)

		var respObj []structs.TestApiResponse
		_ = json.Unmarshal(response.Body.Bytes(), &respObj)

		obj26 := respObj[0]

		// Then get the first 100 quotes, i.e. first page with pagesize 100
		pageSize = 100
		page = 0
		jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s", "id": %d, "pageSize":%d, "page":%d}`, user.ApiKey, topicId, pageSize, page))
		request, _ = http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(jsonStr))
		response = httptest.NewRecorder()
		GetTopic(response, request)

		var newRespObj []structs.TestApiResponse
		_ = json.Unmarshal(response.Body.Bytes(), &newRespObj)

		if len(newRespObj) != pageSize {
			t.Fatalf("got %d number of quotes, expected %d as the pagesize", len(newRespObj), pageSize)
		}

		//Compare the 26th object from the 100pagesize request with the 1st object from the 2nd page where pagesize is 25.
		if newRespObj[25] != obj26 {
			t.Fatalf("got %+v but expected %+v", newRespObj[25], obj26)
		}

	})

	t.Cleanup(func() {
		log.Println("CLEANUP TestTopics!")
	})
}
