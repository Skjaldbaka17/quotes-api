package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Skjaldbaka17/quotes-api/structs"
)

const defaultPageSize = 25
const maxPageSize = 200
const maxQuotes = 50
const defaultMaxQuotes = 1

//ValidateRequestBody takes in the request and validates all the input fields, returns an error with reason for validation-failure
//if validation fails.
//TODO: Make validation better! i.e. make it "real"
func GetRequestBody(rw http.ResponseWriter, r *http.Request, requestBody *structs.Request) error {
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		log.Printf("Got error when decoding: %s", err)
		err = errors.New("request body is not structured correctly. Please refer to the /docs page for information on how to structure the request body")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
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
				parsedDate, err = time.Parse(layout, requestBody.Qods[idx].Date)
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
				parsedDate, err = time.Parse(layout, requestBody.Aods[idx].Date)
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

		parseDate, err := time.Parse(layout, requestBody.Minimum)
		if err != nil {
			log.Printf("Got error when decoding: %s", err)
			err = fmt.Errorf("the minimum date is not structured correctly, should be in %s format", layout)
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
			return err
		}
		requestBody.Minimum = parseDate.Format("01-02-2006")
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
