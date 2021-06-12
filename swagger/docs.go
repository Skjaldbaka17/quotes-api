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

// swagger:parameters getQODHistory
type getQODHistoryWrapper struct {
	// The structure of the request for getting the QOD history
	// in: body
	Body struct {
		// The language of the QOD. If left empty the english QOD is returned
		//
		// Example: English
		Language string `json:"language"`
		//The minimum date to retrieve the history
		// example: 2020-06-21
		Minimum string `json:"minimum"`
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
