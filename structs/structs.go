package structs

type QuoteView struct {
	// The author's id
	//Unique: true
	//example: 24952
	Authorid int `json:"author_id"`
	// Name of author
	//example: Muhammad Ali
	Name string `json:"name"`
	// The quote's id
	//Unique: true
	//example: 582676
	Quoteid int `json:"quote_id" `
	// The topic's id
	//Unique: true
	//example: 10
	Topicid int `json:"topic_id" `
	// The topic's name
	//Unique: true
	//example: inspirational
	Topicname string `json:"topic_name" `
	// The quote
	//example: Float like a butterfly, sting like a bee.
	Quote string `json:"quote"`
	// Whether or not this quote is in Icelandic or not
	// example: false
	Isicelandic bool `json:"is_icelandic"`
	//swagger:ignore
	Id int `json:"id"`
	//swagger:ignore
	Hasicelandicquotes bool `json:"has_icelandic_quotes"`
	//swagger:ignore
	Nroficelandicquotes int `json:"nr_of_icelandic_quotes"`
	//swagger:ignore
	Nrofenglishquotes int `json:"nr_of_english_quotes"`
	//swagger:ignore Date
	Date string `json:"date"`
	//swagger:ignore
	Quotecount int `json:"quote_count"`
	//swagger:ignore
	Authorcount int `json:"author_count"`
}

type ListItem struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Isicelandic string `json:"is_icelandic"`
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
	ApiKey       string      `json:"apiKey"`
}

type RequestEvent struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Route       string `json:"route"`
	RequestBody string `json:"request_body"`
	ApiKey      string `json:"api_key"`
}

type ErrorEvent struct {
	Id           int    `json:"id"`
	UserId       int    `json:"user_id"`
	Route        string `json:"route"`
	RequestBody  string `json:"request_body"`
	ErrorMessage string `json:"error_message"`
	ExtraInfo    string `json:"extra_info"`
}

type UserResponse struct {
	// The user id
	// example: 1
	Id int `json:"id"`
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
	StatusCode int    `json:"status_code"`
}
