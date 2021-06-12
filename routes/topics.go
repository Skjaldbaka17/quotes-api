package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Skjaldbaka17/quotes-api/handlers"
	"github.com/Skjaldbaka17/quotes-api/structs"
	"gorm.io/gorm/clause"
)

// swagger:route POST /topics TOPICS GetTopics
// List the available topics, english / icelandic or both
// responses:
//	200: listTopicsResponse

// GetTopics handles POST requests for listing the available quote-topics
func GetTopics(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	var results []structs.ListItem
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := handlers.Db.Table("topics")

	dbPointer = quoteLanguageSQL(requestBody.Language, dbPointer)
	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err := dbPointer.Find(&results).Error
	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&results)
}

// swagger:route POST /topic TOPICS GetTopic
// Get quotes from a particular topic
// responses:
//	200: multipleQuotesTopicResponse

// GetTopic handles POST requests for getting quotes from a particular topic
func GetTopic(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	var results []structs.QuoteView
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	//Order by quoteid to have definitive order (when for examplke some quotes rank the same for plain, phrase and general)
	dbPoint := handlers.Db.Table("topicsview").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "quoteid DESC", Vars: []interface{}{}, WithoutParentheses: true},
	})

	if requestBody.Topic != "" {
		dbPoint = dbPoint.Where("lower(topicname) = lower(?)", requestBody.Topic)
	} else {
		dbPoint = dbPoint.Where("topicid = ?", requestBody.Id)
	}

	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err := pagination(requestBody, dbPoint).Find(&results).Error

	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//Update popularity in background!
	go handlers.DirectFetchTopicCountIncrement(requestBody.Id, requestBody.Topic)

	json.NewEncoder(rw).Encode(&results)
}
