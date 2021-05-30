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
	Name   string   `json:"name"`
	Count  int      `json:"count"`
	Quotes []Quotes `json:"quotes" gorm:"foreignKey:authorid"`
}

type SearchView struct {
	//swagger:ignore AuthorId then gorm cant map the values correctly, but works with Authorid and Quoteid etc. Why? TODO

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
	// The quote
	//example: Float like a butterfly, sting like a bee.
	Quote       string `json:"quote"`
	Isicelandic bool   `json:"isicelandic"`
}

type ListItem struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Isicelandic string `json:"isicelandic"`
}

type Request struct {
	Ids          []int  `json:"ids,omitempty"`
	Id           int    `json:"id,omitempty"`
	Page         int    `json:"page,omitempty"`
	SearchString string `json:"searchString,omitempty"`
	PageSize     int    `json:"pageSize,omitempty"`
	Language     string `json:"language,omitempty"`
	Topic        string `json:"topic,omitempty"`
}

type TopicsView struct {
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
	//example: 582676
	Topicid int `json:"topicid" `
	// The topic's name
	//Unique: true
	//example: inspirational
	Topicname string `json:"topicname" `
	// The quote
	//example: Float like a butterfly, sting like a bee.
	Quote       string `json:"quote"`
	Isicelandic bool   `json:"-"`
}
