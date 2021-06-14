package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/Skjaldbaka17/quotes-api/structs"
)

const defaultPageSize = 25
const maxPageSize = 200
const maxQuotes = 50
const defaultMaxQuotes = 1

//returns error and the body as a string
func getBody(rw http.ResponseWriter, r *http.Request, requestBody *structs.Request) (error, string) {
	buf, _ := ioutil.ReadAll(r.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	//Save the state back into the body for later use (Especially useful for getting the AOD/QOD because if the AOD has not been set a random AOD is set and the function called again)
	r.Body = rdr2
	if err := json.NewDecoder(rdr1).Decode(&requestBody); err != nil {
		log.Printf("Got error when decoding: %s", err)
		err = errors.New("request body is not structured correctly. Please refer to the /docs page for information on how to structure the request body")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
		return err, ""
	}
	return nil, string(buf)
}

// ValidateRequestApiKey checks if the ApiKey supplied exists and wether the user has finished his allowed request in the past
// hour. Also adds to the requestHistory... Maybe move that to the end of a request?
func validateRequestApiKey(rw http.ResponseWriter, r *http.Request) error {
	var requestBody structs.Request
	err, bodyAsString := getBody(rw, r, &requestBody)
	if err != nil {
		return err
	}

	if requestBody.ApiKey == "" {
		log.Printf("no ApiKey given when accessing resource")
		err := errors.New("you need to supply an apiKey to access this resource. Create a user and get a free-tier apiKey here: https://www.example.com")
		rw.WriteHeader(http.StatusForbidden)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
		return err
	}

	var user structs.UserDBModel

	err = Db.Table("users").Where("api_key = ?", requestBody.ApiKey).First(&user).Error
	// Err==nil if user with given api_key does not exist or internal server error
	if err != nil {
		m1 := regexp.MustCompile(`record not found`)
		if m1.Match([]byte(err.Error())) {
			log.Printf("the api-key that the requester supplied does not exist")
			err := errors.New("you need a valid ApiKey to access this resource. Create a user and get a free-tier ApiKey here: https://www.example.com")
			rw.WriteHeader(http.StatusForbidden)
			json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
			return err
		}
		log.Printf("error when searching for user with the given api key (api key validation): %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: InternalServerError})
		return err
	}

	//Check if requests from this api-key the past hour are less than allowed for the users-tier (i.e. if this next request is
	// allowed then save the request to request-history)
	type countStruct struct {
		Count int `json:"count"`
	}
	var count countStruct
	if err := Db.Table("requesthistory").Select("count(*)").
		Where("created_at >= (NOW() - INTERVAL '1 hour')").
		Where("user_id = ?", user.Id).
		First(&count).Error; err != nil {
		log.Printf("error when counting request history: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: InternalServerError})
		return err
	}

	if float64(count.Count) >= REQUESTS_PER_HOUR[user.Tier] {
		err := fmt.Errorf(
			"you have used all the requests per hour that your tier %s allows for, i.e. %f request per hour. See https://www.example.com for more info and pricing plans to upgrade your tier if necessary", user.Tier, REQUESTS_PER_HOUR[user.Tier])
		rw.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
		return err
	}

	//TODO: Put the following in its own golang function and run as a separate process!
	requestEvent := structs.RequestEvent{
		UserId:      user.Id,
		RequestBody: bodyAsString,
		Route:       r.URL.String(),
		ApiKey:      user.ApiKey,
	}
	result := Db.Table("requesthistory").Create(&requestEvent)
	if result.Error != nil {
		log.Printf("error when inserting into requestHistory: %s", result.Error)
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: InternalServerError})
		return err
	}

	return nil
}

//ValidateRequestBody takes in the request and validates all the input fields, returns an error with reason for validation-failure
//if validation fails.
//TODO: Make validation better! i.e. make it "real"
func GetRequestBody(rw http.ResponseWriter, r *http.Request, requestBody *structs.Request) error {
	if err := validateRequestApiKey(rw, r); err != nil {
		return err
	}
	if err, _ := getBody(rw, r, requestBody); err != nil {
		return err
	}

	if requestBody.PageSize < 1 || requestBody.PageSize > maxPageSize {
		requestBody.PageSize = defaultPageSize
	}

	if requestBody.Page < 0 {
		requestBody.Page = 0
	}

	if requestBody.MaxQuotes < 0 || requestBody.MaxQuotes > maxQuotes {
		requestBody.MaxQuotes = maxQuotes
	}

	if requestBody.MaxQuotes <= 0 {
		requestBody.MaxQuotes = defaultMaxQuotes
	}

	const layout = "2006-01-02"
	//Set date into correct format, if supplied, otherwise input today's date in the correct format for all qods
	if len(requestBody.Qods) != 0 {
		for idx, _ := range requestBody.Qods {
			if requestBody.Qods[idx].Date == "" {
				requestBody.Qods[idx].Date = time.Now().UTC().Format(layout)
			} else {
				var parsedDate time.Time
				parsedDate, err := time.Parse(layout, requestBody.Qods[idx].Date)
				if err != nil {
					log.Printf("Got error when decoding: %s", err)
					err = fmt.Errorf("the date is not structured correctly, should be in %s format", layout)
					rw.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
					return err
				}

				requestBody.Qods[idx].Date = parsedDate.UTC().Format(layout)
			}
		}
	}

	//Set date into correct format, if supplied, otherwise input today's date in the correct format for all qods
	if len(requestBody.Aods) != 0 {
		for idx, _ := range requestBody.Aods {
			if requestBody.Aods[idx].Date == "" {
				requestBody.Aods[idx].Date = time.Now().UTC().Format(layout)
			} else {
				var parsedDate time.Time
				parsedDate, err := time.Parse(layout, requestBody.Aods[idx].Date)
				if err != nil {
					log.Printf("Got error when decoding: %s", err)
					err = fmt.Errorf("the date is not structured correctly, should be in %s format", layout)
					rw.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
					return err
				}

				requestBody.Aods[idx].Date = parsedDate.UTC().Format(layout)
			}
		}
	}

	if requestBody.Minimum != "" {

		_, err := time.Parse(layout, requestBody.Minimum)
		if err != nil {
			log.Printf("Got error when decoding: %s", err)
			err = fmt.Errorf("the minimum date is not structured correctly, should be in %s format", layout)
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
			return err
		}
		// requestBody.Minimum = parseDate.Format("01-02-2006")
	}

	if requestBody.Maximum != "" {

		parseDate, err := time.Parse(layout, requestBody.Maximum)
		if err != nil {
			log.Printf("Got error when decoding: %s", err)
			err = fmt.Errorf("the maximum date is not structured correctly, should be in %s format", layout)
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
			return err
		}
		requestBody.Minimum = parseDate.Format("01-02-2006")
	}

	return nil
}

//ValidateUserRequestBody takes in the request and validates all the input fields, returns an error with reason for validation-failure
//if validation fails.
//TODO: Make validation better! i.e. make it "real"
func GetUserRequestBody(rw http.ResponseWriter, r *http.Request, requestBody *structs.UserApiModel) error {
	//Save the state back into the body for later use (Especially useful for getting the AOD/QOD because if the AOD has not been set a random AOD is set and the function called again)
	buf, _ := ioutil.ReadAll(r.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	//Save the state back into the body for later use (Especially useful for getting the AOD/QOD because if the AOD has not been set a random AOD is set and the function called again)
	r.Body = rdr2
	err := json.NewDecoder(rdr1).Decode(&requestBody)

	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		err = errors.New("request body is not structured correctly. Please refer to the /docs page for information on how to structure the request body")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
		return err
	}

	return nil
}

func ValidateUserInformation(rw http.ResponseWriter, r *http.Request, requestBody *structs.UserApiModel) error {
	//TODO: Add email validation
	if requestBody.Email == "" {
		err := errors.New("email should not be empty")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
		return err
	}

	if requestBody.Name == "" {
		err := errors.New("name should not be empty")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
		return err
	}

	if requestBody.Password == "" {
		err := errors.New("password should not be empty")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
		return err
	}

	if len(requestBody.Password) < 8 {
		err := errors.New("password should be at least 8 characters long")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
		return err
	}

	if requestBody.PasswordConfirmation == "" {
		err := errors.New("password confirmation should not be empty")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
		return err
	}

	if requestBody.PasswordConfirmation != requestBody.Password {
		err := errors.New("passwords do not match")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
		return err
	}
	return nil
}

// Check whether user has GOD-tier permissions
func AuthorizeGODApiKey(rw http.ResponseWriter, r *http.Request) error {
	var requestBody structs.Request
	if err, _ := getBody(rw, r, &requestBody); err != nil {
		return err
	}

	var user structs.UserDBModel
	if err := Db.Table("users").Where("api_key = ?", requestBody.ApiKey).First(&user).Error; err != nil {
		log.Printf("error when searching for user with the given api key in AuthorIzeGOD (api key validation): %s", err)
		rw.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "You need special privileges to access this route."})
		return err
	}

	if user.Tier != TIERS[len(TIERS)-1] {
		err := errors.New("you do not have the authorization to perform this action. Is your name Bassi Maraj? This is not meant for you... Sorry for the inconvenience")
		rw.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
		return err
	}

	return nil
}
