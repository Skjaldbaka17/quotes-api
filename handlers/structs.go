package handlers

import "gorm.io/gorm"

type Quotes struct {
	gorm.Model
	Id          int    `json:"id"`
	Authorid    int    `json:"authorid"`
	Quote       string `json:"quote"`
	Count       int    `json:"count"`
	IsIcelandic bool   `json:"isicelandic"`
}

type Authors struct {
	gorm.Model
	Id     int      `json:"id"`
	Name   string   `json:"name"`
	Count  int      `json:"count"`
	Quotes []Quotes `json:"quotes" gorm:"foreignKey:authorid"`
}

type AuthorsView struct {
	// The author's id
	//Unique: true
	//example: 24952
	Id int `json:"id"`
	// Name of author
	//example: Muhammad Ali
	Name               string `json:"name"`
	Hasicelandicquotes bool   `json:"hasicelandicquotes"`
}

type QuoteView struct {
	// The author's id
	//Unique: true
	//example: 24952
	Authorid int `json:"authorid"`
	// Name of author
	//example: Muhammad Ali
	Name string `json:"name"`
	// The quote's id
	//Unique: true
	//example: 582676
	Quoteid int `json:"quoteid" `
	// The topic's id
	//Unique: true
	//example: 10
	Topicid int `json:"topicid" `
	// The topic's name
	//Unique: true
	//example: inspirational
	Topicname string `json:"topicname" `
	// The quote
	//example: Float like a butterfly, sting like a bee.
	Quote string `json:"quote"`
	// Whether or not this quote is in Icelandic or not
	// example: false
	Isicelandic bool `json:"isicelandic"`
	//swagger:ignore
	Id                 int  `json:"id"`
	Hasicelandicquotes bool `json:"hasicelandicquotes"`
}

type ListItem struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Isicelandic string `json:"isicelandic"`
}

type Request struct {
	Ids          []int       `json:"ids,omitempty"`
	Id           int         `json:"id,omitempty"`
	Page         int         `json:"page,omitempty"`
	SearchString string      `json:"searchString,omitempty"`
	PageSize     int         `json:"pageSize,omitempty"`
	Language     string      `json:"language,omitempty"`
	Topic        string      `json:"topic,omitempty"`
	AuthorId     int         `json:"authorId"`
	QuoteId      int         `json:"quoteId"`
	TopicId      int         `json:"topicId"`
	MaxQuotes    int         `json:"maxQuotes"`
	OrderConfig  OrderConfig `json:"orderConfig"`
}

type OrderConfig struct {
	// What to order by, 'alphabetical', 'popularity' or 'nrOfQuotes'
	// example: popularity
	OrderBy string `json:"orderBy"`
	// Where to start the ordering (if empty it starts from beginning, for example start at 'A' for alphabetical ascending order)
	// example: F
	StartFrom string `json:"startFrom"`
	// Whether to order the list in reverse or not
	// example: true
	Reverse bool `json:"reverse"`
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}
