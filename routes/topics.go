package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Skjaldbaka17/quotes-api/handlers"
	"github.com/Skjaldbaka17/quotes-api/structs"
	"gorm.io/gorm/clause"
)

// swagger:route POST /topics TOPICS getTopics
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

	pointer := handlers.Db.Table("topics")

	switch strings.ToLower(requestBody.Language) {
	case "english":
		pointer = pointer.Not("isicelandic")
	case "icelandic":
		pointer = pointer.Where("isicelandic")
	}

	err := pointer.Find(&results).Error
	log.Println(results)
	if err != nil {
		//TODO: Respond with better error -- and put into swagger -- and add tests
		log.Printf("Got error when decoding: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(&results)
}

// swagger:route POST /topic TOPICS getTopic
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

	//Order by quoteid to have definitive order (when for examplke some quotes rank the same for plain, phrase and general)
	dbPoint := handlers.Db.Table("topicsview")

	if requestBody.Topic != "" {
		dbPoint = dbPoint.Where("lower(topicname) = lower(?)", requestBody.Topic)
	} else {
		dbPoint = dbPoint.Where("topicid = ?", requestBody.Id)
	}

	err := dbPoint.Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "quoteid DESC", Vars: []interface{}{}, WithoutParentheses: true},
	}).
		Limit(requestBody.PageSize).
		Offset(requestBody.Page * requestBody.PageSize).
		Find(&results).Error

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