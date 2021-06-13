package docs

import "github.com/Skjaldbaka17/quotes-api/structs"

// Data structure representing the response for a authors
// swagger:response authorsResponse
type authorsResponseWrapper struct {
	// A usual authors response
	// in: body
	Body []structs.AuthorsView
}

// Data structure representing the response for a random author
// swagger:response randomAuthorResponse
type randomAuthorResponseWrapper struct {
	// A quote struct
	// in: body
	Body []baseQuotesResponseModel //model
}

// Data structure representing the response for the quotes
// swagger:response quotesResponse
type quotesResponseWrapper struct {
	// A quote struct
	// in: body
	Body []baseQuotesResponseModel //model
}

// Data structure representing the response for a quote
// swagger:response randomQuoteResponse
type quoteResponseWrapper struct {
	// A quote struct
	// in: body
	Body structs.QuoteView
}

// Data structure representing the response for the quote of the day
// swagger:response quoteOfTheDayResponse
type quoteOfTheDayResponseWrapper struct {
	// The response to the author of the day request
	// in: body
	Body qodResponseModel
}

// Data structure representing the response for the history of QODS
// swagger:response qodHistoryResponse
type qodHistoryResponseWrapper struct {
	// The response to the author of the day request
	// in: body
	Body []qodResponseModel
}

// Data structure representing the response for the author of the day
// swagger:response authorOfTheDayResponse
type authorOfTheDayResponseWrapper struct {
	// The response to the author of the day request
	// in: body
	Body struct {
		// The author's id
		//Unique: true
		//example: 24952
		Id int `json:"id"`
		// Name of the author
		//example: Muhammad Ali
		Name string `json:"name"`
		// The date when this author was the author of the day
		// example: 2021-06-12T00:00:00Z
		Date string `json:"date"`
	}
}

// Data structure representing the response for the history of AODs
// swagger:response aodHistoryResponse
type aodHistoryResponseWrapper struct {
	// The response to the author of the day request
	// in: body
	Body []struct {
		// The author's id
		//Unique: true
		//example: 24952
		Id int `json:"id"`
		// Name of the author
		//example: Muhammad Ali
		Name string `json:"name"`
		// The date when this author was the author of the day
		// example: 2021-06-12T00:00:00Z
		Date string `json:"date"`
	}
}

// swagger:response successResponse
type successResponseWrapper struct {
	// The successful response to a successful setting of a QOD
	// in: body
	Body struct {
		// Example: This request was a success
		Message string `json:"message"`
		// HTTP status code
		//
		// Example: 200
		StatusCode int `json:"statusCode"`
	}
}

// Data structure representing the error response to a wrongly structured request body
// swagger:response incorrectBodyStructureResponse
type incorrectBodyStructureResponseWrapper struct {
	// The error response to a wrongly structured request body
	// in: body
	Body struct {
		// The error message
		// Example: request body is not structured correctly.
		Message string `json:"message"`
	}
}

// Data structure representing the error response to an internal server error
// swagger:response internalServerErrorResponse
type internalServerErrorResponseWrapper struct {
	// The error response to an internal server
	// in: body
	Body struct {
		// The error message
		// Example: Please try again later.
		Message string `json:"message"`
	}
}

// Data structure representing the error response to an incorrect Credentials error
// swagger:response incorrectCredentialsResponse
type incorrectCredentialsResponseWrapper struct {
	// The error response to an unothorized access
	// in: body
	Body struct {
		// The error message
		// Example: Valar Dohaeris
		Message string `json:"message"`
	}
}

// Data structure for supported languages information
// swagger:response listOfStrings
type listOfStringsWrapper struct {
	// The languages supported by the api
	// in: body
	Body []struct {
		// The languages supported
		// example: ["English", "Icelandic"]
		Languages []string `json:"languages"`
	}
}

// Data structure representing a list response for topics
// swagger:response listTopicsResponse
type listTopicsResponseWrapper struct {
	// List of topics
	// in: body
	Body []struct {
		// The id of the topic
		// example: 10
		Id int `json:"id"`
		// Name of the topics
		// example: inspirational
		Name string `json:"name"`
		// Boolean whether or not this quote is in icelandic
		// example: true
		Isicelandic bool `json:"isicelandic"`
	}
}

// Data structure representing topic quotes response
// swagger:response multipleQuotesTopicResponse
type multipleQuotesTopicResponseWrapper struct {
	// List of quotes with their topic attached
	// in: body
	Body []struct {
		// The author's id
		//Unique: true
		//example: 26214
		Authorid int `json:"authorid"`
		// Name of author
		//example: John D. Rockefeller
		Name string `json:"name"`
		// The quote's id
		//Unique: true
		//example: 625402
		Quoteid int `json:"quoteid" `
		// The topic's id
		//Unique: true
		//example: 6
		Topicid int `json:"topicid" `
		// The topic's name
		// Unique: true
		// example: motivational
		Topicname string `json:"topicname"`
		// The quote
		//example: If you want to succeed you should strike out on new paths, rather than travel the worn paths of accepted success.
		Quote       string `json:"quote"`
		Isicelandic bool   `json:"-"`
	}
}

// Data structure representing a user response
// swagger:response userResponse
type userResponseWrapper struct {
	// The necessary data for the user
	// in: body
	Body structs.UserResponse
}
