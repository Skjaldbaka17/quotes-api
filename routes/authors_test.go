package routes

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/Skjaldbaka17/quotes-api/handlers"
	"github.com/Skjaldbaka17/quotes-api/structs"
	"github.com/google/uuid"
)

func createUser(t *testing.T) structs.UserResponse {
	password := "1234567890"
	passwordConfirmation := "1234567890"
	random, _ := uuid.NewRandom()
	name := "Þórður Ágústsson"
	email := random.String() + "@gmail.com"
	var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "passwordConfirmation":"%s", "email":"%s"}`, name, password, passwordConfirmation, email))
	userResponse, response := basicRequestReturnSingle(jsonStr, CreateUser)
	if response.Result().StatusCode != 200 {
		t.Fatalf("Failed creating user, got statust code %d but expected %d", response.Result().StatusCode, 200)
	}

	t.Cleanup(func() {
		log.Println("CLEANUP:", userResponse)
		handlers.Db.Table("users").Delete(&structs.User{Id: userResponse.Id})
		handlers.Db.Table("requesthistory").Where("user_id = ?", userResponse.Id).Delete(structs.RequestEvent{})
	})
	return userResponse
}

func getGODModeUser(t *testing.T) structs.UserResponse {
	var user structs.UserResponse
	handlers.Db.Table("users").Where("tier = 'GOD'").First(&user)
	t.Cleanup(func() {
		log.Println("CLEANUPGOD:", user)
		handlers.Db.Table("requesthistory").Where("user_id = ?", user.Id).Delete(structs.RequestEvent{})
	})
	return user
}

func TestAuthors(t *testing.T) {
	user := createUser(t)
	godUser := getGODModeUser(t)

	t.Run("Get authors", func(t *testing.T) {
		t.Run("should return Author with id 1", func(t *testing.T) {

			authorId := Set{1}
			var jsonStr = []byte(fmt.Sprintf(`{"ids": [%s], "apiKey":"%s"}`, authorId.toString(), user.ApiKey))

			respObj, _ := requestAndReturnArray(jsonStr, GetAuthorsById)
			firstAuthor := respObj[0]
			if firstAuthor.Id != authorId[0] {
				t.Fatalf("got %d, want %d", firstAuthor.Authorid, authorId[0])
			}
		})

		t.Run("Should return 400 error, request body not structured correctly, the ids in body is array of strings", func(t *testing.T) {
			authorId := []string{"1", "2", "3"}
			var jsonStr = []byte(fmt.Sprintf(`{"ids": [%s], "apiKey":"%s"}`, authorId, user.ApiKey))

			_, response := requestAndReturnArray(jsonStr, GetAuthorsById)
			if response.StatusCode != 400 {
				t.Fatalf("got statusCode %d, want 400", response.StatusCode)
			}
		})

	})

	t.Run("Authorlist Test", func(t *testing.T) {

		t.Run("Should return first 50 authors (alphabetically)", func(t *testing.T) {

			pageSize := 50
			var jsonStr = []byte(fmt.Sprintf(`{"pageSize": %d,"apiKey":"%s"}`, pageSize, user.ApiKey))

			respObj, _ := requestAndReturnArray(jsonStr, GetAuthorsList)

			if len(respObj) != 50 {
				t.Fatalf("got list of length %d, but expected list of length %d", len(respObj), pageSize)
			}

			firstAuthor := respObj[0]
			if firstAuthor.Name[0] != 'A' {
				t.Fatalf("got %s, want name that starts with 'A'", firstAuthor.Name)
			}

		})

		t.Run("Should return first authors, with only English quotes, (alphabetically)", func(t *testing.T) {

			language := "english"
			var jsonStr = []byte(fmt.Sprintf(`{"language": "%s","apiKey":"%s"}`, language, user.ApiKey))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Hasicelandicquotes {
				t.Fatalf("got %+v, but expected an author that has no icelandic quotes", firstAuthor)
			}

			if firstAuthor.Name[0] != 'A' {
				t.Fatalf("got %s, want name that starts with 'A'", firstAuthor.Name)
			}

		})

		t.Run("Should return first English authors in reverse alphabetical order (i.e. first author starts with Z)", func(t *testing.T) {

			language := "english"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language": "%s", "orderConfig":{"orderBy":"alphabetical", "reverse":true}}`, user.ApiKey, language))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Hasicelandicquotes {
				t.Fatalf("got %+v, but expected an author that has no icelandic quotes", firstAuthor)
			}

			if firstAuthor.Name[0] != 'Z' {
				t.Fatalf("got %s, want name that starts with 'Z'", firstAuthor.Name)
			}

		})

		t.Run("Should return first authors starting from 'F' (i.e. greater than or equal to 'F' alphabetically)", func(t *testing.T) {
			language := "english"
			minimum := "f"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language": "%s", "orderConfig":{"orderBy":"alphabetical","minimum":"%s"}}`, user.ApiKey, language, minimum))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Hasicelandicquotes {
				t.Fatalf("got %+v, but expected an author that has no icelandic quotes", firstAuthor)
			}

			if firstAuthor.Name[0] != strings.ToUpper(minimum)[0] {
				t.Fatalf("got %s, want name that starts with 'F'", firstAuthor.Name)
			}

		})

		t.Run("Should return first authors starting from 'Y' in reverse order (i.e. first authors gotten should start with Z and the last will end with Y)", func(t *testing.T) { t.Skip() })

		t.Run("Should return first authors starting from 'F' and Ending at (including) 'H' in reverse order, i.e. start at H and end at F", func(t *testing.T) { t.Skip() })

		t.Run("Should return authors with less than or equal to 1 quotes in total", func(t *testing.T) {

			maximum := 1
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s", "orderConfig":{"orderBy":"nrOfQuotes","maximum":"%d"}}`, user.ApiKey, maximum))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Nroficelandicquotes+firstAuthor.Nrofenglishquotes > 1 {
				t.Fatalf("got %+v, but expected an author that has no more than 1 quotes", firstAuthor)
			}

		})

		t.Run("Should return first authors with more than 10 quotes but less than or equal to 11 in total", func(t *testing.T) {

			minimum := 10
			maximum := 11
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","orderConfig":{"orderBy":"nrOfQuotes","maximum":"%d", "minimum":"%d"}}`, user.ApiKey, maximum, minimum))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Nroficelandicquotes+firstAuthor.Nrofenglishquotes != 10 {
				t.Fatalf("got %+v, but expected an author that has no fewer than 10 quotes", firstAuthor)
			}

		})

		t.Run("Should return first authors with less than 10 quotes in total in reversed order (start with those with 10 quotes)", func(t *testing.T) {

			maximum := 10
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","orderConfig":{"orderBy":"nrOfQuotes","maximum":"%d","reverse":true}}`, user.ApiKey, maximum))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Nroficelandicquotes+firstAuthor.Nrofenglishquotes != 10 {
				t.Fatalf("got %+v, but expected an author that has 10 quotes", firstAuthor)
			}

		})

		t.Run("Should return first authors (reverse order DESC by nr of quotes) only icelandic quotes", func(t *testing.T) {
			language := "icelandic"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language":"%s", "orderConfig":{"orderBy":"nrOfQuotes","reverse":true}}`, user.ApiKey, language))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Nroficelandicquotes != 160 {
				t.Fatalf("got %+v, but expected an author that has 10 quotes", firstAuthor)
			}
		})

		t.Run("Should return first 50 authors (ordered by most popular, i.e. DESC count)", func(t *testing.T) {
			handlers.DirectFetchAuthorsCountIncrement([]int{1})

			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","orderConfig":{"orderBy":"%s"}}`, user.ApiKey, "popularity"))

			respObj, errResponse := requestAndReturnArrayAuthors(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Count == 0 {
				t.Fatalf("got %+v, but expected an author that does not have 0 popularity count", firstAuthor)
			}

		})

		t.Run("Should return first 50 authors in reverse popularity order (i.e. least popular first i.e. ASC count)", func(t *testing.T) {

			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","orderConfig":{"orderBy":"%s","reverse":true}}`, user.ApiKey, "popularity"))

			respObj, errResponse := requestAndReturnArrayAuthors(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			firstAuthor := respObj[0]

			if firstAuthor.Count != 0 {
				t.Fatalf("got %+v, but expected an author that has 0 popularity count", firstAuthor)
			}

		})

		t.Run("Should return first 100 authors", func(t *testing.T) {
			pageSize := 100
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","pageSize":%d}}`, user.ApiKey, pageSize))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			if len(respObj) != 100 {
				t.Fatalf("got %d nr of authors, but expected %d authors", len(respObj), pageSize)
			}
		})

		t.Run("Should return the next 50 authors starting from 'F' (i.e. pagination, page 1, alphabetical order)", func(t *testing.T) {

			pageSize := 100
			minimum := "F"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","pageSize":%d, "orderConfig":{"minimum":"%s"}}}`, user.ApiKey, pageSize, minimum))

			respObj, errResponse := requestAndReturnArray(jsonStr, GetAuthorsList)

			objToFetch := respObj[50]

			if errResponse.StatusCode != 200 {
				t.Fatalf("got error %s, but expected an empty errormessage", errResponse.Message)
			}

			if respObj[0].Name[0] != minimum[0] {
				t.Fatalf("got %+v, but expected author starting with '%s'", len(respObj), minimum)
			}

			pageSize = 50
			page := 1
			jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s", "pageSize":%d, "page":%d, "orderConfig":{"minimum":"%s"}}}`, user.ApiKey, pageSize, page, minimum))

			respObj, errResponse = requestAndReturnArray(jsonStr, GetAuthorsList)

			if objToFetch != respObj[0] {
				t.Fatalf("got %+v, but expected %+v", respObj[0], objToFetch)
			}

		})

	})

	t.Run("Random author", func(t *testing.T) {
		t.Run("Should return a random author with only a single quote (i.e. default)", func(t *testing.T) {

			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s"}`, user.ApiKey))
			firstRespObj, _ := requestAndReturnArray(jsonStr, GetRandomAuthor)

			if len(firstRespObj) != 1 {
				t.Fatalf("Expected only a single quote from the random author but got %d", len(firstRespObj))
			}

			firstAuthor := firstRespObj[0]
			if firstAuthor.Name == "" {
				t.Fatalf("Expected a random author but got an empty name for author")
			}

			secondRespObj, _ := requestAndReturnArray(jsonStr, GetRandomAuthor)
			secondAuthor := secondRespObj[0]
			if firstAuthor.Authorid == secondAuthor.Authorid {
				t.Fatalf("Expected two different authors but got the same author twice which is higly improbable, got author with id %d and name %s", firstAuthor.Authorid, firstAuthor.Name)
			}

		})

		t.Run("Should return a random Author with only quotes from him in Icelandic", func(t *testing.T) {
			language := "icelandic"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language":"%s"}`, user.ApiKey, language))
			firstRespObj, _ := requestAndReturnArray(jsonStr, GetRandomAuthor)

			firstAuthor := firstRespObj[0]
			if firstAuthor.Name == "" {
				t.Fatalf("Expected a random author but got an empty name for author")
			}

			if !firstAuthor.Isicelandic {
				t.Fatalf("Expected the quotes returned to be in icelandic")
			}

			secondRespObj, _ := requestAndReturnArray(jsonStr, GetRandomAuthor)
			secondAuthor := secondRespObj[0]
			if firstAuthor.Authorid == secondAuthor.Authorid {
				t.Fatalf("Expected two different authors but got the same author twice which is higly improbable, got author with id %d and name %s", firstAuthor.Authorid, firstAuthor.Name)
			}
		})

		t.Run("Should return a random Author with only quotes from him in English", func(t *testing.T) {

			language := "english"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language":"%s"}`, user.ApiKey, language))
			firstRespObj, _ := requestAndReturnArray(jsonStr, GetRandomAuthor)

			firstAuthor := firstRespObj[0]
			if firstAuthor.Name == "" {
				t.Fatalf("Expected a random author but got an empty name for author")
			}

			if firstAuthor.Isicelandic {
				t.Fatalf("Expected the quotes returned to be in English")
			}

			secondRespObj, _ := requestAndReturnArray(jsonStr, GetRandomAuthor)
			secondAuthor := secondRespObj[0]
			if firstAuthor.Authorid == secondAuthor.Authorid {
				t.Fatalf("Expected two different authors but got the same author twice which is higly improbable, got author with id %d and name %s", firstAuthor.Authorid, firstAuthor.Name)
			}

		})

		t.Run("Should return author with a maximum of 2 of his quotes", func(t *testing.T) {
			maxQuotes := 2
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","maxQuotes":%d}`, user.ApiKey, maxQuotes))
			firstRespObj, _ := requestAndReturnArray(jsonStr, GetRandomAuthor)

			firstAuthor := firstRespObj[0]
			if firstAuthor.Name == "" {
				t.Fatalf("Expected a random author but got an empty name for author")
			}

			if len(firstRespObj) != 2 {
				t.Fatalf("Expected 2 quotes but got %d", len(firstRespObj))
			}
		})

	})

	t.Run("Author of the day", func(t *testing.T) {

		t.Run("Should set / Overwrite Author of the day", func(t *testing.T) {

			authorId := 1
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s", "aods": [{"id":%d, "date":""}]}`, godUser.ApiKey, authorId))
			_, response := requestAndReturnArray(jsonStr, SetAuthorOfTheDay)
			if response.StatusCode != 200 {
				t.Fatalf("Expected a succesful insert but got %+v", response)
			}
		})

		t.Run("Should set AOD for 12-22-2020 and 12-21-2020", func(t *testing.T) {

			//TODO: add to test that the quotes where actually input into the DB
			authorId1 := 2
			date1 := "2020-12-22"
			date2 := "2020-12-21"
			authorId2 := 3
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s", "aods": [{"id":%d, "date":"%s"},{"id":%d, "date":"%s"}]}`, godUser.ApiKey, authorId1, date1, authorId2, date2))
			_, response := requestAndReturnArray(jsonStr, SetAuthorOfTheDay)
			if response.StatusCode != 200 {
				t.Fatalf("Expected a succesful insert but got %+v", response)
			}

		})

		t.Run("Should get Author of the day", func(t *testing.T) {

			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s", "language":"%s"}`, user.ApiKey, "english"))
			author := requestAndReturnSingle(jsonStr, GetAuthorOfTheDay)

			if author.Name == "" {
				t.Fatalf("Expected the author of the day but got an empty author %+v", author)
			}

			if author.Id == 0 {
				t.Fatalf("Expected the autho to have id > 0 but got: %+v", author)
			}

			const layout = "2006-01-02T15:04:05Z" //The date needed for reference always
			date, _ := time.Parse(layout, author.Date)
			if date.Format("01-02-2006") != time.Now().Format("01-02-2006") {
				t.Fatalf("Expected the author for the date %s but got AOD for date %s i.e. %+v", time.Now().Format("01-02-2006"), date.Format("01-02-2006"), author)
			}

		})

		t.Run("Should get complete history of Author of the day", func(t *testing.T) {

			//Input a quote in history for testing
			authorId := 1111
			date := "1998-06-16"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s", "aods": [{"id":%d, "date":"%s"}]}`, godUser.ApiKey, authorId, date))
			_, response := requestAndReturnArray(jsonStr, SetAuthorOfTheDay)
			if response.StatusCode != 200 {
				t.Fatalf("Expected a succesful insert but got %+v", response)
			}

			//Get History:

			jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language":"%s"}`, user.ApiKey, "english"))
			authors, _ := requestAndReturnArray(jsonStr, GetAODHistory)

			if len(authors) == 0 {
				t.Fatalf("Expected the history of AOD but got an empty list: %+v", authors)
			}

			containsBirfdayAuthor := false
			containsTodayAuthor := false
			const layout = "2006-01-02T15:04:05Z" //The date needed for reference always
			for _, author := range authors {
				if author.Id == 0 {
					t.Fatalf("Expected all authors to have id > 0 but got: %+v", authors)
				}
				date, _ := time.Parse(layout, author.Date)
				if date.Format("01-02-2006") == time.Now().Format("01-02-2006") {
					containsTodayAuthor = true
				}

				if date.Format("01-02-2006") == "06-16-1998" {
					containsBirfdayAuthor = true
				}
			}

			if !containsBirfdayAuthor {
				t.Fatalf("AOD history should contain the AOD for birfday but does not: %+v", authors)
			}

			if !containsTodayAuthor {
				t.Fatalf("AOD history should contain the AOD for today but does not: %+v", authors)
			}

		})

		t.Run("Should get history of AOD starting from June 4th 2021", func(t *testing.T) {

			//Input a quote in history for testing
			authorId := 666
			date := "2021-06-04"
			var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","aods": [{"id":%d, "date":"%s"}]}`, godUser.ApiKey, authorId, date))
			_, response := requestAndReturnArray(jsonStr, SetAuthorOfTheDay)
			if response.StatusCode != 200 {
				t.Fatalf("Expected a succesful insert but got %+v", response)
			}

			//Get History:

			minimum := "2021-06-04"
			jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language":"%s", "minimum":"%s"}`, user.ApiKey, "english", minimum))
			authors, _ := requestAndReturnArray(jsonStr, GetAODHistory)

			if len(authors) == 0 {
				t.Fatalf("Expected the history of AOD but got an empty list: %+v", authors)
			}

			const layout = "2006-01-02T15:04:05Z" //The date needed for reference always
			compareDate, _ := time.Parse(layout, "2021-06-04")
			compareYear := compareDate.Year()
			compareMonth := compareDate.Month()
			compareDay := compareDate.Day()
			containsAuthorNotInRange := false
			containsFourthOfJuneAuthor := false
			for _, author := range authors {
				date, _ := time.Parse(layout, author.Date)
				yearOfAuthor := date.Year()
				monthOfAuthor := date.Month()
				dayOfAuthor := date.Day()

				if yearOfAuthor < compareYear || (yearOfAuthor == compareYear && monthOfAuthor < compareMonth) || (yearOfAuthor == compareYear && monthOfAuthor == compareMonth && dayOfAuthor < compareDay) {
					containsAuthorNotInRange = true
				}

				if date.Format("2006-01-02") == "2021-06-04" {
					containsFourthOfJuneAuthor = true
				}

				if author.Id == 0 {
					t.Fatalf("Expected all authors to have id > 0 but got: %+v", authors)
				}

			}

			if containsAuthorNotInRange {
				t.Fatalf("AOD history contains an earlier quote than was requested: %+v", authors)
			}

			if !containsFourthOfJuneAuthor {
				t.Fatalf("QOD history should contain the QOD for 4th of june 2021 but does not: %+v", authors)
			}

		})

	})

	t.Cleanup(func() {
		log.Println("CLEANUP TestAuthors!")
		// Delete from aod
		handlers.Db.Exec("DELETE FROM aod")
		handlers.Db.Exec("DELETE FROM aodice")
		// Set popularity of authors to 0
		handlers.Db.Exec("Update authors set count = 0 where count > 0")
		// Set popularity of quotes to 0
		handlers.Db.Exec("Update quotes set count = 0 where count > 0")
		// Set popularity of topics to 0
		handlers.Db.Exec("Update topics set count = 0 where count > 0")
	})
}
