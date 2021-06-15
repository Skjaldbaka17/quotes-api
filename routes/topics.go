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
//	200: topicsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// GetTopics handles POST requests for listing the available quote-topics
func GetTopics(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	var results []structs.TopicDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := handlers.Db.Table("topics")

	dbPointer = quoteLanguageSQL(requestBody.Language, dbPointer)
	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err := dbPointer.Find(&results).Error
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when querying DB in GetTopics: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	apiResults := structs.ConvertToTopicsAPIModel(results)
	json.NewEncoder(rw).Encode(apiResults)
}

// swagger:route POST /topic TOPICS GetTopic
// Get quotes from a particular topic
// responses:
//	200: topicViewsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// GetTopic handles POST requests for getting quotes from a particular topic
func GetTopic(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.Request
	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
		return
	}
	var results []structs.TopicViewDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	//Order by quoteid to have definitive order (when for examplke some quotes rank the same for plain, phrase and general)
	dbPoint := handlers.Db.Table("topicsview").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "quote_id DESC", Vars: []interface{}{}, WithoutParentheses: true},
	})

	if requestBody.Topic != "" {
		dbPoint = dbPoint.Where("lower(topic_name) = lower(?)", requestBody.Topic)
	} else {
		dbPoint = dbPoint.Where("topic_id = ?", requestBody.Id)
	}

	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err := pagination(requestBody, dbPoint).Find(&results).Error

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when querying DB in GetTopic: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	//Update popularity in background!
	go handlers.DirectFetchTopicCountIncrement(requestBody.Id, requestBody.Topic)
	topicViewsAPI := structs.ConvertToTopicViewsAPIModel(results)

	json.NewEncoder(rw).Encode(topicViewsAPI)
}
