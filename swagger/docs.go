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
