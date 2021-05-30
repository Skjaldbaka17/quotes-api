// Package classification quotes-api.
//
// Documentation of our quotes API.
//	tags:
//		-name: QUOTES
//		description: Access random quote service. Use this to get random quotes , quotes filtered by authors or tags etc.
//
//     Schemes: http
//     BasePath: /api/
//     Version: 1.0.0
//     Host: quotel-api.com
//	   Contact: Þórður Ágústsson<skjaldbaka17@gmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - basic
//
//    SecurityDefinitions:
//    basic:
//      type: basic
//
// swagger:meta
package docs

import "github.com/Skjaldbaka17/quotes-api/handlers"

// Data structure representing most responses
// swagger:response multipleQuotesResponse
type authorsResponseWrapper struct {
	// List of authors / quotes
	// in: body
	Body []handlers.QuoteView
}

// Data structure representing the response for a random quote
// swagger:response randomQuoteResponse
type randomQuoteResponseWrapper struct {
	// A quote struct
	// in: body
	Body handlers.QuoteView
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

// Data structure for a list of strings
// swagger:response listOfStrings
type listOfStringsWrapper struct {
	// List of languages supported by the api
	// in: body
	Body []struct {
		// The languages supported
		// example: ["English", "Icelandic"]
		Languages []string `json:"languages"`
	}
}

// swagger:parameters getAuthorsByIds
type getAuthorsByIdsWrapper struct {
	// The structure of the request for authors by their ids
	// in: body
	// required: true
	Body struct {
		// The list of author's ids you want
		//
		// Required: true
		// Example: [24952,19161]
		Ids []int `json:"ids"`
	}
}

// swagger:parameters getQuotesByIds
type getQuotesByIdsWrapper struct {
	// The structure of the request for quotes by their ids
	// in: body
	// required: true
	Body struct {
		// The list of quotes's ids you want
		//
		// Required: true
		// Example: [582676,443976]
		Ids []int `json:"ids"`
	}
}

// swagger:parameters generalSearchByString searchQuotesByString
type getSearchByStringWrapper struct {
	// The structure of the request for searching quotes/authors
	// in: body
	// required: true
	Body struct {
		// The string to be used in the search
		//
		// Required: true
		// Example: sting like butterfly
		SearchString string `json:"searchString"`
		// The number of quotes to be returned on each "page"
		//
		// Maximum: 200
		// Minimum: 1
		// Default: 25
		// Example: 30
		PageSize int `json:"pageSize"`
		// The page you are asking for, starts with 0.
		//
		// Minimum: 0
		// Example: 0
		Page int `json:"page"`
		// The particular language that the quote should be in
		// example: English
		Language string `json:"language"`
	}
}

// swagger:parameters searchAuthorsByString
type getSearchAuthorsByStringWrapper struct {
	// The structure of the request for searching authors
	// in: body
	// required: true
	Body struct {
		// The string to be used in the search
		//
		// Required: true
		// Example: Ali Muhammad
		SearchString string `json:"searchString"`
		// The number of quotes to be returned on each "page"
		//
		// Maximum: 200
		// Minimum: 1
		// Default: 25
		// Example: 30
		PageSize int `json:"pageSize"`
		// The page you are asking for, starts with 0.
		//
		// Minimum: 0
		// Example: 0
		Page int `json:"page"`
		// The particular language that the quote should be in
		// example: English
		Language string `json:"language"`
	}
}

// swagger:parameters getTopics
type listTopicsWrapper struct {
	// The structure of the request for listing topics
	// in: body
	Body struct {
		// The language of the topics. If left empty all topics from all languages are returned
		//
		// Example: English
		Language string `json:"language"`
	}
}

// swagger:parameters getTopic
type quotesFromTopicWrapper struct {
	// The structure of the request for listing topics
	// in: body
	Body struct {
		// Name of the topic, if left empty then the id is used
		//
		// required: false
		// Example: Motivational
		Topic string `json:"topic"`
		// The topic's id, if left empty then the topic name is used
		//
		// Example: 10
		Id string `json:"id"`
		// The number of quotes to be returned on each "page"
		//
		// Maximum: 200
		// Minimum: 1
		// Default: 25
		// Example: 30
		PageSize int `json:"pageSize"`
		// The page you are asking for, starts with 0.
		//
		// Minimum: 0
		// Example: 0
		Page int `json:"page"`
	}
}

// swagger:parameters getRandomQuote
type getRandomQuoteResponseWrapper struct {
	// The structure of the response to random Quote post request
	// in: body
	Body struct {
		// The random quote returned must be in the given language
		//
		// Example: English
		Language string `json:"language"`
		// The random quote returned must contain a match with the searchstring
		//
		// Example: float
		SearchString string `json:"searchString"`
		// The random quote returned must be a part of the topic with the given topicId
		//
		// Example: 10
		TopicId string `json:"topicId"`
		// The random quote returned must be from the author with the given authorId
		//
		//example: 24952
		Authorid int `json:"authorid"`
	}
}
