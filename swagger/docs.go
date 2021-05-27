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

// Data structure representing a single product
// swagger:response multipleQuotesResponse
type authorsResponseWrapper struct {
	// Newly created product
	// in: body
	Body []handlers.SearchView
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
	}
}
