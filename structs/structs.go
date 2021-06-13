package structs

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
	// Name of the author
	//example: Muhammad Ali
	Name string `json:"name"`
	// Whether or not this author has some icelandic quotes
	//example: true
	Hasicelandicquotes bool `json:"hasicelandicquotes"`
	// How many quotes in icelandic this author has
	//example: 6
	Nroficelandicquotes int `json:"nroficelandicquotes"`
	// How many quotes in icelandic this author has
	//example: 78
	Nrofenglishquotes int `json:"nrofenglishquotes"`
	//swagger:ignore
	Count int `json:"count"`
	//swagger:ignore
	Date string `json:"date"`
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
	Id int `json:"id"`
	//swagger:ignore
	Hasicelandicquotes bool `json:"hasicelandicquotes"`
	//swagger:ignore
	Nroficelandicquotes int `json:"nroficelandicquotes"`
	//swagger:ignore
	Nrofenglishquotes int `json:"nrofenglishquotes"`
	//swagger:ignore Date
	Date string `json:"date"`
	//swagger:ignore
	Quotecount int `json:"quotecount"`
	//swagger:ignore
	Authorcount int `json:"authorcount"`
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
	Date         string      `json:"date"`
	Minimum      string      `json:"minimum"`
	Maximum      string      `json:"maximum"`
	Qods         []Qod       `json:"qods"`
	Aods         []Qod       `json:"aods"`
}

type UserRequest struct {
	ApiKey               string `json:"apiKey"`
	Id                   int    `json:"id"`
	Email                string `json:"email"`
	Name                 string `json:"name"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
	Tier                 string `json:"tier"`
}

type User struct {
	Id           int    `json:"id"`
	ApiKey       string `json:"api_key"`
	Message      string `json:"message"`
	PasswordHash string `json:"password_hash"`
	Tier         string `json:"tier"`
	Name         string `json:"name"`
	Email        string `json:"email"`
}

type UserResponse struct {
	// The api-key that the user should send to get access to the api
	// example: 1d8db1d2-6f5b-4254-8b74-44f5e5229add
	ApiKey string `json:"api_key"`
}

type Qod struct {
	// the date for which this quote is the QOD, if left empty this quote is today's QOD.
	//
	// Example: 12-22-2020
	Date string `json:"date"`
	// The id of the quote to be set as this dates QOD
	//
	// Example: 1
	Id int `json:"id"`
	// The language of the QOD
	// Example: icelandic
	Language string `json:"language"`
}

type OrderConfig struct {
	// What to order by, 'alphabetical', 'popularity' or 'nrOfQuotes'
	// example: popularity
	OrderBy string `json:"orderBy"`
	// Where to start the ordering (if empty it starts from beginning, for example start at 'A' for alphabetical ascending order)
	// example: F
	Minimum string `json:"minimum"`
	// Where to end the ordering (if empty it ends at the logical end, for example end at 'Z' for alphabetical ascending order)
	// example: Z
	Maximum string `json:"maximum"`
	// Whether to order the list in reverse or not (true is Descending and false is Ascending, false is default)
	// example: true
	Reverse bool `json:"reverse"`
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}
