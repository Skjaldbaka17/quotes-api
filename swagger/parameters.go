package docs

// swagger:parameters GetAuthors
type getAuthorsWrapper struct {
	// The structure of the request for getting authors by their ids
	// in: body
	// required: true
	Body struct {
		// A list of the authors's ids that you want
		//
		// Required: true
		// Example: [24952,19161]
		Ids []int `json:"ids"`
	}
}

// swagger:parameters ListAuthors
type authorsListWrapper struct {
	// The structure of the request for getting a list of authors
	// in: body
	Body struct {
		// Response is paged. This parameter controls the number of Authors to be returned on each "page"
		//
		// Maximum: 200
		// Minimum: 1
		// Default: 25
		// Example: 30
		PageSize int `json:"pageSize"`
		// Response is paged. This parameter controls the page you are asking for, starts with 0.
		//
		// Minimum: 0
		// Example: 0
		Page int `json:"page"`
		// Only return authors that have quotes in the given language ("english" or "icelandic") if left empty then no constraint
		// is set on the quotes' language. Note if ordering by nrOfQuotes if this parameter is set then only the amount of
		// quotes the author has in the given language counts towards the final ordering.
		// Example: English
		Language string `json:"language"`
		//Model
		OrderConfig orderConfigListAuthorsModel
	}
}

// swagger:parameters GetRandomAuthor
type randomAuthorWrapper struct {
	// The structure of the request for getting a random author
	// in: body
	Body struct {
		// The random author must have quotes in the given language ("english" or "icelandic") if left empty then no
		// constraint on language is set
		//
		// Example: English
		Language string `json:"language"`
		// How many quotes, maximum, to be returned from this author
		//
		// Example: 10
		// Maximum: 50
		// default: 1
		MaxQuotes int `json:"maxQuotes"`
	}
}

// swagger:parameters GetAuthorOfTheDay
type authorOfTheDayWrapper struct {
	// The structure of the request for getting the author of the day
	// in: body
	Body struct {
		// Get the author of the day for the given language ("icelandic" or "english")
		//
		// Default: English
		// Example: English
		Language string `json:"language"`
	}
}

// swagger:parameters GetAODHistory
type historyAODWrapper struct {
	// The structure of the request for getting the history of AODs
	// in: body
	Body []struct {
		// Get the history of theAODS for the given language ("icelandic" or "english")
		//
		// Default: English
		// Example: icelandic
		Language string `json:"language"`
		// The earliest date to return. All authors between minimum and today will be returned.
		// Example: 2020-12-21
		Minimum string `json:"minimum"`
	}
}

// swagger:parameters SetAuthorOfTheDay
type setAODWrapper struct {
	// The structure of the request for setting AODs
	// in: body
	Body []struct {
		Aods []ofTheDayModel
	}
}
