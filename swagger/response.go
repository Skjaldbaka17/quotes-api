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
